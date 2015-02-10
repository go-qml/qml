package types

// #include "networkCookie.h"
import "C"

import (
	"time"
)

// NetworkCookie corresponds to a QNetworkCookie
// It provides methods to be converted to/from QNetworkCookie
type NetworkCookie struct {
	Domain         string
	ExpirationDate time.Time
	HttpOnly       bool
	Secure         bool
	SessionCookie  bool
	Name           []byte
	Path           string
}
