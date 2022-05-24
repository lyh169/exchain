package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/go-amino"

	"github.com/okex/exchain/libs/tendermint/mempool"
)

var (
	cdc *amino.Codec
)

func init() {
	cdc = amino.NewCodec()
	mempool.RegisterMessages(cdc)
}

var hexKeys = []string{
	"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17", //0xbbE4733d85bc2b90682147779DA49caB38C0aA1F
	"171786c73f805d257ceb07206d851eea30b3b41a2170ae55e1225e0ad516ef42", //0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0
	"b7700998b973a2cae0cb8e8a328171399c043e57289735aca5f2419bd622297a", //0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64
	"00dcf944648491b3a822d40bf212f359f699ed0dd5ce5a60f1da5e1142855949", //0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044
}

func genPrivateKey(hexPrivateKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexPrivateKey)
}

func genAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

func genAddressHexPrivateKey(hexPrivateKey string) common.Address {
	privateKey, _ := genPrivateKey(hexPrivateKey)

	return genAddress(privateKey)
}

func genRandAddress() common.Address {
	privateKay, _ := crypto.GenerateKey()
	log.Printf("private key: %v \n", fmt.Sprintf("%x", crypto.FromECDSA(privateKay)))
	return crypto.PubkeyToAddress(privateKay.PublicKey)
}

func genTx(hexPrivateKey string, nonce uint64) []byte {
	fromPrivateKey, _ := genPrivateKey(hexPrivateKey)
	to := genRandAddress()
	unsignedTx := types.NewTransaction(nonce, to, big.NewInt(100), 3000000, big.NewInt(1), nil)
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(67)), fromPrivateKey)
	if err != nil {
		panic(err)
	}
	ret, err := signedTx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return ret
}

func writeTxMessage(w *bufio.Writer, tx []byte) error {
	msg := mempool.TxMessage{Tx: tx}
	if _, err := w.WriteString(hex.EncodeToString(cdc.MustMarshalBinaryBare(&msg))); err != nil {
		return err
	}
	return w.WriteByte('\n')
}

func write(hexPrivateKey string, nonce uint64, rule func(string, uint64) []byte) {
	f, err := os.OpenFile(fmt.Sprintf("../TxMessage-%s.txt", genAddressHexPrivateKey(hexPrivateKey)), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	txWriter := bufio.NewWriter(f)
	defer txWriter.Flush()
	writeTxMessage(txWriter, rule(hexPrivateKey, nonce))
}

func main() {
	//	for _, v := range hexKeys {
	//		write(v, 0, genTx)
	//	}
	for i := 0; i < 10; i++ {
		log.Println(genRandAddress())
	}
}
