package infura

import (
	"errors"
	"fmt"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/infura/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newStreamEngine(cfg *types.Config, logger log.Logger) (types.IStreamEngine, error) {
	if cfg.MysqlUrl == "" {
		return nil, errors.New("infura.mysql-url is empty")
	}
	return newMySQLEngine(cfg.MysqlUrl, cfg.MysqlUser, cfg.MysqlPass, cfg.MysqlDB, logger)
}

type MySQLEngine struct {
	db     *gorm.DB
	logger log.Logger
}

func newMySQLEngine(url, user, pass, dbName string, l log.Logger) (types.IStreamEngine, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, url, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&types.TransactionReceipt{}, &types.TransactionLog{},
		&types.LogTopic{}, &types.Block{}, &types.Transaction{}, &types.ContractCode{})
	return &MySQLEngine{
		db:     db,
		logger: l,
	}, nil
}

func (e *MySQLEngine) Write(streamData types.IStreamData) bool {
	e.logger.Debug("Begin MySqlEngine write")
	data := streamData.ConvertEngineData()
	trx := e.db.Begin()
	// write TransactionReceipts
	for _, receipt := range data.TransactionReceipts {
		ret := trx.Create(receipt)
		if ret.Error != nil {
			return e.rollbackWithError(trx, ret.Error)
		}
	}

	// write Block
	ret := trx.Create(data.Block)
	if ret.Error != nil {
		return e.rollbackWithError(trx, ret.Error)
	}

	// write contract code
	for _, code := range data.ContractCodes {
		ret := trx.Create(code)
		if ret.Error != nil {
			return e.rollbackWithError(trx, ret.Error)
		}
	}

	trx.Commit()
	e.logger.Debug("End MySqlEngine write")
	return true
}

func (e *MySQLEngine) rollbackWithError(trx *gorm.DB, err error) bool {
	trx.Rollback()
	e.logger.Error(err.Error())
	return false
}
