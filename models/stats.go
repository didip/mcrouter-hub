package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

func NewMcRouterStatsManager(mcRouterAddr string) *McRouterStatsManager {
	if mcRouterAddr == "" {
		mcRouterAddr = "localhost:5000"
	}

	mc := &McRouterStatsManager{}
	mc.McRouterAddr = mcRouterAddr

	return mc
}

type McRouterStatsManager struct {
	McRouterAddr string
}

func (sm *McRouterStatsManager) Stats() (*Stats, error) {
	c1 := exec.Command("echo", "stats")
	c2 := exec.Command("nc", "localhost", "5000")

	var c2Output bytes.Buffer

	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = &c2Output

	c2.Start()

	err := c1.Run()
	if err != nil {
		return nil, err
	}

	err = c2.Wait()
	if err != nil {
		return nil, err
	}

	return NewStats(c2Output.Bytes()), nil
}

func (sm *McRouterStatsManager) StatsFromFile() (map[string]interface{}, error) {
	addrParts := strings.Split(sm.McRouterAddr, ":")

	statsFile := fmt.Sprintf("/var/mcrouter/stats/libmcrouter.mcrouter.%v.stats", addrParts[1])

	statsJson, err := ioutil.ReadFile(statsFile)
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	err = json.Unmarshal(statsJson, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func NewStats(data []byte) *Stats {
	stats := &Stats{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		command := parts[0]

		if command != "STAT" {
			continue
		}

		statsKey := parts[1]
		value := strings.Join(parts[2:], " ")

		if statsKey == "version" {
			stats.Version = value
		}
		if statsKey == "commandargs" {
			stats.CommandArgs = value
		}
		if statsKey == "pid" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.PID = valueInt64
			}
		}
		if statsKey == "parent_pid" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.ParentPID = valueInt64
			}
		}
		if statsKey == "time" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.Time = valueInt64
			}
		}
		if statsKey == "uptime" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.Uptime = valueInt64
			}
		}
		if statsKey == "num_servers" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumServers = valueInt64
			}
		}
		if statsKey == "num_servers_new" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumServersNew = valueInt64
			}
		}
		if statsKey == "num_servers_up" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumServersUp = valueInt64
			}
		}
		if statsKey == "num_servers_down" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumServersDown = valueInt64
			}
		}
		if statsKey == "num_servers_closed" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumServersClosed = valueInt64
			}
		}
		if statsKey == "num_clients" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumClients = valueInt64
			}
		}
		if statsKey == "num_suspect_servers" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.NumSuspectServers = valueInt64
			}
		}
		if statsKey == "mcc_txbuf_reqs" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.MccTxbufReqs = valueInt64
			}
		}
		if statsKey == "mcc_waiting_replies" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.MccWaitingReplies = valueInt64
			}
		}
		if statsKey == "destination_batch_size" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.DestinationBatchSize = valueInt64
			}
		}
		if statsKey == "asynclog_requests" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.AsynclogRequests = valueInt64
			}
		}
		if statsKey == "proxy_reqs_processing" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.ProxyReqsProcessing = valueInt64
			}
		}
		if statsKey == "proxy_reqs_waiting" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.ProxyReqsWaiting = valueInt64
			}
		}
		if statsKey == "client_queue_notify_period" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.ClientQueueNotifyPeriod = valueInt64
			}
		}
		if statsKey == "rusage_system" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats.RusageSystem = valueFloat64
			}
		}
		if statsKey == "rusage_user" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats.RusageUser = valueFloat64
			}
		}
		if statsKey == "ps_num_minor_faults" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.PsNumMinorFaults = valueInt64
			}
		}
		if statsKey == "ps_num_major_faults" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.PsNumMajorFaults = valueInt64
			}
		}
		if statsKey == "ps_user_time_sec" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats.PsUserTimeSec = valueFloat64
			}
		}
		if statsKey == "ps_system_time_sec" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats.PsSystemTimeSec = valueFloat64
			}
		}
		if statsKey == "ps_vsize" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.PsVsize = valueInt64
			}
		}
		if statsKey == "ps_rss" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.PsRss = valueInt64
			}
		}
		if statsKey == "fibers_allocated" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.FibersAllocated = valueInt64
			}
		}
		if statsKey == "fibers_pool_size" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.FibersPoolSize = valueInt64
			}
		}
		if statsKey == "fibers_stack_high_watermark" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.FibersStackHighWatermark = valueInt64
			}
		}
		if statsKey == "successful_client_connections" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats.SuccessfulClientConnections = valueInt64
			}
		}

	}

	return stats
}

type Stats struct {
	Version                     string
	CommandArgs                 string
	PID                         int64
	ParentPID                   int64
	Time                        int64
	Uptime                      int64
	NumServers                  int64
	NumServersNew               int64
	NumServersUp                int64
	NumServersDown              int64
	NumServersClosed            int64
	NumClients                  int64
	NumSuspectServers           int64
	MccTxbufReqs                int64
	MccWaitingReplies           int64
	DestinationBatchSize        int64
	AsynclogRequests            int64
	ProxyReqsProcessing         int64
	ProxyReqsWaiting            int64
	ClientQueueNotifyPeriod     int64
	RusageSystem                float64
	RusageUser                  float64
	PsNumMinorFaults            int64
	PsNumMajorFaults            int64
	PsUserTimeSec               float64
	PsSystemTimeSec             float64
	PsVsize                     int64
	PsRss                       int64
	FibersAllocated             int64
	FibersPoolSize              int64
	FibersStackHighWatermark    int64
	SuccessfulClientConnections int64
	// cycles_avg                        int64
	// cycles_min                        int64
	// cycles_max                        int64
	// cycles_p01                        int64
	// cycles_p05                        int64
	// cycles_p50                        int64
	// cycles_p95                        int64
	// cycles_p99                        int64
	// cycles_num                        int64
	// duration_us                       int64
	// cmd_cas_count                     int64
	// cmd_delete_count                  int64
	// cmd_get_count                     int64
	// cmd_gets_count                    int64
	// cmd_set_count                     int64
	// cmd_cas_outlier_count             int64
	// cmd_cas_outlier_failover_count    int64
	// cmd_cas_outlier_shadow_count      int64
	// cmd_cas_outlier_all_count         int64
	// cmd_delete_outlier_count          int64
	// cmd_delete_outlier_failover_count int64
	// cmd_delete_outlier_shadow_count   int64
	// cmd_delete_outlier_all_count      int64
	// cmd_get_outlier_count             int64
	// cmd_get_outlier_failover_count    int64
	// cmd_get_outlier_shadow_count      int64
	// cmd_get_outlier_all_count         int64
	// cmd_gets_outlier_count            int64
	// cmd_gets_outlier_failover_count   int64
	// cmd_gets_outlier_shadow_count     int64
	// cmd_gets_outlier_all_count        int64
	// cmd_set_outlier_count             int64
	// cmd_set_outlier_failover_count    int64
	// cmd_set_outlier_shadow_count      int64
	// cmd_set_outlier_all_count         int64
	// cmd_other_outlier_count           int64
	// cmd_other_outlier_failover_count  int64
	// cmd_other_outlier_shadow_count    int64
	// cmd_other_outlier_all_count       int64
	// cmd_cas                           int64
	// cmd_delete                        int64
	// cmd_get                           int64
	// cmd_gets                          int64
	// cmd_set                           int64
	// cmd_cas_outlier                   int64
	// cmd_cas_outlier_failover          int64
	// cmd_cas_outlier_shadow            int64
	// cmd_cas_outlier_all               int64
	// cmd_delete_outlier                int64
	// cmd_delete_outlier_failover       int64
	// cmd_delete_outlier_shadow         int64
	// cmd_delete_outlier_all            int64
	// cmd_get_outlier                   int64
	// cmd_get_outlier_failover          int64
	// cmd_get_outlier_shadow            int64
	// cmd_get_outlier_all               int64
	// cmd_gets_outlier                  int64
	// cmd_gets_outlier_failover         int64
	// cmd_gets_outlier_shadow           int64
	// cmd_gets_outlier_all              int64
	// cmd_set_outlier                   int64
	// cmd_set_outlier_failover          int64
	// cmd_set_outlier_shadow            int64
	// cmd_set_outlier_all               int64
	// cmd_other_outlier                 int64
	// cmd_other_outlier_failover        int64
	// cmd_other_outlier_shadow          int64
	// cmd_other_outlier_all             int64
}
