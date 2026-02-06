package metrics

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// PgxPoolCollector implements prometheus.Collector for pgxpool stats.
type PgxPoolCollector struct {
	pool *pgxpool.Pool

	acquiredConns          *prometheus.Desc
	idleConns              *prometheus.Desc
	totalConns             *prometheus.Desc
	maxConns               *prometheus.Desc
	constructingConns      *prometheus.Desc
	acquireCountTotal      *prometheus.Desc
	emptyAcquireCountTotal *prometheus.Desc
	canceledAcquireTotal   *prometheus.Desc
	acquireDurationTotal   *prometheus.Desc
}

func newPgxPoolCollector(pool *pgxpool.Pool) *PgxPoolCollector {
	return &PgxPoolCollector{
		pool: pool,
		acquiredConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "acquired_conns"),
			"Number of currently acquired connections in pool",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "idle_conns"),
			"Number of currently idle connections in pool",
			nil, nil,
		),
		totalConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "total_conns"),
			"Total number of connections in pool",
			nil, nil,
		),
		maxConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "max_conns"),
			"Maximum number of connections allowed in pool",
			nil, nil,
		),
		constructingConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "constructing_conns"),
			"Number of connections currently being constructed",
			nil, nil,
		),
		acquireCountTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "acquire_count_total"),
			"Total number of successful connection acquires",
			nil, nil,
		),
		emptyAcquireCountTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "empty_acquire_count_total"),
			"Total number of acquires from an empty pool",
			nil, nil,
		),
		canceledAcquireTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "canceled_acquire_count_total"),
			"Total number of acquires canceled by context",
			nil, nil,
		),
		acquireDurationTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pgpool", "acquire_duration_seconds_total"),
			"Total cumulative time spent acquiring connections",
			nil, nil,
		),
	}
}

func (c *PgxPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquiredConns
	ch <- c.idleConns
	ch <- c.totalConns
	ch <- c.maxConns
	ch <- c.constructingConns
	ch <- c.acquireCountTotal
	ch <- c.emptyAcquireCountTotal
	ch <- c.canceledAcquireTotal
	ch <- c.acquireDurationTotal
}

func (c *PgxPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.pool.Stat()

	ch <- prometheus.MustNewConstMetric(c.acquiredConns, prometheus.GaugeValue, float64(stats.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stats.MaxConns()))
	ch <- prometheus.MustNewConstMetric(c.constructingConns, prometheus.GaugeValue, float64(stats.ConstructingConns()))
	ch <- prometheus.MustNewConstMetric(c.acquireCountTotal, prometheus.CounterValue, float64(stats.AcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.emptyAcquireCountTotal, prometheus.CounterValue, float64(stats.EmptyAcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.canceledAcquireTotal, prometheus.CounterValue, float64(stats.CanceledAcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.acquireDurationTotal, prometheus.CounterValue, stats.AcquireDuration().Seconds())
}

// RedisPoolCollector implements prometheus.Collector for redis pool stats.
type RedisPoolCollector struct {
	client *redis.Client

	totalConns *prometheus.Desc
	idleConns  *prometheus.Desc
	staleConns *prometheus.Desc
	hitsTotal  *prometheus.Desc
	missTotal  *prometheus.Desc
	timeouts   *prometheus.Desc
}

func newRedisPoolCollector(client *redis.Client) *RedisPoolCollector {
	return &RedisPoolCollector{
		client: client,
		totalConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "total_conns"),
			"Total number of connections in Redis pool",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "idle_conns"),
			"Number of idle connections in Redis pool",
			nil, nil,
		),
		staleConns: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "stale_conns"),
			"Number of stale connections removed from Redis pool",
			nil, nil,
		),
		hitsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "hits_total"),
			"Total number of times a free connection was found in Redis pool",
			nil, nil,
		),
		missTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "misses_total"),
			"Total number of times a free connection was not found in Redis pool",
			nil, nil,
		),
		timeouts: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "redis_pool", "timeouts_total"),
			"Total number of wait timeouts in Redis pool",
			nil, nil,
		),
	}
}

func (c *RedisPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalConns
	ch <- c.idleConns
	ch <- c.staleConns
	ch <- c.hitsTotal
	ch <- c.missTotal
	ch <- c.timeouts
}

func (c *RedisPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.client.PoolStats()

	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns))
	ch <- prometheus.MustNewConstMetric(c.staleConns, prometheus.GaugeValue, float64(stats.StaleConns))
	ch <- prometheus.MustNewConstMetric(c.hitsTotal, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.missTotal, prometheus.CounterValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.timeouts, prometheus.CounterValue, float64(stats.Timeouts))
}

// RegisterPoolCollectors registers pgx and redis pool collectors with prometheus.
func RegisterPoolCollectors(pool *pgxpool.Pool, client *redis.Client) {
	if pool != nil {
		prometheus.MustRegister(newPgxPoolCollector(pool))
	}
	if client != nil {
		prometheus.MustRegister(newRedisPoolCollector(client))
	}
}
