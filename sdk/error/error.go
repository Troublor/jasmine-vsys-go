package error

import (
	"fmt"
)

type Err struct {
	Msg string
}

func NewError(msg string) Err {
	return Err{Msg: msg}
}

func (e Err) Error() string {
	return fmt.Sprintf("SDK error: %s", e.Msg)
}

var InvalidCallerTxFailure = NewTransactionFailureErr("ContractInvalidCaller")

type TransactionFailureErr struct {
	Err
	Reason string
}

func NewTransactionFailureErr(reason string) TransactionFailureErr {
	return TransactionFailureErr{
		Err:    Err{Msg: reason},
		Reason: reason,
	}
}

func (e TransactionFailureErr) Error() string {
	return fmt.Sprintf("Transaction failure: %s", e.Reason)
}

var (
	NotFoundErr    = NewError("not found")
	UnconfirmedErr = NewError("do not have enough confirmation")
)
