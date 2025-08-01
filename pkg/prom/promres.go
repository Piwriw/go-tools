package prom

import (
	"errors"
	"fmt"

	prommodel "github.com/prometheus/common/model"
)

const (
	PROMVALUEKET          = "__prom_value__"
	PROMTIMESTAMPVALUEKET = "__prom_timestamp__"
)

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
