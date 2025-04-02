package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/rulefmt"
)

func TestAddAlertRule(t *testing.T) {
	// 示例：创建一个高 CPU 使用率的告警规则
	// cpuRule := CreateAlertingRule(
	// 	"HighCPUUsage",
	// 	"node_cpu_usage_seconds_total{mode=\"idle\"} < 10",
	// 	map[string]string{
	// 		"severity": "critical",
	// 	},
	// 	map[string]string{
	// 		"summary":     "High CPU usage on {{ $labels.instance }}",
	// 		"description": "CPU usage is above 90% for more than 5 minutes on {{ $labels.instance }}",
	// 	},
	// 	model.Duration(5*time.Minute), // 5 minutes
	// )
	var rules []rulefmt.RuleNode
	for i := range 5000 {
		rules = append(rules, CreateAlertingRule(fmt.Sprintf("TESTJoohwanHPA%d", i),
			" db_stats_counter{category=\"database.mysql\",instance=\"instance-7\",Variable_name=\"Threads_connected\"}\n      >= 3",
			map[string]string{
				"alertConfigID": "10058",
				"category":      "database.mysql",
				"groupId":       "0",
				"instance":      "instance-7",
			},
			map[string]string{
				"description": "'数据库{{ $labels.name }} 连接数已达: {{$value}}'",
				"name":        "'{{ $labels.name }}'",
				"suggestion":  "可清理释放部分非活跃连接",
				"summary":     "'数据库{{ $labels.name }} 连接数已达: {{$value}}'",
				"value":       "'{{$value}}'",
			},
			model.Duration(0*time.Second), // 5 minutes
		))
	}
	// 创建一个规则组
	ruleGroup := GenerateRuleGroup("node-rules", 0, rules...)
	// 创建规则文件
	if err := AddAlertRule("/Users/joohwan/GolandProjects/go-tools/go-alertmanager/client/rules.yml", ruleGroup); err != nil {
		t.Fatal(err)
	}

}

func TestGetAlertGroups(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	alertGroups, err := client.Silenced(false).Inhibited(false).GetAlertGroups()
	if err != nil {
		t.Error(err)
	}
	for _, group := range alertGroups {
		t.Logf("%+v\n", group)
	}
}

func TestGetReceivers(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	receivers, err := client.GetReceivers()
	if err != nil {
		t.Fatal(err)
	}
	for _, receiver := range receivers {
		t.Log(*receiver.Name)
	}
}

func TestGetAlert(t *testing.T) {
	client := NewClient("10.0.0.195:9003")
	alerts, err := client.Silenced(false).
		Filter("severity=\"致命:#4A4A4A\"\n").Inhibited(false).GetAlerts()
	if err != nil {
		t.Error(err)
	}
	for _, alert := range alerts {
		t.Log(alert)
	}
}

func TestGetAlerts(t *testing.T) {
	// name, err2 := ExtractServerName("http://notice:9003")
	// if err2 != nil {
	//	t.Error(err2)
	// }
	// t.Log(name)
	// false, false, 0, , "severity=\"致命:#4A4A4A\""
	client := NewClient("10.0.0.195:9003")
	_, err := client.Silenced(false).
		Filter("instance=~\"instance-75|instance-56|instance-58\"").Inhibited(false).GetAlerts()
	if err != nil {
		t.Error(err)
	}

	go func(managerClient *AlertManagerClient) {
		for i := 0; i < 10; i++ {
			alerts, err := managerClient.Inhibited(false).GetAlerts()
			if err != nil {
				t.Error(err)
			}

			t.Log("All alerts len", len(alerts))
		}
	}(client)

	for i := 0; i < 10; i++ {
		alerts, err := client.Filter("instance=\"instance-46\"").GetAlerts()
		if err != nil {
			t.Error(err)
		}

		t.Log("All alerts len", len(alerts))

	}
	time.Sleep(1 * time.Minute)

}
