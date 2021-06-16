package wos

import (
	"fmt"
	"strings"
)

type extensionOptions interface{}
type extensionHeaders func(headers map[string][]string, isWos bool) error

func setHeaderPrefix(key string, value string) extensionHeaders {
	return func(headers map[string][]string, isWos bool) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("set header %s with empty value", key)
		}
		setHeaders(headers, key, []string{value}, isWos)
		return nil
	}
}

// WithReqPaymentHeader sets header for requester-pays
func WithReqPaymentHeader(requester PayerType) extensionHeaders {
	return setHeaderPrefix(REQUEST_PAYER, string(requester))
}
