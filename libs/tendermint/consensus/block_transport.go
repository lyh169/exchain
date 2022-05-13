package consensus

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"time"
)

type BlockTransport struct {
	height int64
	recvProposal time.Time
	firstPart time.Time
	totalElapsed time.Duration
	first2LastPartElapsed time.Duration
	droppedDue2NotExpected int
	droppedDue2NotAdded int
	droppedDue2WrongHeight int
	totalParts int
	Logger  log.Logger
}

func (bt *BlockTransport) onProposal(height int64)  {
	if bt.height == height || bt.height == 0 {
		bt.recvProposal = time.Now()
		bt.height = height
	} else {
		//panic("invalid height")
	}
}

func (bt *BlockTransport) reset(height int64) {
	bt.height = height
	bt.droppedDue2NotExpected = 0
	bt.droppedDue2NotAdded = 0
	bt.droppedDue2WrongHeight = 0
	bt.totalParts = 0
}

func (bt *BlockTransport) on1stPart(height int64)  {
	if bt.height == height || bt.height == 0 {
		bt.firstPart = time.Now()
		bt.height = height
	} else {
		//panic("invalid height")
	}
}

func (bt *BlockTransport) onRecvBlock(height int64)  {
	if bt.height == height {
		bt.totalElapsed = time.Now().Sub(bt.recvProposal)
		bt.first2LastPartElapsed = time.Now().Sub(bt.firstPart)
		trace.GetElapsedInfo().AddInfo(trace.RecvBlock,
			fmt.Sprintf("%d<%dms>", height, bt.totalElapsed.Milliseconds()))
		trace.GetElapsedInfo().AddInfo(trace.First2LastPart,
			fmt.Sprintf("%d<%dms>", height, bt.first2LastPartElapsed.Milliseconds()))
	} else {
		//panic("invalid height")
	}
}
