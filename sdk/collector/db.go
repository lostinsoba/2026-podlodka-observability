package collector

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type DBMetricCollector struct {
	db *sql.DB

	maxOpenDesc           *prometheus.Desc
	openDesc              *prometheus.Desc
	inUseDesc             *prometheus.Desc
	idleDesc              *prometheus.Desc
	waitedForDesc         *prometheus.Desc
	blockedSecondsDesc    *prometheus.Desc
	closedMaxIdleDesc     *prometheus.Desc
	closedMaxLifetimeDesc *prometheus.Desc
	closedMaxIdleTimeDesc *prometheus.Desc
}

func RegisterDBMetricCollector(mr prometheus.Registerer, db *sql.DB) error {
	dbMetricCollector := &DBMetricCollector{
		db: db,
		maxOpenDesc: prometheus.NewDesc(
			"db_max_open",
			"Maximum number of open connections to the database.",
			nil,
			nil,
		),
		openDesc: prometheus.NewDesc(
			"db_open_connections",
			"The number of established connections both in use and idle.",
			nil,
			nil,
		),
		inUseDesc: prometheus.NewDesc(
			"db_in_use_connections",
			"The number of connections currently in use.",
			nil,
			nil,
		),
		idleDesc: prometheus.NewDesc(
			"db_idle_connections",
			"The number of idle connections.",
			nil,
			nil,
		),
		waitedForDesc: prometheus.NewDesc(
			"db_waited_for_connections_total",
			"The total number of connections waited for.",
			nil,
			nil,
		),
		blockedSecondsDesc: prometheus.NewDesc(
			"db_blocked_seconds",
			"The total time blocked waiting for a new connection.",
			nil,
			nil,
		),
		closedMaxIdleDesc: prometheus.NewDesc(
			"db_closed_max_idle_connections_total",
			"The total number of connections closed due to SetMaxIdleConns.",
			nil,
			nil,
		),
		closedMaxLifetimeDesc: prometheus.NewDesc(
			"db_closed_max_lifetime_connections_total",
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil,
			nil,
		),
		closedMaxIdleTimeDesc: prometheus.NewDesc(
			"db_closed_max_idle_time_connections_total",
			"The total number of connections closed due to SetConnMaxIdleTime.",
			nil,
			nil,
		),
	}
	return mr.Register(dbMetricCollector)
}

func (mc *DBMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mc.maxOpenDesc
	ch <- mc.openDesc
	ch <- mc.inUseDesc
	ch <- mc.idleDesc
	ch <- mc.waitedForDesc
	ch <- mc.blockedSecondsDesc
	ch <- mc.closedMaxIdleDesc
	ch <- mc.closedMaxLifetimeDesc
	ch <- mc.closedMaxIdleTimeDesc
}

func (mc *DBMetricCollector) Collect(ch chan<- prometheus.Metric) {
	stats := mc.db.Stats()

	ch <- prometheus.MustNewConstMetric(
		mc.maxOpenDesc,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.openDesc,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.inUseDesc,
		prometheus.GaugeValue,
		float64(stats.InUse),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.idleDesc,
		prometheus.GaugeValue,
		float64(stats.Idle),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.waitedForDesc,
		prometheus.CounterValue,
		float64(stats.WaitCount),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.blockedSecondsDesc,
		prometheus.CounterValue,
		stats.WaitDuration.Seconds(),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.closedMaxIdleDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.closedMaxLifetimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		mc.closedMaxIdleTimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleTimeClosed),
	)
}
