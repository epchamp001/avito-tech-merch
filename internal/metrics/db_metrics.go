package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// DBQueryDuration - гистограмма времени выполнения запросов к базе данных
	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query"},
	)

	// DBActiveConnections - gauge для количества активных соединений с базой данных
	DBActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_active_connections",
			Help: "Number of active database connections.",
		},
	)

	// DBErrorsTotal - счетчик количества ошибок при выполнении запросов к базе данных
	DBErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors.",
		},
		[]string{"query"},
	)
)

func init() {
	prometheus.MustRegister(DBQueryDuration, DBActiveConnections, DBErrorsTotal)
}

func RecordDBQueryDuration(query string, duration float64) {
	DBQueryDuration.WithLabelValues(query).Observe(duration)
}

func RecordDBActiveConnections(count int) {
	DBActiveConnections.Set(float64(count))
}

func RecordDBError(query string) {
	DBErrorsTotal.WithLabelValues(query).Inc()
}
