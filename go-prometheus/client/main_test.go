package client

import (
	"context"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/promql/parser"
)

var prometheusUrl = "http://10.0.0.195:9002"

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
	query, err := client.Query(context.TODO(), "is_alive{category=\"database.gaussdb\"}", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(query)

}

func TestQueryRange(t *testing.T) {
	active := new(bool)
	t.Logf("%v", *active)
}

func TestValidate(t *testing.T) {
	// 要校验的 PromQL 查询
	query := "111"
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
	expr, err := parser.ParseExpr("db{}")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(expr)

}

func TestPrettify(t *testing.T) {
	expr, err := parser.ParseExpr("a @ 1")
	if err != nil {
		slog.Error("ParseExpr is failed", slog.Any("err", err))
	}
	t.Log(parser.Prettify(expr))
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
