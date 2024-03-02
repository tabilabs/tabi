package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgrades "github.com/tabilabs/tabi/app/upgrades"
)

var (
	plans = []upgrades.Upgrade{}
)

// RegisterUpgradePlans register a handler of upgrade plan
func (app *Tabi) RegisterUpgradePlans() {
	for _, u := range plans {
		app.registerUpgradeHandler(u.UpgradeName,
			u.StoreUpgrades,
			u.UpgradeHandlerConstructor(
				app.mm,
				app.configurator,
				app.appKeepers(),
			),
		)
	}
}

func (app *Tabi) appKeepers() upgrades.AppKeepers {
	return upgrades.AppKeepers{
		AppCodec:        app.AppCodec(),
		BankKeeper:      app.BankKeeper,
		AccountKeeper:   app.AccountKeeper,
		GetKey:          app.GetKey,
		ModuleManager:   app.mm,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		ReaderWriter:    app,
	}
}

// registerUpgradeHandler implements the upgrade execution logic of the upgrade module
func (app *Tabi) registerUpgradeHandler(
	planName string,
	upgrades *storetypes.StoreUpgrades,
	upgradeHandler upgradetypes.UpgradeHandler,
) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		app.Logger().Info("not found upgrade plan", "planName", planName, "err", err.Error())
		return
	}

	if upgradeInfo.Name == planName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// this configures a no-op upgrade handler for the planName upgrade
		app.UpgradeKeeper.SetUpgradeHandler(planName, upgradeHandler)
		// configure store loader that checks if version+1 == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrades))
	}
}
