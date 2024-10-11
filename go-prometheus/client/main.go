package client

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/pkg/errors"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/teststorage"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PrometheusClient 封装 Prometheus API 客户端
type PrometheusClient struct {
	client v1.API
}

// NewPrometheusClient 初始化 PrometheusClient
func NewPrometheusClient(address string) (*PrometheusClient, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return &PrometheusClient{
		client: v1.NewAPI(client),
	}, nil
}

// Query 执行 Prometheus 查询
func (p *PrometheusClient) Query(query string, ts time.Time) (model.Value, error) {
	result, warnings, err := p.client.Query(context.Background(), query, ts)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// QueryRange 执行 Prometheus 时间范围查询
func (p *PrometheusClient) QueryRange(query string, start, end time.Time, step time.Duration) (model.Value, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	result, warnings, err := p.client.QueryRange(context.Background(), query, r)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// Validate 验证 PromQL 查询是否合法
func (p *PrometheusClient) Validate(querySQL string) error {

	// 创建查询引擎
	opts := promql.EngineOpts{
		MaxSamples:    50000000, // 设置查询的最大样本
		Timeout:       0,        // 超时时间
		Logger:        nil,      // 可以传递 logger
		LookbackDelta: 0,
	}

	engine := promql.NewEngine(opts)

	// 使用 Prometheus 的测试存储来执行查询（这里不真正评估数据，只做语法校验）
	storage := teststorage.New(nil)
	defer storage.Close()
	_, err := engine.NewInstantQuery(context.TODO(), storage, nil, querySQL, time.Now())
	if err != nil {
		return errors.Wrapf(err, "Valide query is failed")
	}
	return nil
}

// Prettify 格式化 PromQL 查询
func (p *PrometheusClient) Prettify(querySQL string) (string, error) {
	expr, err := parser.ParseExpr(querySQL)
	if err != nil {
		return "", err
	}
	return parser.Prettify(expr), nil
}

// PushGateway 推送指标到 PushGateway
func (p *PrometheusClient) PushGateway(pushGatewayUrl, jobName string, groups map[string]string, collect prometheus.Collector) error {
	pusher := push.New(pushGatewayUrl, jobName).
		Collector(collect)
	for k, v := range groups {
		pusher.Grouping(k, v)
	}
	return pusher.Push()
}
