package wasm_hooks_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecaddress "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/skip-mev/slinky/x/oracle"
	oraclekeeper "github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	wasm.AppModuleBasic{},
	oracle.AppModuleBasic{},
)

var (
	initiaSupply = math.NewInt(100_000_000_000)
	testDenoms   = []string{
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	}
)

type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

func MakeTestCodec(t testing.TB) codec.Codec {
	return MakeEncodingConfig(t).Codec
}

func MakeEncodingConfig(_ testing.TB) EncodingConfig {
	interfaceRegistry, _ := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          codecaddress.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
			ValidatorAddressCodec: codecaddress.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		},
	})
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	txConfig := tx.NewTxConfig(appCodec, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(legacyAmino)

	ModuleBasics.RegisterLegacyAminoCodec(legacyAmino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             appCodec,
		TxConfig:          txConfig,
		Amino:             legacyAmino,
	}
}

var bondDenom = sdk.DefaultBondDenom

func initialTotalSupply() sdk.Coins {
	faucetBalance := sdk.NewCoins(sdk.NewCoin(bondDenom, initiaSupply))
	for _, testDenom := range testDenoms {
		faucetBalance = faucetBalance.Add(sdk.NewCoin(testDenom, initiaSupply))
	}

	return faucetBalance
}

type TestFaucet struct {
	t                testing.TB
	bankKeeper       bankkeeper.Keeper
	sender           sdk.AccAddress
	balance          sdk.Coins
	minterModuleName string
}

func NewTestFaucet(t testing.TB, ctx sdk.Context, bankKeeper bankkeeper.Keeper, minterModuleName string, initiaSupply ...sdk.Coin) *TestFaucet {
	require.NotEmpty(t, initiaSupply)
	r := &TestFaucet{t: t, bankKeeper: bankKeeper, minterModuleName: minterModuleName}
	_, _, addr := keyPubAddr()
	r.sender = addr
	r.Mint(ctx, addr, initiaSupply...)
	r.balance = initiaSupply
	return r
}

func (f *TestFaucet) Mint(parentCtx sdk.Context, addr sdk.AccAddress, amounts ...sdk.Coin) {
	amounts = sdk.Coins(amounts).Sort()
	require.NotEmpty(f.t, amounts)
	ctx := parentCtx.WithEventManager(sdk.NewEventManager()) // discard all faucet related events
	err := f.bankKeeper.MintCoins(ctx, f.minterModuleName, amounts)
	require.NoError(f.t, err)
	err = f.bankKeeper.SendCoinsFromModuleToAccount(ctx, f.minterModuleName, addr, amounts)
	require.NoError(f.t, err)
	f.balance = f.balance.Add(amounts...)
}

func (f *TestFaucet) Fund(parentCtx sdk.Context, receiver sdk.AccAddress, amounts ...sdk.Coin) {
	require.NotEmpty(f.t, amounts)
	// ensure faucet is always filled
	if !f.balance.IsAllGTE(amounts) {
		f.Mint(parentCtx, f.sender, amounts...)
	}
	ctx := parentCtx.WithEventManager(sdk.NewEventManager()) // discard all faucet related events
	err := f.bankKeeper.SendCoins(ctx, f.sender, receiver, amounts)
	require.NoError(f.t, err)
	f.balance = f.balance.Sub(amounts...)
}

func (f *TestFaucet) NewFundedAccount(ctx sdk.Context, amounts ...sdk.Coin) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	f.Fund(ctx, addr, amounts...)
	return addr
}

type TestKeepers struct {
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	WasmKeeper    wasmkeeper.Keeper
	OracleKeeper  oraclekeeper.Keeper

	EncodingConfig EncodingConfig
	Faucet         *TestFaucet
	MultiStore     storetypes.CommitMultiStore
}

// createDefaultTestInput common settings for createTestInput
func createDefaultTestInput(t testing.TB) (sdk.Context, TestKeepers) {
	return createTestInput(t, false)
}

// createTestInput encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func createTestInput(t testing.TB, isCheckTx bool) (sdk.Context, TestKeepers) {
	// Load default move config
	return _createTestInput(t, isCheckTx, dbm.NewMemDB())
}

var keyCounter uint64

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func _createTestInput(
	t testing.TB,
	isCheckTx bool,
	db dbm.DB,
) (sdk.Context, TestKeepers) {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distributiontypes.StoreKey, wasmtypes.StoreKey,
		oracletypes.StoreKey,
	)
	ms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	for _, v := range keys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeIAVL, db)
	}
	memKeys := storetypes.NewMemoryStoreKeys()
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger()).WithHeaderHash(make([]byte, 32))

	encodingConfig := MakeEncodingConfig(t)
	appCodec := encodingConfig.Codec

	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName: nil,

		// for testing
		authtypes.Minter: {authtypes.Minter, authtypes.Burner},
	}

	ac := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]), // target store
		authtypes.ProtoBaseAccount,                          // prototype
		maccPerms,
		ac,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ctx.Logger().With("module", "x/"+banktypes.ModuleName),
	)
	require.NoError(t, bankKeeper.SetParams(ctx, banktypes.DefaultParams()))

	oracleKeeper := oraclekeeper.NewKeeper(
		runtime.NewKVStoreService(keys[oracletypes.StoreKey]),
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	queryAllowlist := make(map[string]proto.Message)
	queryAllowlist["/slinky.oracle.v1.Query/GetAllCurrencyPairs"] = &oracletypes.GetAllCurrencyPairsResponse{}
	queryAllowlist["/slinky.oracle.v1.Query/GetPrice"] = &oracletypes.GetPriceResponse{}
	queryAllowlist["/slinky.oracle.v1.Query/GetPrices"] = &oracletypes.GetPricesResponse{}

	// use accept list stargate querier
	wasmOpts := []wasmkeeper.Option{}
	wasmOpts = append(wasmOpts, wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: wasmkeeper.AcceptListStargateQuerier(queryAllowlist, queryRouter, appCodec),
	}))

	wasmKeeper := wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		msgRouter,
		queryRouter,
		t.TempDir(),
		wasmtypes.DefaultWasmConfig(),
		"iterator,stargate,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3,cosmwasm_1_4",
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)
	wasmParams := wasmtypes.DefaultParams()
	require.NoError(t, wasmKeeper.SetParams(ctx, wasmParams))

	faucet := NewTestFaucet(t, ctx, bankKeeper, authtypes.Minter, initialTotalSupply()...)

	// register query service
	am := module.NewManager( // minimal module set that we use for message/ query tests
		oracle.NewAppModule(appCodec, oracleKeeper),
	)
	am.RegisterServices(module.NewConfigurator(appCodec, msgRouter, queryRouter)) //nolint:errcheck

	keepers := TestKeepers{
		AccountKeeper:  accountKeeper,
		WasmKeeper:     wasmKeeper,
		BankKeeper:     bankKeeper,
		OracleKeeper:   oracleKeeper,
		EncodingConfig: encodingConfig,
		Faucet:         faucet,
		MultiStore:     ms,
	}
	return ctx, keepers
}
