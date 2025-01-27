package simulation_test

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/types/kv"
	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	// "github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	ibctmtypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp"
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	clientID := "clientidone"

	height := types.NewHeight(0, 10)

	clientState := &ibctmtypes.ClientState{
		FrozenHeight: height,
	}

	consState := &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   host.FullClientStateKey(clientID),
				Value: app.IBCKeeper.ClientKeeper.MustMarshalClientState(clientState),
			},
			{
				Key:   host.FullConsensusStateKey(clientID, height),
				Value: app.IBCKeeper.ClientKeeper.MustMarshalConsensusState(consState),
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"ClientState", fmt.Sprintf("ClientState A: %v\nClientState B: %v", clientState, clientState)},
		{"ConsensusState", fmt.Sprintf("ConsensusState A: %v\nConsensusState B: %v", consState, consState)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			// res, found := simulation.NewDecodeStore(app.IBCKeeper.ClientKeeper, kvPairs.Pairs[i], kvPairs.Pairs[i])
			kvA := tmkv.Pair{
				Key:   kvPairs.Pairs[i].GetKey(),
				Value: kvPairs.Pairs[i].GetValue(),
			}
			res, found := simulation.NewDecodeStore(app.IBCKeeper.ClientKeeper, kvA, kvA)
			if i == len(tests)-1 {
				require.False(t, found, string(kvPairs.Pairs[i].Key))
				require.Empty(t, res, string(kvPairs.Pairs[i].Key))
			} else {
				require.True(t, found, string(kvPairs.Pairs[i].Key))
				require.Equal(t, tt.expectedLog, res, string(kvPairs.Pairs[i].Key))
			}
		})
	}
}
