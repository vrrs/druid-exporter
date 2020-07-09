package collector

import (
	"net"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	healthURL      = "/status/health"
	segmentDataURL = "/druid/coordinator/v1/datasources?simple"
	tasksURL       = "/druid/indexer/v1/tasks"
	workersURL     = "/druid/indexer/v1/workers"
	supervisorURL  = "/druid/indexer/v1/supervisor?full"
)

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus         *prometheus.Desc
	DataSourceCount           *prometheus.Desc
	DruidTasks                *prometheus.Desc
	DruidSupervisors          *prometheus.Desc
	DruidSegmentCount         *prometheus.Desc
	DruidSegmentSize          *prometheus.Desc
	DruidSegmentReplicateSize *prometheus.Desc
}

// SegementInterface is the interface for parsing segments data
type SegementInterface []struct {
	Name       string `json:"name"`
	Properties struct {
		Tiers struct {
			DefaultTier struct {
				Size           int64 `json:"size"`
				ReplicatedSize int64 `json:"replicatedSize"`
				SegmentCount   int   `json:"segmentCount"`
			} `json:"_default_tier"`
		} `json:"tiers"`
		Segments struct {
			MaxTime        time.Time `json:"maxTime"`
			Size           int64     `json:"size"`
			MinTime        time.Time `json:"minTime"`
			Count          int       `json:"count"`
			ReplicatedSize int64     `json:"replicatedSize"`
		} `json:"segments"`
	} `json:"properties"`
}

// TasksInterface is the interface for parsing druid tasks data
type TasksInterface []struct {
	ID               string  `json:"id"`
	GroupID          string  `json:"groupId"`
	Type             string  `json:"type"`
	CreatedTime      string  `json:"createdTime"`
	StatusCode       string  `json:"statusCode"`
	Status           string  `json:"status"`
	RunnerStatusCode string  `json:"runnerStatusCode"`
	Duration         float64 `json:"duration"`
	DataSource       string  `json:"dataSource"`
}

type worker struct {
	Worker struct {
		Host    string `json:"host"`
		Version string `json:"version"`
		IP      string `json:"ip"`
	}
	RunningTasks []string `json:"runningTasks"`
}

func (w worker) podName() string {
	return strings.Split(w.Worker.IP, ".")[0]
}

// ReverseDNSLookup returns hostname if no results are found
func ReverseDNSLookup(host string, dnsCache *cache.Cache) string {
	if cacheValue, found := dnsCache.Get(host); found {
		return cacheValue.(string)
	}
	names, err := net.LookupAddr(host)
	if err != nil || len(names) == 0 {
		dnsCache.Set(host, host, cache.DefaultExpiration)
		return host
	}
	dnsResult := strings.Split(names[0], ".")[0]
	dnsCache.Set(host, dnsResult, cache.DefaultExpiration)
	return dnsResult
}
