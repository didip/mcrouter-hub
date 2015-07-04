package models

import (
	"testing"
)

func mockdata() []byte {
	return []byte(`STAT version mcrouter 1.0
STAT commandargs -p 5000 -f /opt/mcrouter-config/mcrouter.json
STAT pid 17901
STAT parent_pid 65072
STAT time 1435350338
STAT uptime 46
STAT num_servers 9
STAT num_servers_new 9
STAT num_servers_up 0
STAT num_servers_down 0
STAT num_servers_closed 0
STAT num_clients 1
STAT num_suspect_servers 0
STAT mcc_txbuf_reqs 0
STAT mcc_waiting_replies 0
STAT destination_batch_size 0
STAT asynclog_requests 0
STAT proxy_reqs_processing 1
STAT proxy_reqs_waiting 0
STAT client_queue_notify_period 0
STAT rusage_system 0.022264
STAT rusage_user 0.034848
STAT ps_num_minor_faults 3839
STAT ps_num_major_faults 0
STAT ps_user_time_sec 0.03
STAT ps_system_time_sec 0.02
STAT ps_vsize 581939200
STAT ps_rss 15605760
STAT fibers_allocated 0
STAT fibers_pool_size 0
STAT fibers_stack_high_watermark 0
STAT successful_client_connections 1
STAT cycles_avg 0
STAT cycles_min 0
STAT cycles_max 0
STAT cycles_p01 0
STAT cycles_p05 0
STAT cycles_p50 0
STAT cycles_p95 0
STAT cycles_p99 0
STAT cycles_num 0
STAT duration_us 0
STAT cmd_cas_count 0
STAT cmd_delete_count 0
STAT cmd_get_count 0
STAT cmd_gets_count 0
STAT cmd_set_count 0
STAT cmd_cas_outlier_count 0
STAT cmd_cas_outlier_failover_count 0
STAT cmd_cas_outlier_shadow_count 0
STAT cmd_cas_outlier_all_count 0
STAT cmd_delete_outlier_count 0
STAT cmd_delete_outlier_failover_count 0
STAT cmd_delete_outlier_shadow_count 0
STAT cmd_delete_outlier_all_count 0
STAT cmd_get_outlier_count 0
STAT cmd_get_outlier_failover_count 0
STAT cmd_get_outlier_shadow_count 0
STAT cmd_get_outlier_all_count 0
STAT cmd_gets_outlier_count 0
STAT cmd_gets_outlier_failover_count 0
STAT cmd_gets_outlier_shadow_count 0
STAT cmd_gets_outlier_all_count 0
STAT cmd_set_outlier_count 0
STAT cmd_set_outlier_failover_count 0
STAT cmd_set_outlier_shadow_count 0
STAT cmd_set_outlier_all_count 0
STAT cmd_other_outlier_count 0
STAT cmd_other_outlier_failover_count 0
STAT cmd_other_outlier_shadow_count 0
STAT cmd_other_outlier_all_count 0
STAT cmd_cas 0
STAT cmd_delete 0
STAT cmd_get 0
STAT cmd_gets 0
STAT cmd_set 0
STAT cmd_cas_outlier 0
STAT cmd_cas_outlier_failover 0
STAT cmd_cas_outlier_shadow 0
STAT cmd_cas_outlier_all 0
STAT cmd_delete_outlier 0
STAT cmd_delete_outlier_failover 0
STAT cmd_delete_outlier_shadow 0
STAT cmd_delete_outlier_all 0
STAT cmd_get_outlier 0
STAT cmd_get_outlier_failover 0
STAT cmd_get_outlier_shadow 0
STAT cmd_get_outlier_all 0
STAT cmd_gets_outlier 0
STAT cmd_gets_outlier_failover 0
STAT cmd_gets_outlier_shadow 0
STAT cmd_gets_outlier_all 0
STAT cmd_set_outlier 0
STAT cmd_set_outlier_failover 0
STAT cmd_set_outlier_shadow 0
STAT cmd_set_outlier_all 0
STAT cmd_other_outlier 0
STAT cmd_other_outlier_failover 0
STAT cmd_other_outlier_shadow 0
STAT cmd_other_outlier_all 0`)
}

func TestStatsConstructor(t *testing.T) {
	stats := NewStats(mockdata())
	if stats.Version != "mcrouter 1.0" {
		t.Error("Failed to assign correctly. Key: %v, Value: %v", "stats.Version", stats.Version)
	}
}

func TestGetStats(t *testing.T) {
	statsManager := NewMcRouterStatsManager("")
	stats, err := statsManager.Stats()
	if err != nil {
		t.Fatalf("Failed to get stats. McRouter may not be up. You must run McRouter on localhost 5000. Error: %v", err)
	}
	if stats == nil {
		t.Errorf("Failed to get stats.")
	}
	if stats.Version != "mcrouter 1.0" {
		t.Errorf("Failed to get stats values correctly.")
	}
}
