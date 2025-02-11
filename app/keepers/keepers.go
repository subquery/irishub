package keepers

import (
	"github.com/spf13/cast"

	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	srvflags "github.com/evmos/ethermint/server/flags"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	nfttransfer "github.com/bianjieai/nft-transfer"
	ibcnfttransferkeeper "github.com/bianjieai/nft-transfer/keeper"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	tibcmttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer"
	tibcmttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/keeper"
	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer"
	tibcnfttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibcroutingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibccli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"
	tibckeeper "github.com/bianjieai/tibc-go/modules/tibc/core/keeper"

	coinswapkeeper "github.com/irisnet/irismod/modules/coinswap/keeper"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/modules/farm"
	farmkeeper "github.com/irisnet/irismod/modules/farm/keeper"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	mtkeeper "github.com/irisnet/irismod/modules/mt/keeper"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	nftkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	randomkeeper "github.com/irisnet/irismod/modules/random/keeper"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	recordkeeper "github.com/irisnet/irismod/modules/record/keeper"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"

	guardiankeeper "github.com/irisnet/irishub/v3/modules/guardian/keeper"
	guardiantypes "github.com/irisnet/irishub/v3/modules/guardian/types"
	"github.com/irisnet/irishub/v3/modules/internft"
	mintkeeper "github.com/irisnet/irishub/v3/modules/mint/keeper"
	minttypes "github.com/irisnet/irishub/v3/modules/mint/types"
	iristypes "github.com/irisnet/irishub/v3/types"
)

// AppKeepers defines a structure used to consolidate all
// the keepers needed to run an iris app.
type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	scopedIBCKeeper         capabilitykeeper.ScopedKeeper
	scopedTransferKeeper    capabilitykeeper.ScopedKeeper
	scopedIBCMockKeeper     capabilitykeeper.ScopedKeeper
	scopedNFTTransferKeeper capabilitykeeper.ScopedKeeper
	scopedICAHostKeeper     capabilitykeeper.ScopedKeeper
	scopedTIBCKeeper        capabilitykeeper.ScopedKeeper
	scopedTIBCMockKeeper    capabilitykeeper.ScopedKeeper

	FeeGrantKeeper        feegrantkeeper.Keeper
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	EvidenceKeeper        *evidencekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCTransferKeeper     ibctransferkeeper.Keeper
	IBCNFTTransferKeeper  ibcnfttransferkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	GuardianKeeper        guardiankeeper.Keeper
	TokenKeeper           tokenkeeper.Keeper
	RecordKeeper          recordkeeper.Keeper
	NFTKeeper             nftkeeper.Keeper
	MTKeeper              mtkeeper.Keeper
	HTLCKeeper            htlckeeper.Keeper
	CoinswapKeeper        coinswapkeeper.Keeper
	ServiceKeeper         servicekeeper.Keeper
	OracleKeeper          oraclekeeper.Keeper
	RandomKeeper          randomkeeper.Keeper
	FarmKeeper            farmkeeper.Keeper
	TIBCKeeper            *tibckeeper.Keeper
	TIBCNFTTransferKeeper tibcnfttransferkeeper.Keeper
	TIBCMTTransferKeeper  tibcmttransferkeeper.Keeper
	EvmKeeper             *evmkeeper.Keeper
	FeeMarketKeeper       feemarketkeeper.Keeper

	TransferModule       transfer.AppModule
	ICAModule            ica.AppModule
	NftTransferModule    tibcnfttransfer.AppModule
	MtTransferModule     tibcmttransfer.AppModule
	IBCNftTransferModule nfttransfer.AppModule
}

// New initializes a new instance of AppKeepers.
//
// It takes in various parameters including appCodec, bApp, legacyAmino, maccPerms, modAccAddrs, blockedAddress, skipUpgradeHeights, homePath, invCheckPeriod, logger, and appOpts.
// It returns an instance of AppKeepers.
func New(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	modAccAddrs map[string]bool,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	logger log.Logger,
	appOpts servertypes.AppOptions,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.genStoreKeys()

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, appKeepers.keys); err != nil {
		tmos.Exit(err.Error())
	}

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)
	appKeepers.ConsensusParamsKeeper = consensuskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[consensustypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)
	appKeepers.scopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.scopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.scopedNFTTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcnfttransfertypes.ModuleName)
	appKeepers.scopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		ethermint.ProtoAccount,
		maccPerms,
		iristypes.Bech32PrefixAccAddr,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		blockedAddress,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[crisistypes.StoreKey],
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	// set the governance module account as the authority for conducting upgrades
	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	appKeepers.scopedTIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(tibchost.ModuleName)
	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.scopedIBCKeeper,
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.ICAModule = ica.NewAppModule(nil, &appKeepers.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(appKeepers.ICAHostKeeper)

	// register the proposal types
	appKeepers.TIBCKeeper = tibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[tibchost.StoreKey],
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.NFTKeeper = nftkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[nfttypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.TIBCNFTTransferKeeper = tibcnfttransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[tibcnfttypes.StoreKey],
		appKeepers.GetSubspace(tibcnfttypes.ModuleName),
		appKeepers.AccountKeeper,
		nftkeeper.NewLegacyKeeper(appKeepers.NFTKeeper),
		appKeepers.TIBCKeeper.PacketKeeper,
		appKeepers.TIBCKeeper.ClientKeeper,
	)

	appKeepers.MTKeeper = mtkeeper.NewKeeper(
		appCodec, appKeepers.keys[mttypes.StoreKey],
	)

	appKeepers.TIBCMTTransferKeeper = tibcmttransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[tibcnfttypes.StoreKey],
		appKeepers.GetSubspace(tibcnfttypes.ModuleName),
		appKeepers.AccountKeeper, appKeepers.MTKeeper,
		appKeepers.TIBCKeeper.PacketKeeper,
		appKeepers.TIBCKeeper.ClientKeeper,
	)

	appKeepers.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.scopedTransferKeeper,
	)
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.IBCTransferKeeper)
	transferIBCModule := transfer.NewIBCModule(appKeepers.IBCTransferKeeper)

	appKeepers.IBCNFTTransferKeeper = ibcnfttransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcnfttransfertypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		internft.NewInterNftKeeper(appCodec, appKeepers.NFTKeeper, appKeepers.AccountKeeper),
		appKeepers.scopedNFTTransferKeeper,
	)
	appKeepers.IBCNftTransferModule = nfttransfer.NewAppModule(appKeepers.IBCNFTTransferKeeper)
	nfttransferIBCModule := nfttransfer.NewIBCModule(appKeepers.IBCNFTTransferKeeper)

	// routerModule := router.NewAppModule(app.RouterKeeper, transferIBCModule)
	// create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(ibcnfttransfertypes.ModuleName, nfttransferIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	appKeepers.NftTransferModule = tibcnfttransfer.NewAppModule(appKeepers.TIBCNFTTransferKeeper)
	appKeepers.MtTransferModule = tibcmttransfer.NewAppModule(appKeepers.TIBCMTTransferKeeper)

	tibcRouter := tibcroutingtypes.NewRouter()
	tibcRouter.AddRoute(tibcnfttypes.ModuleName, appKeepers.NftTransferModule).
		AddRoute(tibcmttypes.ModuleName, appKeepers.MtTransferModule)
	appKeepers.TIBCKeeper.SetRouter(tibcRouter)

	// create evidence keeper with router
	appKeepers.EvidenceKeeper = evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	appKeepers.GuardianKeeper = guardiankeeper.NewKeeper(
		appCodec,
		appKeepers.keys[guardiantypes.StoreKey],
	)

	appKeepers.TokenKeeper = tokenkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[tokentypes.StoreKey],
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	).WithSwapRegistry(tokenv1.SwapRegistry{
		iristypes.NativeToken.MinUnit: tokenv1.SwapParams{
			MinUnit: iristypes.EvmToken.MinUnit,
			Ratio:   sdk.OneDec(),
		},
		iristypes.EvmToken.MinUnit: tokenv1.SwapParams{
			MinUnit: iristypes.NativeToken.MinUnit,
			Ratio:   sdk.OneDec(),
		},
	})

	appKeepers.RecordKeeper = recordkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[recordtypes.StoreKey],
	)

	appKeepers.HTLCKeeper = htlckeeper.NewKeeper(
		appCodec, appKeepers.keys[htlctypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.CoinswapKeeper = coinswapkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[coinswaptypes.StoreKey],
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.ServiceKeeper = servicekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[servicetypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		servicetypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[oracletypes.StoreKey],
		appKeepers.ServiceKeeper,
	)

	appKeepers.RandomKeeper = randomkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[randomtypes.StoreKey],
		appKeepers.BankKeeper,
		appKeepers.ServiceKeeper,
	)

	govConfig := govtypes.DefaultConfig()
	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.FarmKeeper = farmkeeper.NewKeeper(appCodec,
		appKeepers.keys[farmtypes.StoreKey],
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
		appKeepers.GovKeeper,
		appKeepers.CoinswapKeeper,
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(tibchost.RouterKey, tibccli.NewProposalHandler(appKeepers.TIBCKeeper)).
		AddRoute(farmtypes.RouterKey, farm.NewCommunityPoolCreateFarmProposalHandler(appKeepers.FarmKeeper))

	appKeepers.GovKeeper.SetHooks(govtypes.NewMultiGovHooks(
		farmkeeper.NewGovHook(appKeepers.FarmKeeper),
	))

	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	// Create Ethermint keepers
	appKeepers.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		appKeepers.keys[feemarkettypes.StoreKey],
		appKeepers.tkeys[feemarkettypes.TransientKey],
		appKeepers.GetSubspace(feemarkettypes.ModuleName),
	)
	appKeepers.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evmtypes.StoreKey],
		appKeepers.tkeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.FeeMarketKeeper,
		nil,
		geth.NewEVM,
		cast.ToString(appOpts.Get(srvflags.EVMTracer)),
		appKeepers.GetSubspace(evmtypes.ModuleName),
	)
	return appKeepers
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	paramsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())
	paramsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())
	paramsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(tokentypes.ModuleName).WithKeyTable(tokenv1.ParamKeyTable())
	paramsKeeper.Subspace(recordtypes.ModuleName)
	paramsKeeper.Subspace(htlctypes.ModuleName).WithKeyTable(htlctypes.ParamKeyTable())
	paramsKeeper.Subspace(coinswaptypes.ModuleName).WithKeyTable(coinswaptypes.ParamKeyTable())
	paramsKeeper.Subspace(servicetypes.ModuleName).WithKeyTable(servicetypes.ParamKeyTable())
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(farmtypes.ModuleName).WithKeyTable(farmtypes.ParamKeyTable())
	paramsKeeper.Subspace(tibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName).WithKeyTable(evmtypes.ParamKeyTable())
	paramsKeeper.Subspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())

	return paramsKeeper
}
