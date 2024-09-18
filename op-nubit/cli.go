package nubit

import (
	"log"
	"time"

	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
)

const (
	// The URL of Nubit DA node
	NodeRPCFlagName = "da.node_rpc"
	// The auth token of Nubit DA node
	AuthTokenFlagName = "da.auth_token"
	// The namespace of running Layer 2
	NamespaceFlagName = "da.namespace"
	// Enable the backup functionality
	EnableETHBackupFlagName = "da.enable_eth_backup"

	// NamespaceSize is the size of the hex encoded namespace string
	NamespaceSize = 29 * 2

	// Default local deployed Nubit Node
	DefaultNodeRPC = "http://localhost:26658"

	// 30*time.Duration(l.RollupConfig.BlockTime)*time.Second
	DefaultFetchTimeout = time.Minute

	// 30*time.Duration(l.RollupConfig.BlockTime)*time.Second
	DefaultSubmitTimeout = time.Minute

	// Prefix
	NubitDataPrefix = 0xda
)

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    NodeRPCFlagName,
			Usage:   "dial address of the Nubit node rpc client; supports grpc, http, https, ws",
			Value:   DefaultNodeRPC,
			EnvVars: opservice.PrefixEnvVar(envPrefix, "DA_NODE_RPC"),
		},
		&cli.StringFlag{
			Name:    AuthTokenFlagName,
			Usage:   "authentication token of the Nubit node",
			EnvVars: opservice.PrefixEnvVar(envPrefix, "DA_AUTH_TOKEN"),
		},
		&cli.StringFlag{
			Name:    NamespaceFlagName,
			Usage:   "namespace stands for the running Layer 2",
			EnvVars: opservice.PrefixEnvVar(envPrefix, "DA_NAMESPACE"),
		},
		&cli.BoolFlag{
			Name:    EnableETHBackupFlagName,
			Usage:   "enable eth as the backup rollup data store, default is true",
			Value:   true,
			EnvVars: opservice.PrefixEnvVar(envPrefix, "DA_ENABLE_ETH_BACKUP"),
		},
	}
}

type CLIConfig struct {
	Rpc             string
	AuthToken       string
	Namespace       string
	EnableETHBackup bool
	FetchTimeout    time.Duration
	SubmitTimeout   time.Duration
}

func (c CLIConfig) Check() error {
	log.Println("nubit: checking", "namespace", c.Namespace)
	log.Println("nubit: checking", "rpc", c.Rpc)
	return nil
}

func NewCLIConfig() CLIConfig {
	return CLIConfig{
		Rpc:       DefaultNodeRPC,
		AuthToken: "",
		// "op-stack" = "6F702D737461636B"
		Namespace:       "0000000000000000000000000000000000000000006F702D737461636B",
		EnableETHBackup: true,
		FetchTimeout:    DefaultFetchTimeout,
		SubmitTimeout:   DefaultSubmitTimeout,
	}
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		Rpc:             ctx.String(NodeRPCFlagName),
		AuthToken:       ctx.String(AuthTokenFlagName),
		Namespace:       ctx.String(NamespaceFlagName),
		EnableETHBackup: ctx.Bool(EnableETHBackupFlagName),
		FetchTimeout:    DefaultFetchTimeout,
		SubmitTimeout:   DefaultSubmitTimeout,
	}
}
