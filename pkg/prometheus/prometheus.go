package prometheus

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/ahilsend/vaultify/pkg"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	buildInfo = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_build_info",
			Help: "Build information like version and commit",
		},
		[]string{"version", "goversion", "commit_hash"},
	)
	authLeaseRenewed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vaultify_auth_lease_renewed",
			Help: "Counter for renewed auth leases",
		},
		[]string{"role", "warnings"},
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
		[]string{"role", "secret", "warnings"},
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
	prometheus.MustRegister(buildInfo)
	prometheus.MustRegister(authLeaseRenewed)
	prometheus.MustRegister(authLeaseFailed)
	prometheus.MustRegister(secretLeaseRenewed)
	prometheus.MustRegister(secretLeaseFailed)

	buildInfo.With(prometheus.Labels{
		"version":     pkg.Version,
		"commit_hash": pkg.CommitHash,
		"goversion":   runtime.Version(),
	}).Inc()
}

func IncAuthLeaseRenewed(role string, hasWarnings bool) {
	authLeaseRenewed.With(prometheus.Labels{
		"role":     role,
		"warnings": strconv.FormatBool(hasWarnings),
	}).Inc()
}

func IncAuthLeaseFailed(role string) {
	authLeaseFailed.With(prometheus.Labels{
		"role": role,
	}).Inc()
}

func IncSecretLeaseRenewed(role string, secret string, hasWarnings bool) {
	secretLeaseRenewed.With(prometheus.Labels{
		"role":     role,
		"secret":   secret,
		"warnings": strconv.FormatBool(hasWarnings),
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
