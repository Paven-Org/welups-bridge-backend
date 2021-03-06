package wel

import (
	"bridge/service-managers/logger"
	"crypto/sha256"
	"fmt"
	"math/big"

	welclient "github.com/Paven-Org/gotron-sdk/pkg/client"
	"github.com/Paven-Org/gotron-sdk/pkg/proto/api"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
)

// Contract caller for the ERC20 interface

type WelExport struct {
	address string
	cli     *welclient.GrpcClient
}

func MkWelExport(welcli *welclient.GrpcClient, contract string) *WelExport {
	return &WelExport{cli: welcli, address: contract}
}

func (ex *WelExport) TriggerConstantContract(from, contractAddress, method, jsonString string) (*api.TransactionExtention, error) {
	return ex.cli.TriggerConstantContract(from, contractAddress, method, jsonString)
}

func (ex *WelExport) TriggerContract(from, contractAddress, method, jsonString string, feeLimit, tAmount int64, tTokenID string, tTokenAmount int64) (*api.TransactionExtention, error) {
	return ex.cli.TriggerContract(from, contractAddress, method, jsonString, feeLimit, tAmount, tTokenID, tTokenAmount)
}

func (ex *WelExport) Withdraw(opts *CallOpts, tokenAddr string, account string, networkID *big.Int, value int64) (*api.TransactionExtention, error) {
	jsonString := fmt.Sprintf(`[{"address": "%s"},{"address":"%s"},{"uint256": "%s"},{"uint256":"%d"}]`, tokenAddr, account, networkID.String(), value)
	res, err := ex.cli.TriggerContract(opts.From, ex.address, "withdraw(address,address,uint256,uint256)", jsonString, opts.Fee_limit, opts.T_amount, "", opts.T_amount)
	if err != nil {
		logger.Get().Err(err).Msgf("RPC to make contract call failed")
		return nil, err
	}

	// signing
	logger.Get().Info().Msgf("Signing transaction...")

	tx := res.Transaction
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to sign transaction")
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, opts.Prikey)
	if err != nil {
		return nil, err
	}
	tx.Signature = append(tx.Signature, signature)

	// Broadcast
	ret, err := ex.cli.Broadcast(tx)
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	fmt.Println(ret)
	if ret.GetCode() != api.Return_SUCCESS {
		err = fmt.Errorf(api.ReturnResponseCode_name[int32(ret.GetCode())])
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	// Optionally loop to confirm transaction right here
	// it's better to do this inside of a workflow anyway, so kinda not compelling

	// return
	res.Transaction = tx
	res.Result = ret

	return res, nil
}

func (ex *WelExport) Claim(opts *CallOpts, tokenAddr string, account string, requestID *big.Int, amount *big.Int, sig []byte) (*api.TransactionExtention, error) {
	jsonString := fmt.Sprintf(`[{"address": "%s"},{"address":"%s"},{"uint256": "%s"},{"uint256":"%s"},{"bytes":"%x"}]`, tokenAddr, account, requestID.String(), amount.String(), sig)
	res, err := ex.cli.TriggerContract(opts.From, ex.address, "claim(address,address,uint256,uint256,bytes)", jsonString, opts.Fee_limit, opts.T_amount, "", amount.Int64())
	if err != nil {
		logger.Get().Err(err).Msgf("RPC to make contract call failed")
		return nil, err
	}

	// signing
	logger.Get().Info().Msgf("Signing transaction...")

	tx := res.Transaction
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to sign transaction")
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, opts.Prikey)
	if err != nil {
		return nil, err
	}
	tx.Signature = append(tx.Signature, signature)

	// Broadcast
	ret, err := ex.cli.Broadcast(tx)
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	fmt.Println(ret)
	if ret.GetCode() != api.Return_SUCCESS {
		err = fmt.Errorf(api.ReturnResponseCode_name[int32(ret.GetCode())])
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	// Optionally loop to confirm transaction right here
	// it's better to do this inside of a workflow anyway, so kinda not compelling

	// return
	res.Transaction = tx
	res.Result = ret

	return res, nil
}
