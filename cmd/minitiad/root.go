package main

import (
	"errors"
	"io"
	"os"
	"path"

	rosettaCmd "cosmossdk.io/tools/rosetta/cmd"

	dbm "github.com/cometbft/cometbft-db"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	cosmosgenutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	minitiaapp "github.com/initia-labs/miniwasm/app"
	"github.com/initia-labs/miniwasm/app/params"
)

// NewRootCmd creates a new root command for initiad. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := minitiaapp.MakeEncodingConfig()

	sdkConfig := sdk.GetConfig()
	sdkConfig.SetCoinType(minitiaapp.CoinType)

	accountPubKeyPrefix := minitiaapp.AccountAddressPrefix + "pub"
	validatorAddressPrefix := minitiaapp.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := minitiaapp.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := minitiaapp.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := minitiaapp.AccountAddressPrefix + "valconspub"

	sdkConfig.SetBech32PrefixForAccount(minitiaapp.AccountAddressPrefix, accountPubKeyPrefix)
	sdkConfig.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	sdkConfig.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	sdkConfig.SetAddressVerifier(wasmtypes.VerifyAddressLen()) // same with wasm one
	sdkConfig.Seal()

	// Get the executable name and configure the viper instance so that environmental
	// variables are checked based off that name. The underscore character is used
	// as a separator
	executableName, err := os.Executable()
	if err != nil {
		panic(err)
	}

	basename := path.Base(executableName)

	// Configure the viper instance
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(minitiaapp.DefaultNodeHome).
		WithViper(minitiaapp.EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   basename,
		Short: "minitia App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			// read envs before reading persistent flags
			// TODO - should we handle this for tx flags & query flags?
			initClientCtx, err := readEnv(initClientCtx)
			if err != nil {
				return err
			}

			// read persistent flags if they changed, and override the env configs.
			initClientCtx, err = client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// unsafe-reset-all is not working without viper set
			viper.Set(tmcli.HomeFlag, initClientCtx.HomeDir)

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			initiaappTemplate, initiaappConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, initiaappTemplate, initiaappConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	// TODO check gaia before make release candidate
	// authclient.Codec = encodingConfig.Marshaler

	rootCmd.AddCommand(
		InitCmd(minitiaapp.ModuleBasics, minitiaapp.DefaultNodeHome),
		AddGenesisAccountCmd(minitiaapp.DefaultNodeHome),
		AddGenesisValidatorCmd(minitiaapp.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, minitiaapp.DefaultNodeHome),
		cosmosgenutilcli.ValidateGenesisCmd(minitiaapp.ModuleBasics),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
	)

	a := appCreator{encodingConfig}
	server.AddCommands(rootCmd, minitiaapp.DefaultNodeHome, a.newApp, a.appExport, addModuleInitFlags)
	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(minitiaapp.DefaultNodeHome),
	)

	// add rosetta commands
	rootCmd.AddCommand(rosettaCmd.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Marshaler))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	wasm.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	minitiaapp.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
	)

	minitiaapp.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encodingConfig params.EncodingConfig
}

// newApp is an AppCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	return minitiaapp.NewMinitiaApp(
		logger, db, traceStore, true,
		wasmOpts,
		appOpts,
		baseappOptions...,
	)
}

func (a appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	var initiaApp *minitiaapp.MinitiaApp
	if height != -1 {
		initiaApp = minitiaapp.NewMinitiaApp(logger, db, traceStore, false, []wasmkeeper.Option{}, appOpts)

		if err := initiaApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		initiaApp = minitiaapp.NewMinitiaApp(logger, db, traceStore, true, []wasmkeeper.Option{}, appOpts)
	}

	return initiaApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func readEnv(clientCtx client.Context) (client.Context, error) {
	if outputFormat := clientCtx.Viper.GetString(tmcli.OutputFlag); outputFormat != "" {
		clientCtx = clientCtx.WithOutputFormat(outputFormat)
	}

	if homeDir := clientCtx.Viper.GetString(flags.FlagHome); homeDir != "" {
		clientCtx = clientCtx.WithHomeDir(homeDir)
	}

	if clientCtx.Viper.GetBool(flags.FlagDryRun) {
		clientCtx = clientCtx.WithSimulation(true)
	}

	if keyringDir := clientCtx.Viper.GetString(flags.FlagKeyringDir); keyringDir != "" {
		clientCtx = clientCtx.WithKeyringDir(clientCtx.Viper.GetString(flags.FlagKeyringDir))
	}

	if chainID := clientCtx.Viper.GetString(flags.FlagChainID); chainID != "" {
		clientCtx = clientCtx.WithChainID(chainID)
	}

	if keyringBackend := clientCtx.Viper.GetString(flags.FlagKeyringBackend); keyringBackend != "" {
		kr, err := client.NewKeyringFromBackend(clientCtx, keyringBackend)
		if err != nil {
			return clientCtx, err
		}

		clientCtx = clientCtx.WithKeyring(kr)
	}

	if nodeURI := clientCtx.Viper.GetString(flags.FlagNode); nodeURI != "" {
		clientCtx = clientCtx.WithNodeURI(nodeURI)

		client, err := client.NewClientFromNode(nodeURI)
		if err != nil {
			return clientCtx, err
		}

		clientCtx = clientCtx.WithClient(client)
	}

	return clientCtx, nil
}
