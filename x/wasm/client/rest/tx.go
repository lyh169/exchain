package rest

import (
	"github.com/gorilla/mux"
	"github.com/okex/exchain/x/wasm/ioutils"
	"net/http"
	"strconv"

	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	wasmUtils "github.com/okex/exchain/x/wasm/client/utils"
	"github.com/okex/exchain/x/wasm/types"
)

func registerTxRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/code", storeCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/code/{codeId}", instantiateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/{contractAddr}", executeContractHandlerFn(cliCtx)).Methods("POST")
}

type storeCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type instantiateContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Label   string       `json:"label" yaml:"label"`
	Deposit sdk.Coins    `json:"deposit" yaml:"deposit"`
	Admin   string       `json:"admin,omitempty" yaml:"admin"`
	Msg     []byte       `json:"msg" yaml:"msg"`
}

type executeContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	ExecMsg []byte       `json:"exec_msg" yaml:"exec_msg"`
	Amount  sdk.Coins    `json:"coins" yaml:"coins"`
}

func storeCodeHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeCodeReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var err error
		wasm := req.WasmBytes

		// gzip the wasm file
		if ioutils.IsWasm(wasm) {
			wasm, err = ioutils.GzipIt(wasm)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if !ioutils.IsGzip(wasm) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.MsgStoreCode{
			Sender:       req.BaseReq.From,
			WASMByteCode: wasm,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wasmUtils.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func instantiateContractHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req instantiateContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		vars := mux.Vars(r)

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get the id of the code to instantiate
		codeID, err := strconv.ParseUint(vars["codeId"], 10, 64)
		if err != nil {
			return
		}

		msg := types.MsgInstantiateContract{
			Sender: req.BaseReq.From,
			CodeID: codeID,
			Label:  req.Label,
			Funds:  sdk.CoinsToCoinAdapters(req.Deposit),
			Msg:    req.Msg,
			Admin:  req.Admin,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wasmUtils.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func executeContractHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req executeContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars["contractAddr"]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgExecuteContract{
			Sender:   req.BaseReq.From,
			Contract: contractAddr,
			Msg:      req.ExecMsg,
			Funds:    sdk.CoinsToCoinAdapters(req.Amount),
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wasmUtils.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
