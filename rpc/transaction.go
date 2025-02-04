package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

var (
	feltZero = new(felt.Felt).SetUint64(0)
	feltOne  = new(felt.Felt).SetUint64(1)
	feltTwo  = new(felt.Felt).SetUint64(2)
)

func adaptTransaction(t TXN) (Transaction, error) {
	txMarshalled, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TransactionType_Invoke:
		var tx InvokeTxnV1
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_Declare:
		switch {
		case t.Version.Equal(feltZero):
			var tx DeclareTxnV0
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case t.Version.Equal(feltOne):
			var tx DeclareTxnV1
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		case t.Version.Equal(feltTwo):
			var tx DeclareTxnV2
			json.Unmarshal(txMarshalled, &tx)
			return tx, nil
		}
	case TransactionType_Deploy:
		var tx DeployTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_DeployAccount:
		var tx DeployAccountTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	case TransactionType_L1Handler:
		var tx L1HandlerTxn
		json.Unmarshal(txMarshalled, &tx)
		return tx, nil
	}
	return nil, errors.New(fmt.Sprint("internal error with adaptTransaction() : unknown transaction type ", t.Type))

}

// TransactionByHash gets the details and status of a submitted transaction.
func (provider *Provider) TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, error) {
	// todo: update to return a custom Transaction type, then use adapt function
	var tx TXN
	if err := do(ctx, provider.c, "starknet_getTransactionByHash", &tx, hash); err != nil {
			return nil, tryUnwrapToRPCErr(err,ErrHashNotFound)	
}
	return adaptTransaction(tx)
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (provider *Provider) TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error) {
	var tx TXN
	if err := do(ctx, provider.c, "starknet_getTransactionByBlockIdAndIndex", &tx, blockID, index); err != nil {
		
		return nil,tryUnwrapToRPCErr(err,  ErrInvalidTxnIndex ,ErrBlockNotFound)

	}
	return adaptTransaction(tx)
}

// TxnReceipt gets the transaction receipt by the transaction hash.
func (provider *Provider) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (TransactionReceipt, error) {
	var receipt UnknownTransactionReceipt
	err := do(ctx, provider.c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err,ErrHashNotFound)
	}
	return receipt.TransactionReceipt, nil
}

// GetTransactionStatus gets the transaction status (possibly reflecting that the tx is still in the mempool, or dropped from it)
// Parameters:
// - ctx: the context.Context object for cancellation and timeouts.
// - transactionHash: the transaction hash as a felt
// Returns:
// - *GetTxnStatusResp: The transaction status
// - error, if one arose.
func (provider *Provider) GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResp, error) {
	var receipt TxnStatusResp
	err := do(ctx, provider.c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound)
	}
	return &receipt, nil
}
