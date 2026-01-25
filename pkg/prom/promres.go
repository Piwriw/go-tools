package prom

import (
	"errors"
	"fmt"
	"math"

	prommodel "github.com/prometheus/common/model"
)

const (
	PROMVALUEKET          = "__prom_value__"
	PROMTIMESTAMPVALUEKET = "__prom_timestamp__"
	MAXPoint              = 1024 // 默认最大数据点数量
)

// DownSampleMode 降采样模式
type DownSampleMode string

const (
	// DownSampleModeUniform 均匀采样模式（等间隔采样）
	DownSampleModeUniform DownSampleMode = "uniform"
	// DownSampleModeMax 最大值采样模式（每个区间取最大值）
	DownSampleModeMax DownSampleMode = "max"
	// DownSampleModeMin 最小值采样模式（每个区间取最小值）
	DownSampleModeMin DownSampleMode = "min"
	// DownSampleModeAverage 平均值采样模式（每个区间取平均值）
	DownSampleModeAverage DownSampleMode = "average"
	// DownSampleModeLTTB LTTB 采样模式（Largest Triangle Three Buckets，保留视觉特征）
	DownSampleModeLTTB DownSampleMode = "lttb"
)

// DownSampleOptions 降采样配置选项
type DownSampleOptions struct {
	// MaxPoints 最大数据点数量，默认为 MAXPoint
	MaxPoints int
	// Mode 降采样模式，默认为 DownSampleModeUniform
	Mode DownSampleMode
	// Enabled 是否启用降采样，默认为 true
	Enabled bool
}

type ResPromQL struct {
	Val prommodel.Value
}

func (r *ResPromQL) String() string {
	return r.Val.String()
}

func (r *ResPromQL) Rows() ([]Row, error) {
	resp, err := convertToFieldMap(r.Val)
	if err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

type Data struct {
	ResultType string       `json:"result_type"`
	Result     []ResultPair `json:"result"`
}

// DownSample 对数据进行降采样处理
// 遍历 Result 列表，并对每个元素调用其 DownSample 方法。
// 参数:
//   - maxPoints: 可选参数，表示最大允许的数据点数量，默认为 MAXPoint
func (r *Data) DownSample(maxPoints ...int) {
	for i := range r.Result {
		r.Result[i].DownSample(maxPoints...)
	}
}

// DownSampleWithOptions 使用配置选项对数据进行降采样处理
func (r *Data) DownSampleWithOptions(opts *DownSampleOptions) {
	for i := range r.Result {
		r.Result[i].DownSampleWithOptions(opts)
	}
}

// DownSample 对数据点进行等比缩小（使用默认的均匀采样模式）
// 当数据点数量超过 maxPoints 时，按等间隔采样，保留原始 Value 值
// 参数:
//   - maxPoints: int 类型，表示最大允许的数据点数量，默认为 MAXPoint
//
// 返回值:
//   - ResultPair: 缩小后的结果，如果原始点数 <= maxPoints，则返回原始数据
func (r *ResultPair) DownSample(maxPoints ...int) {
	var maxPoint int
	if len(maxPoints) == 0 {
		maxPoint = MAXPoint
	} else {
		maxPoint = maxPoints[0]
	}

	opts := &DownSampleOptions{
		MaxPoints: maxPoint,
		Mode:      DownSampleModeUniform,
		Enabled:   true,
	}
	r.DownSampleWithOptions(opts)
}

// DownSampleWithOptions 使用配置选项对数据点进行降采样
// 支持多种降采样模式：均匀采样、最大值、最小值、平均值、LTTB
// 参数:
//   - opts: 降采样配置选项
func (r *ResultPair) DownSampleWithOptions(opts *DownSampleOptions) {
	if r == nil || opts == nil || !opts.Enabled {
		return
	}

	// 设置默认值
	if opts.MaxPoints <= 0 {
		opts.MaxPoints = MAXPoint
	}
	if opts.Mode == "" {
		opts.Mode = DownSampleModeUniform
	}

	// 如果数据点数量不超过最大值，直接返回原始数据
	if len(r.Values) <= opts.MaxPoints {
		return
	}

	// 根据不同模式进行降采样
	sampler := setSampler(opts.Mode)
	r.Values = sampler.Sample(r.Values, opts.MaxPoints)
}

// Sampler 降采样器接口
type Sampler interface {
	// Sample 对数据进行采样
	// 参数:
	//   - values: 原始数据点
	//   - maxPoints: 最大数据点数量
	// 返回:
	//   - 采样后的数据点
	Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair
}

// setSampler 根据模式获取对应的采样器
func setSampler(mode DownSampleMode) Sampler {
	switch mode {
	case DownSampleModeUniform:
		return &uniformSampler{}
	case DownSampleModeMax:
		return &maxSampler{}
	case DownSampleModeMin:
		return &minSampler{}
	case DownSampleModeAverage:
		return &avgSampler{}
	case DownSampleModeLTTB:
		return &lttbSampler{}
	default:
		return &uniformSampler{}
	}
}

// ensureLastPoint 确保包含最后一个数据点
func ensureLastPoint(sampledValues, originalValues []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	if len(sampledValues) == 0 || len(originalValues) == 0 {
		return sampledValues
	}

	lastSampled := sampledValues[len(sampledValues)-1]
	lastOriginal := originalValues[len(originalValues)-1]

	if lastSampled.Timestamp != lastOriginal.Timestamp {
		if len(sampledValues) < maxPoints {
			sampledValues = append(sampledValues, lastOriginal)
		} else {
			sampledValues[len(sampledValues)-1] = lastOriginal
		}
	}

	return sampledValues
}

// uniformSampler 均匀采样器（等间隔采样）
type uniformSampler struct{}

func (s *uniformSampler) Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	if len(values) <= maxPoints {
		return values
	}

	// 计算采样步长
	step := len(values) / maxPoints
	if step < 1 {
		step = 1
	}

	// 按步长采样
	sampledValues := make([]prommodel.SamplePair, 0, maxPoints)
	for i := 0; i < len(values); i += step {
		sampledValues = append(sampledValues, values[i])
		if len(sampledValues) >= maxPoints {
			break
		}
	}

	// 确保包含最后一个数据点
	return ensureLastPoint(sampledValues, values, maxPoints)
}

type AggType string

const (
	AggTypeMax AggType = "max"
	AggTypeMin AggType = "min"
	AggTypeAvg AggType = "avg"
)

// maxSampler 最大值采样器
type maxSampler struct{}

func (s *maxSampler) Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	return aggregationSample(values, maxPoints, AggTypeMax)
}

// minSampler 最小值采样器
type minSampler struct{}

func (s *minSampler) Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	return aggregationSample(values, maxPoints, AggTypeMin)
}

// avgSampler 平均值采样器
type avgSampler struct{}

func (s *avgSampler) Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	return aggregationSample(values, maxPoints, AggTypeAvg)
}

// aggregationSample 聚合采样（支持 max、min、avg）
func aggregationSample(values []prommodel.SamplePair, maxPoints int, aggType AggType) []prommodel.SamplePair {
	if len(values) <= maxPoints {
		return values
	}

	// 计算每个桶的大小
	bucketSize := len(values) / maxPoints
	if bucketSize < 1 {
		bucketSize = 1
	}

	sampledValues := make([]prommodel.SamplePair, 0, maxPoints)

	for i := 0; i < len(values); i += bucketSize {
		end := i + bucketSize
		if end > len(values) {
			end = len(values)
		}

		// 获取桶内的数据点
		bucket := values[i:end]
		if len(bucket) == 0 {
			continue
		}

		// 根据聚合类型计算值
		var aggValue prommodel.SampleValue
		var aggTimestamp prommodel.Time

		switch aggType {
		case AggTypeMax:
			maxVal := bucket[0].Value
			maxIdx := 0
			for j, point := range bucket {
				if point.Value > maxVal {
					maxVal = point.Value
					maxIdx = j
				}
			}
			aggValue = maxVal
			aggTimestamp = bucket[maxIdx].Timestamp

		case AggTypeMin:
			minVal := bucket[0].Value
			minIdx := 0
			for j, point := range bucket {
				if point.Value < minVal {
					minVal = point.Value
					minIdx = j
				}
			}
			aggValue = minVal
			aggTimestamp = bucket[minIdx].Timestamp

		case AggTypeAvg:
			var sum prommodel.SampleValue
			for _, point := range bucket {
				sum += point.Value
			}
			aggValue = sum / prommodel.SampleValue(len(bucket))
			// 使用桶的中间时间戳
			aggTimestamp = bucket[len(bucket)/2].Timestamp
		}

		sampledValues = append(sampledValues, prommodel.SamplePair{
			Timestamp: aggTimestamp,
			Value:     aggValue,
		})

		if len(sampledValues) >= maxPoints {
			break
		}
	}

	// 确保包含最后一个数据点
	return ensureLastPoint(sampledValues, values, maxPoints)
}

// lttbSampler LTTB 采样器（Largest Triangle Three Buckets）
type lttbSampler struct{}

// Sample 使用 LTTB 算法进行采样
// 这是一种视觉优化的降采样算法，能够保留数据的视觉特征
// 参考: https://github.com/sveinn-steinarsson/flot-downsample
func (s *lttbSampler) Sample(values []prommodel.SamplePair, maxPoints int) []prommodel.SamplePair {
	if len(values) <= maxPoints {
		return values
	}

	if maxPoints <= 2 {
		// 至少需要 3 个点才能使用 LTTB 算法，降级为均匀采样
		uniformSampler := &uniformSampler{}
		return uniformSampler.Sample(values, maxPoints)
	}

	originalLength := len(values)
	sampledValues := make([]prommodel.SamplePair, 0, maxPoints)

	// 始终包含第一个点
	sampledValues = append(sampledValues, values[0])

	// 计算每个桶的大小
	bucketSize := float64(originalLength-2) / float64(maxPoints-2)

	// 用于存储上一个选中点的索引
	a := 0

	for i := 0; i < maxPoints-2; i++ {
		// 计算当前桶的平均点
		avgRangeStart := int(math.Floor(float64(i+1)*bucketSize)) + 1
		avgRangeEnd := int(math.Floor(float64(i+2)*bucketSize)) + 1
		if avgRangeEnd > originalLength {
			avgRangeEnd = originalLength
		}

		avgTimestamp := prommodel.Time(0)
		avgValue := prommodel.SampleValue(0)
		avgRangeLength := avgRangeEnd - avgRangeStart

		if avgRangeLength > 0 {
			for j := avgRangeStart; j < avgRangeEnd; j++ {
				avgTimestamp += values[j].Timestamp
				avgValue += values[j].Value
			}
			avgTimestamp /= prommodel.Time(avgRangeLength)
			avgValue /= prommodel.SampleValue(avgRangeLength)
		}

		// 计算当前桶的范围
		rangeOffs := int(math.Floor(float64(i)*bucketSize)) + 1
		rangeTo := int(math.Floor(float64(i+1)*bucketSize)) + 1

		// 点 a 的坐标
		pointATimestamp := float64(values[a].Timestamp)
		pointAValue := float64(values[a].Value)

		maxArea := -1.0
		var nextA int

		// 在当前桶中找到形成最大三角形面积的点
		for j := rangeOffs; j < rangeTo && j < originalLength; j++ {
			// 计算三角形面积
			area := math.Abs((pointATimestamp-float64(avgTimestamp))*(float64(values[j].Value)-pointAValue)-
				(pointATimestamp-float64(values[j].Timestamp))*(float64(avgValue)-pointAValue)) * 0.5

			if area > maxArea {
				maxArea = area
				nextA = j
			}
		}

		sampledValues = append(sampledValues, values[nextA])
		a = nextA
	}

	// 始终包含最后一个点
	sampledValues = append(sampledValues, values[originalLength-1])

	return sampledValues
}

func (r *ResPromQL) Len() (int, error) {
	resp, err := convertToFieldMap(r.Val)
	if err != nil {
		return 0, err
	}
	return len(resp.Rows), nil
}

type ResultPair struct {
	Metric prommodel.Metric       `json:"metric"`
	Values []prommodel.SamplePair `json:"values"`
}

func (r *ResPromQL) Data() (*Data, error) {
	resultType, err := r.ResultType()
	if err != nil {
		return nil, err
	}

	query := r.Val
	res := make([]ResultPair, 0)
	switch query.Type() {
	case prommodel.ValVector: // 瞬时向量结果
		vector, ok := query.(prommodel.Vector)
		if !ok {
			return nil, errors.New("query is not Vector")
		}
		values := make([]prommodel.SamplePair, 0)
		for _, sample := range vector {
			values = append(values, prommodel.SamplePair{
				Timestamp: sample.Timestamp,
				Value:     sample.Value,
			})
			res = append(res, ResultPair{
				Metric: sample.Metric,
				Values: values,
			})
		}

	case prommodel.ValMatrix: // 时间序列结果
		matrix, ok := query.(prommodel.Matrix)
		if !ok {
			return nil, errors.New("query is not Matrix")
		}
		values := make([]prommodel.SamplePair, 0)
		for _, series := range matrix {
			values = append(values, series.Values...)
			res = append(res, ResultPair{
				Metric: series.Metric,
				Values: values,
			})
		}
	default:
		return nil, fmt.Errorf("unsupported query result type: %T", query)
	}
	data := &Data{
		ResultType: resultType,
		Result:     res,
	}
	return data, nil
}

// ResultType 获取 PromQL 结果的类型
func (r *ResPromQL) ResultType() (string, error) {
	resp, err := convertToFieldMap(r.Val)
	if err != nil {
		return "", err
	}
	// 如果没找到，返回空切片和错误
	return resp.ResultType, nil
}

// Metric 获取 PromQL 结果中的 Metric 信息（所有的标签）
func (r *ResPromQL) Metric() ([]prommodel.Metric, error) {
	query := r.Val
	res := make([]prommodel.Metric, 0)
	switch query.Type() {
	case prommodel.ValVector: // 瞬时向量结果
		vector, ok := query.(prommodel.Vector)
		if !ok {
			return nil, errors.New("query is not Vector")
		}
		for _, sample := range vector {
			res = append(res, sample.Metric)
		}
	case prommodel.ValMatrix: // 时间序列结果
		matrix, ok := query.(prommodel.Matrix)
		if !ok {
			return nil, errors.New("query is not Matrix")
		}
		for _, series := range matrix {
			res = append(res, series.Metric)
		}
	default:
		return nil, fmt.Errorf("unsupported query result type: %T", query)
	}
	return res, nil
}

// Values 获取 PromQL 结果中的值
func (r *ResPromQL) Values() ([]prommodel.SamplePair, error) {
	query := r.Val
	res := make([]prommodel.SamplePair, 0)
	switch query.Type() {
	case prommodel.ValVector: // 瞬时向量结果
		vector, ok := query.(prommodel.Vector)
		if !ok {
			return nil, errors.New("query is not Vector")
		}
		for _, sample := range vector {
			res = append(res, prommodel.SamplePair{
				Timestamp: sample.Timestamp,
				Value:     sample.Value,
			})
		}
	case prommodel.ValMatrix: // 时间序列结果
		matrix, ok := query.(prommodel.Matrix)
		if !ok {
			return nil, errors.New("query is not Matrix")
		}
		for _, series := range matrix {
			res = append(res, series.Values...)
		}
	default:
		return nil, fmt.Errorf("unsupported query result type: %T", query)
	}
	return res, nil
}

type Rows []Row
type Resp struct {
	Rows       []Row  `json:"rows"`
	ResultType string `json:"result_type"`
}

var promValueTypeTable = map[prommodel.ValueType]string{
	prommodel.ValScalar: "scalar",
	prommodel.ValVector: "vector",
	prommodel.ValMatrix: "matrix",
	prommodel.ValString: "string",
}

func getPromQLValue(val prommodel.ValueType) string {
	str, ok := promValueTypeTable[val]
	if ok {
		return str
	}
	return ""
}

type Row map[string]any

// GetVal 获取Label值
func (r *Row) GetVal(name string) (any, error) {
	val, ok := (*r)[name]
	if !ok {
		return nil, fmt.Errorf("no field %s found", name)
	}
	return val, nil
}

// GetValStr 获取字符串类型的字段值
func (r *Row) GetValStr(name string) (string, error) {
	val, ok := (*r)[name]
	if !ok {
		return "", fmt.Errorf("no field %s found", name)
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not of type string", name)
	}
	return strVal, nil
}

// GetValue 获取 Prom Value 值
func (r *Row) GetValue() (float64, error) {
	val, ok := (*r)[PROMVALUEKET]
	if !ok {
		return 0, fmt.Errorf("no field value found")
	}
	floatVal, ok := val.(float64)
	if !ok {
		return 0, errors.New("field value is not of type float64")
	}
	return floatVal, nil
}

func convertToFieldMap(query prommodel.Value) (*Resp, error) {
	result := &Resp{
		Rows:       make(Rows, 0),
		ResultType: getPromQLValue(query.Type()),
	}
	switch query.Type() {
	case prommodel.ValVector: // 瞬时向量结果
		vector, ok := query.(prommodel.Vector)
		if !ok {
			return nil, errors.New("query is not Vector")
		}
		for _, sample := range vector {
			entry := make(map[string]interface{})
			// 添加标签字段和值
			for key, value := range sample.Metric {
				entry[string(key)] = string(value)
			}
			// 添加值和时间戳
			entry[PROMVALUEKET] = float64(sample.Value)
			entry[PROMTIMESTAMPVALUEKET] = sample.Timestamp.Time()
			result.Rows = append(result.Rows, Row(entry))
		}
	case prommodel.ValMatrix: // 时间序列结果
		matrix, ok := query.(prommodel.Matrix)
		if !ok {
			return nil, errors.New("query is not Matrix")
		}
		for _, series := range matrix {
			for _, point := range series.Values {
				entry := make(map[string]interface{})
				// 添加标签字段和值
				for key, value := range series.Metric {
					entry[string(key)] = string(value)
				}
				// 添加值和时间戳
				entry[PROMVALUEKET] = float64(point.Value)
				entry[PROMTIMESTAMPVALUEKET] = point.Timestamp.Time()
				result.Rows = append(result.Rows, Row(entry))
			}
		}
	default:
		return nil, fmt.Errorf("unsupported query result type: %T", query)
	}
	return result, nil
}
