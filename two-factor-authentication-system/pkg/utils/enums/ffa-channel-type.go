package enums

import "errors"

type TfaChannelType string

const (
	SMS TfaChannelType = "SMS"
	APP TfaChannelType = "APP"
)

func (tfaChannelType TfaChannelType) IsValid() error {
	switch tfaChannelType {
	case SMS, APP:
		return nil
	}
	return errors.New("Invalid leave type")
}
