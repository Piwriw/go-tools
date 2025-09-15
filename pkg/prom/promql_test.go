package prom

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/prometheus/common/model"

	"github.com/prometheus/client_golang/prometheus"
)

var prometheusUrl = "http://10.0.0.163:9002"

func TestValidateMetricName(t *testing.T) {
	metricName := "a-1"
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	valid := client.ValidateMetric(metricName)
	t.Log(valid)
}

func TestQueryMetric(t *testing.T) {
	metricName := "ALERTS"
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	metrics, err := client.QueryMetric(context.TODO(), metricName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(metrics)
}

func TestQueryAllMetrics(t *testing.T) {
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	metrics, err := client.QueryAllMetrics(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(metrics.Len())

}

func TestQuery(t *testing.T) {
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	query, err := client.Query(context.TODO(), "db_top_activity", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	resultMap, err := convertQueryResultToFieldValueMap(query)

	t.Log(resultMap)

}

// 将 Prometheus 查询结果转换为字段和值的映射
func convertQueryResultToFieldValueMap(query model.Value) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	switch v := query.(type) {
	case model.Vector: // 瞬时向量结果
		for _, sample := range v {
			entry := make(map[string]interface{})
			// 添加标签字段和值
			for key, value := range sample.Metric {
				entry[string(key)] = string(value)
			}
			// 添加值和时间戳
			entry["value"] = float64(sample.Value)
			entry["timestamp"] = sample.Timestamp.Time()
			result = append(result, entry)
		}
	case model.Matrix: // 时间序列结果
		for _, series := range v {
			for _, point := range series.Values {
				entry := make(map[string]interface{})
				// 添加标签字段和值
				for key, value := range series.Metric {
					entry[string(key)] = string(value)
				}
				// 添加值和时间戳
				entry["value"] = float64(point.Value)
				entry["timestamp"] = point.Timestamp.Time()
				result = append(result, entry)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported query result type: %T", query)
	}

	return result, nil
}
func TestQueryRange(t *testing.T) {
	active := new(bool)
	t.Logf("%v", *active)
}

func TestValidate(t *testing.T) {
	// 要校验的 PromQL 查询
	query := "123456"
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Validate(context.TODO(), query); err != nil {
		t.Fatal(err)
	}
	t.Log("PASS")
}

func TestName(t *testing.T) {
	if err := InitPromClient("http://10.0.0.173:9002", WithToken("QWdpbGVYQWRtaW46WVFBZG0xblBAc3MuKg==")); err != nil {
		t.Fatal(err)
	}
	if err := DefaultClient.Reload(context.TODO()); err != nil {
		t.Fatal(err)
	}

}

func TestPrettify(t *testing.T) {
	client, err := NewPrometheusClient(prometheusUrl)
	if err != nil {
		t.Fatal(err)
	}
	prettify, err := client.Prettify("1234567")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prettify)
	fmt.Println(fmt.Sprint("1234567"))
}

func TestPush(t *testing.T) {
	client, err := NewPrometheusClient(prometheusUrl)
	// 创建一个计数器指标
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "example_pushgateway_counter", // 指标名称
		Help: "This is an example counter pushed to Pushgateway",
	})

	// 增加计数器值
	counter.Inc()

	// 定义推送任务的名称和标签
	jobName := "example_job"
	instance := "localhost"
	kv := map[string]string{"job_name": "job1", "instance": instance}
	// 推送到 Pushgateway
	err = client.PushGateway("http://localhost:9091", jobName, kv, counter)

	if err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}

	log.Println("Metrics successfully pushed to Pushgateway")

	// 模拟持续推送的情况，每隔 5 秒推送一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 增加计数器
		counter.Inc()

		// 再次推送到 Pushgateway
		err = client.PushGateway("http://localhost:9091", jobName, kv, counter)
		if err != nil {
			log.Fatalf("Could not push to Pushgateway: %v", err)
		}

		log.Println("Pushed updated metrics to Pushgateway")
	}
}
