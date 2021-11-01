package client

import "context"

//认证参数
type AuthParam struct {
	Username string
	Password string
}

func (a *AuthParam) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": a.Username,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *AuthParam) RequireTransportSecurity() bool {
	return false
}
