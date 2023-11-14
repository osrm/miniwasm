package main

import (
	"fmt"
	"time"

	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/initia-labs/miniwasm/types"
)

// minitiaAppConfig initia specify app config
type minitiaAppConfig struct {
	serverconfig.Config
	WasmConfig wasmtypes.WasmConfig `mapstructure:"wasm"`
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = fmt.Sprintf("0%s", types.BaseDenom)

	minitiaAppConfig := minitiaAppConfig{
		Config:     *srvCfg,
		WasmConfig: wasmtypes.DefaultWasmConfig(),
	}

	minitiaAppTemplate := serverconfig.DefaultConfigTemplate +
		wasmtypes.DefaultConfigTemplate()

	return minitiaAppTemplate, minitiaAppConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// empty block configure
	cfg.Consensus.CreateEmptyBlocks = false
	cfg.Consensus.CreateEmptyBlocksInterval = time.Minute

	// block time from 5s to 0.5s
	cfg.Consensus.TimeoutPropose = cfg.Consensus.TimeoutPropose / 10
	cfg.Consensus.TimeoutProposeDelta = cfg.Consensus.TimeoutProposeDelta / 10
	cfg.Consensus.TimeoutPrevote = cfg.Consensus.TimeoutPrevote / 10
	cfg.Consensus.TimeoutPrevoteDelta = cfg.Consensus.TimeoutPrevoteDelta / 10
	cfg.Consensus.TimeoutPrecommit = cfg.Consensus.TimeoutPrecommit / 10
	cfg.Consensus.TimeoutPrecommitDelta = cfg.Consensus.TimeoutPrecommitDelta / 10
	cfg.Consensus.TimeoutCommit = cfg.Consensus.TimeoutCommit / 10

	return cfg
}
