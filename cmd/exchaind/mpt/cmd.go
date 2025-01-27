package mpt

import (
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt/types"
	"github.com/spf13/cobra"
)

func MptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpt",
		Short: "migrate iavl state to mpt state (if use migrate mpt data, then you should set `--use-composite-key true` when you decide to use mpt to store the coming data)",
	}

	cmd.AddCommand(
		iavl2mptCmd(ctx),
		cleanIavlStoreCmd(ctx),
		mptViewerCmd(ctx),
	)
	cmd.PersistentFlags().UintVar(&types.TrieRocksdbBatchSize, types.FlagTrieRocksdbBatchSize, 100, "Concurrent rocksdb batch size for mpt")

	return cmd
}
