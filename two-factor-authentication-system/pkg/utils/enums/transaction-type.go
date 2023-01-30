package enums

import "errors"

type TransactionType string

const (
	TRADE       TransactionType = "TRADE"
	DISPUTE     TransactionType = "DISPUTE"
	COIN_REDEEM TransactionType = "COIN_REDEEM"
	SEND_COIN   TransactionType = "SEND_COIN"
)

func (transactionType TransactionType) IsValid() error {
	switch transactionType {
	case TRADE, DISPUTE, COIN_REDEEM, SEND_COIN:
		return nil
	}
	return errors.New("invalid Transaction Type")
}
