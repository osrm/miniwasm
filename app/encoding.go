package app

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/initia-labs/initia/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	tempApp := NewMinitiaApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, NewEmptyAppOptions())
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	return encodingConfig
}

func AutoCliOpts() autocli.AppOptions {
	tempApp := NewMinitiaApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, NewEmptyAppOptions())
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range tempApp.ModuleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(tempApp.ModuleManager.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

func BasicManager() module.BasicManager {
	tempApp := NewMinitiaApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, NewEmptyAppOptions())
	return tempApp.BasicModuleManager
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct {
	homeDir string
}

func NewEmptyAppOptions() EmptyAppOptions {
	return EmptyAppOptions{
		homeDir: filepath.Join(os.TempDir(), strconv.Itoa(rand.Int())),
	}
}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	if o == flags.FlagHome {
		if ao.homeDir == "" {
			return DefaultNodeHome
		}

		return ao.homeDir
	}

	return nil
}
