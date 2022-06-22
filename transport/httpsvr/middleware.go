package httpsvr

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
