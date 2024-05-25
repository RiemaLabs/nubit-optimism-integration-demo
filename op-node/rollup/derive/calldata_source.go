package derive

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/eth"

	nubit "github.com/ethereum-optimism/optimism/op-nubit"
)

var daBackend *nubit.NubitDABackend

func SetDABackendSingleton(cfg nubit.CLIConfig) error {
	if daBackend != nil {
		// return errors.New("da client already configured")
		return nil
	}
	n, err := nubit.NewNubitDABackendFromCfg(cfg)
	if err != nil {
		return err
	}
	daBackend = n
	return nil
}

// CalldataSource is a fault tolerant approach to fetching data.
// The constructor will never fail & it will instead re-attempt the fetcher
// at a later point.
type CalldataSource struct {
	// Internal state + data
	open bool
	data []eth.Data
	// Required to re-attempt fetching
	ref     eth.L1BlockRef
	dsCfg   DataSourceConfig
	fetcher L1TransactionFetcher
	log     log.Logger

	batcherAddr common.Address
}

// NewCalldataSource creates a new calldata source. It suppresses errors in fetching the L1 block if they occur.
// If there is an error, it will attempt to fetch the result on the next call to `Next`.
func NewCalldataSource(ctx context.Context, log log.Logger, dsCfg DataSourceConfig, fetcher L1TransactionFetcher, ref eth.L1BlockRef, batcherAddr common.Address) (DataIter, error) {
	_, txs, err := fetcher.InfoAndTxsByHash(ctx, ref.Hash)
	if err != nil {
		return &CalldataSource{
			open:        false,
			ref:         ref,
			dsCfg:       dsCfg,
			fetcher:     fetcher,
			log:         log,
			batcherAddr: batcherAddr,
		}, nil
	}
	data, err := DataFromEVMTransactions(dsCfg, batcherAddr, txs, log.New("origin", ref))
	if err != nil {
		return &CalldataSource{
			open:        false,
			ref:         ref,
			dsCfg:       dsCfg,
			fetcher:     fetcher,
			log:         log,
			batcherAddr: batcherAddr,
		}, err
	}
	return &CalldataSource{
		open: true,
		data: data,
	}, nil
}

// Next returns the next piece of data if it has it. If the constructor failed, this
// will attempt to reinitialize itself. If it cannot find the block it returns a ResetError
// otherwise it returns a temporary error if fetching the block returns an error.
func (ds *CalldataSource) Next(ctx context.Context) (eth.Data, error) {
	if !ds.open {
		if _, txs, err := ds.fetcher.InfoAndTxsByHash(ctx, ds.ref.Hash); err == nil {
			ds.open = true
			ds.data, err = DataFromEVMTransactions(ds.dsCfg, ds.batcherAddr, txs, ds.log)
			// already wrapped
			if err != nil {
				return nil, err
			}
		} else if errors.Is(err, ethereum.NotFound) {
			return nil, NewResetError(fmt.Errorf("failed to open calldata source: %w", err))
		} else {
			return nil, NewTemporaryError(fmt.Errorf("failed to open calldata source: %w", err))
		}
	}
	if len(ds.data) == 0 {
		return nil, io.EOF
	} else {
		data := ds.data[0]
		ds.data = ds.data[1:]
		return data, nil
	}
}

// DataFromEVMTransactions filters all of the transactions and returns the calldata from transactions
// that are sent to the batch inbox address from the batch sender address.
// This will return an empty array if no valid transactions are found.
func DataFromEVMTransactions(dsCfg DataSourceConfig, batcherAddr common.Address, txs types.Transactions, log log.Logger) ([]eth.Data, error) {
	out := []eth.Data{}
	for _, tx := range txs {
		if isValidBatchTx(tx, dsCfg.l1Signer, dsCfg.batchInboxAddress, batcherAddr) {
			data := tx.Data()
			switch len(data) {
			case 0:
				// Corner case
				out = append(out, data)
			default:
				switch data[0] {
				case nubit.NubitDataPrefix:
					log.Info("nubit: blob request", "id", hex.EncodeToString(tx.Data()))
					ctx, cancel := context.WithTimeout(context.Background(), daBackend.FetchTimeout)
					blobs, err := daBackend.Client.Get(ctx, [][]byte{data[1:]}, daBackend.Namespace)
					cancel()
					if err != nil {
						return nil, NewResetError(fmt.Errorf("nubit: failed to resolve the prefix: %w", err))
					}
					if len(blobs) != 1 {
						log.Warn("nubit: unexpected length for blobs", "expected", 1, "got", len(blobs))
						if len(blobs) == 0 {
							log.Warn("nubit: skipping empty blobs")
							continue
						}
					}
					out = append(out, blobs[0])
				default:
					out = append(out, data)
					log.Info("nubit: submitting to eth as the backup")
				}
			}

		}
	}
	return out, nil
}
