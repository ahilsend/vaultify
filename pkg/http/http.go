package http

import (
	ghttp "net/http"

	"github.com/ahilsend/vaultify/pkg/healthz"
)

func Serve(addr string) {
	healthz.Register()
	ghttp.ListenAndServe(addr, nil)
}
