package client

import (
	"context"
	"fmt"
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

const defaultFlushDeadline = 1 * time.Minute

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
func (p *PrometheusClient) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	result, warnings, err := p.client.Query(ctx, query, ts)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// QueryRange 执行 Prometheus 时间范围查询
func (p *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (model.Value, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	result, warnings, err := p.client.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// Validate 验证 PromQL 查询是否合法
func (p *PrometheusClient) Validate(ctx context.Context, querySQL string) error {

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
	_, err := engine.NewInstantQuery(ctx, storage, nil, querySQL, time.Now())
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

// QueryAllMetrics 查询所有指标
func (p *PrometheusClient) QueryAllMetrics(ctx context.Context) (model.LabelValues, error) {
	values, warnings, err := p.client.LabelValues(ctx, "__name__", nil, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return values, nil
}

// QueryMetric 查询指定指标
func (p *PrometheusClient) QueryMetric(ctx context.Context, name string) (model.Value, error) {
	values, warnings, err := p.client.Query(ctx, name, time.Time{})
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return values, nil
}

// 将 Prometheus 查询结果转换为字段和值的映射
func (p *PrometheusClient) convertResToFieldMap(query model.Value) ([]map[string]interface{}, error) {
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
