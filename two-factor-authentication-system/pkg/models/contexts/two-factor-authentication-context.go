package contexts

import "github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"

type TwoFactorAuthenticationContext struct {
	UserId            int64                `json:"userId"`
	SecretKey         string               `json:"secretKey"`
	TfaChannelType    enums.TfaChannelType `json:"tfaChannelType"`
	OtpGenerationTime int64                `json:"otpGenerationTime"`
}
