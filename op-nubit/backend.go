package nubit

import (
	"encoding/hex"
	"time"

	"github.com/rollkit/go-da"
	"github.com/rollkit/go-da/proxy"
)

type NubitDABackend struct {
	Client       da.DA
	FetchTimeout time.Duration

	SubmitTimeout   time.Duration
	Namespace       da.Namespace
	EnableETHBackup bool
}

func NewNubitDABackend(rpc, token, namespace string, FetchTimeout, SubmitTimeout time.Duration, EnableETHBackup bool) (*NubitDABackend, error) {
	client, err := proxy.NewClient(rpc, token)
	if err != nil {
		return nil, err
	}
	ns, err := hex.DecodeString(namespace)
	if err != nil {
		return nil, err
	}
	return &NubitDABackend{
		Client:          client,
		FetchTimeout:    FetchTimeout,
		SubmitTimeout:   SubmitTimeout,
		Namespace:       ns,
		EnableETHBackup: EnableETHBackup,
	}, nil
}

func NewNubitDABackendFromCfg(c CLIConfig) (*NubitDABackend, error) {
	return NewNubitDABackend(c.Rpc, c.AuthToken, c.Namespace, c.FetchTimeout, c.SubmitTimeout, c.EnableETHBackup)
}
