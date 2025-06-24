package prom

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	promtheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	prommodel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/teststorage"
)

const (
	defaultFlushDeadline = 1 * time.Minute
)

var (
	defaultQueryTimeout = 15 * time.Second
)

var DefaultClient *PrometheusClient

func InitPromClient(address string, opts ...Option) error {
	var err error
	DefaultClient, err = NewPrometheusClient(address, opts...)
	if err != nil {
		return err
	}
	return nil
}

// PrometheusClient 封装 Prometheus API 客户端
type PrometheusClient struct {
	address    string
	token      string
	client     promtheusv1.API
	httpClient api.Client
	timeout    time.Duration
}

func (p *PrometheusClient) Client() promtheusv1.API {
	return p.client
}

func (p *PrometheusClient) HTTPClient() api.Client {
	return p.httpClient
}

// NewPrometheusClient 初始化 PrometheusClient，支持 Bearer Token 认证
func NewPrometheusClient(address string, opts ...Option) (*PrometheusClient, error) {
	pc := &PrometheusClient{
		address: address,
	}
	// 应用 Option 参数
	for _, opt := range opts {
		opt(pc)
	}

	// 构建 Transport，根据是否设置 token 决定是否加 Header
	transport := http.DefaultTransport
	if pc.token != "" {
		transport = &authTransport{
			base:  http.DefaultTransport,
			token: pc.token,
		}
	}

	// 创建 Prometheus 客户端
	client, err := api.NewClient(api.Config{
		Address:      address,
		RoundTripper: transport,
	})
	if err != nil {
		return nil, err
	}

	pc.client = promtheusv1.NewAPI(client)
	pc.httpClient = client
	return pc, nil
}

type Option func(*PrometheusClient)

func WithToken(token string) Option {
	return func(pc *PrometheusClient) {
		pc.token = token
	}
}

// WithTimeout 设置 Prometheus 查询超时时间
// 注意：该方法只对 Query 和 QueryRange 方法有效
func WithTimeout(timeout time.Duration) Option {
	return func(pc *PrometheusClient) {
		if timeout <= 0 {
			return
		}
		pc.timeout = timeout
	}
}

// authTransport 用于在请求中注入 Bearer Token
type authTransport struct {
	base  http.RoundTripper
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Basic "+t.token)
	return t.base.RoundTrip(req)
}

// Query 执行 Prometheus 查询
func (p *PrometheusClient) Query(ctx context.Context, query string, ts time.Time) (prommodel.Value, error) {
	// 如果传入的 ctx 没有超时，才设置默认超时
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultQueryTimeout)
		defer cancel()
	}
	result, warnings, err := p.client.Query(ctx, query, ts)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// QueryVector 执行 Prometheus 查询并返回 Vector 类型结果
func (p *PrometheusClient) QueryVector(ctx context.Context, query string, ts time.Time, opts ...promtheusv1.Option) (prommodel.Vector, error) {
	// 如果传入的 ctx 没有超时，才设置默认超时
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultQueryTimeout)
		defer cancel()
	}
	result, warnings, err := p.client.Query(ctx, query, ts, opts...)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}

	vectors, ok := result.(prommodel.Vector)
	if !ok {
		err = errors.New("result is not a vector")
		return nil, err
	}
	return vectors, nil
}

// QueryRange 执行 Prometheus 时间范围查询
func (p *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration, opts ...promtheusv1.Option) (prommodel.Value, error) {
	// 如果传入的 ctx 没有超时，才设置默认超时
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultQueryTimeout)
		defer cancel()
	}
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	result, warnings, err := p.client.QueryRange(ctx, query, r, opts...)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return result, nil
}

// QueryRangeMatrix 执行 Prometheus 时间范围查询
// 注意：该方法返回的是 prommodel.Matrix 类型，需要自行处理
func (p *PrometheusClient) QueryRangeMatrix(ctx context.Context, query string, start, end time.Time, step time.Duration, opts ...promtheusv1.Option) (prommodel.Matrix, error) {
	// 如果传入的 ctx 没有超时，才设置默认超时
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultQueryTimeout)
		defer cancel()
	}
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	result, warnings, err := p.client.QueryRange(ctx, query, r, opts...)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	matrixVal, ok := result.(prommodel.Matrix)
	if !ok {
		err = errors.New("result is not a matrix")
		return nil, err
	}
	return matrixVal, nil
}

// Validate 验证 PromQL 查询是否合法
func (p *PrometheusClient) Validate(ctx context.Context, querySQL string) error {
	// 创建查询引擎
	opts := promql.EngineOpts{
		MaxSamples: 50000000,             // 设置查询的最大样本
		Timeout:    defaultFlushDeadline, // 超时时间
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
func (p *PrometheusClient) QueryAllMetrics(ctx context.Context) (prommodel.LabelValues, error) {
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
func (p *PrometheusClient) QueryMetric(ctx context.Context, name string) (prommodel.Value, error) {
	values, warnings, err := p.client.Query(ctx, name, time.Time{})
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		slog.Warn("PrometheusClient GET Warnings INFO", slog.Any("warnings", warnings))
	}
	return values, nil
}

func (p *PrometheusClient) Reload(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodPost, p.address+"/-/reload", nil)
	if err != nil {
		return err
	}
	do, bytes, err := p.httpClient.Do(ctx, request)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return errors.Errorf("Reload failed,err:%s", bytes)
	}
	return nil
}
