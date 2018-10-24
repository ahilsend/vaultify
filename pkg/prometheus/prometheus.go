package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	authLeaseRenewed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_auth_lease_renewed",
			Help: "Counter for renewed auth leases",
		},
		[]string{"role"},
	)
	authLeaseFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_auth_lease_renewal_failed",
			Help: "Counter for renewed auth leases",
		},
		[]string{"role"},
	)
	secretLeaseRenewed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_secret_lease_renewed",
			Help: "Counter for renewed secrets leases",
		},
		[]string{"role", "secret"},
	)
	secretLeaseFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_secret_lease_renewal_failed",
			Help: "Counter for renewed secrets leases",
		},
		[]string{"role", "secret"},
	)
)

func init() {
	// Register the metrics
	prometheus.MustRegister(authLeaseRenewed)
	prometheus.MustRegister(authLeaseFailed)
	prometheus.MustRegister(secretLeaseRenewed)
	prometheus.MustRegister(secretLeaseFailed)
}

func IncAuthLeaseRenewed(role string) {
	authLeaseRenewed.With(prometheus.Labels{
		"role": role,
	}).Inc()
}

func IncAuthLeaseFailed(role string) {
	authLeaseFailed.With(prometheus.Labels{
		"role": role,
	}).Inc()
}

func IncSecretLeaseRenewed(role string, secret string) {
	secretLeaseRenewed.With(prometheus.Labels{
		"role":   role,
		"secret": secret,
	}).Inc()
}

func IncSecretLeaseFailed(role string, secret string) {
	secretLeaseFailed.With(prometheus.Labels{
		"role":   role,
		"secret": secret,
	}).Inc()
}

func StartServer(metricsAddr string, metricsPath string) error {

	http.Handle(metricsPath, promhttp.Handler())
	return http.ListenAndServe(metricsAddr, nil)
}
