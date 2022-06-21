package watcher

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/suite"
)

type WatcherTestSuite struct {
	suite.Suite
	watcher *Watcher
}

func (suite *WatcherTestSuite) SetupTest() {
	suite.watcher = NewWatcher(nil)
}

func TestWatcherSuite(t *testing.T) {
	suite.Run(t, new(WatcherTestSuite))
}

func (suite *WatcherTestSuite) TestCreateWatchTx() {
	suite.SetupTest()
	wtx := suite.watcher.createWatchTx(auth.NewStdTx(nil, auth.NewStdFee(0, nil), nil, ""))
	if wtx == nil {
		suite.T().Log("----")
	}
	//	suite.Equal(nil, wtx)
	//	var txMsg WatchTx
	//	suite.Equal(nil, txMsg)

	eTx := NewEvmTx(nil, common.Hash{}, common.Hash{}, 0, 0)
	tx := eTx.GetTxWatchMessage()
	suite.Equal(nil, tx)
	if txWatchMessage := eTx.GetTxWatchMessage(); txWatchMessage != nil {
		suite.T().Log("-------")
	}
	suite.watcher.saveTx(eTx)

	batch := []WatchMessage{}
	var evTx *DelAccMsg
	batch = append(batch, evTx)
	suite.watcher.commitBatch(batch)
}
