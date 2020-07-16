package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type VersionInfo struct {
	UpgradeInfo sdk.UpgradeConfig `json:"upgrade_info"`
	Success     bool              `json:"success"`
}

func NewVersionInfo(upgradeConfig sdk.UpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}
