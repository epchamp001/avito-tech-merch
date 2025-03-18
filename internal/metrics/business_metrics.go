package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// RegistrationsTotal - счетчик количества регистраций
	RegistrationsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "registrations_total",
			Help: "Total number of user registrations.",
		},
	)

	// LoginsTotal - счетчик количества авторизаций
	LoginsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logins_total",
			Help: "Total number of user logins.",
		},
	)

	// CoinTransfersTotal - счетчик количества операций по передаче монет
	CoinTransfersTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "coin_transfers_total",
			Help: "Total number of coin transfers.",
		},
	)

	// MerchPurchasesTotal - счетчик количества покупок мерча
	MerchPurchasesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "merch_purchases_total",
			Help: "Total number of merch purchases.",
		},
	)
)

func init() {
	prometheus.MustRegister(RegistrationsTotal, LoginsTotal, CoinTransfersTotal, MerchPurchasesTotal)
}

func RecordRegistration() {
	RegistrationsTotal.Inc()
}

func RecordLogin() {
	LoginsTotal.Inc()
}

func RecordCoinTransfer() {
	CoinTransfersTotal.Inc()
}

func RecordMerchPurchase() {
	MerchPurchasesTotal.Inc()
}
