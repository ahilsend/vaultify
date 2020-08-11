package healthz

import (
	"net/http"
)

func Register() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/readyz", readyz)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("healthy"))
}

func readyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ready"))
}
