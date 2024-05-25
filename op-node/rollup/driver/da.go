package driver

import (
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	nubit "github.com/ethereum-optimism/optimism/op-nubit"
)

func SetDABackend(cfg nubit.CLIConfig) error {
	return derive.SetDABackendSingleton(cfg)
}
