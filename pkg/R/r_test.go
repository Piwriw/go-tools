package common

import (
	"encoding/json"
	"testing"
)

// TestNewR 测试 NewR 构造函数
// 测试场景：创建新的 R 实例
func TestNewR(t *testing.T) {
	t.Run("创建新R实例", func(t *testing.T) {
		r := NewR()
		if r == nil {
			t.Fatal("NewR 返回 nil")
		}
		// 验证默认值
		if r.Code != 0 {
			t.Errorf("期望 Code=0, 实际=%d", r.Code)
		}
		if r.Data != nil {
			t.Error("期望 Data=nil")
		}
		if r.Msg != "" {
			t.Errorf("期望 Msg='', 实际='%s'", r.Msg)
		}
		if r.Detail != "" {
			t.Errorf("期望 Detail='', 实际='%s'", r.Detail)
		}
	})
}

// TestR_SetData 测试 SetData 方法
// 测试场景：设置各种类型的数据
func TestR_SetData(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{name: "字符串数据", data: "test data"},
		{name: "整数数据", data: 42},
		{name: "浮点数数据", data: 3.14},
		{name: "布尔数据", data: true},
		{name: "切片数据", data: []int{1, 2, 3}},
		{name: "map数据", data: map[string]interface{}{"key": "value"}},
		{name: "nil数据", data: nil},
		{name: "结构体数据", data: struct{ Name string }{"Alice"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewR()
			result := r.SetData(tt.data)

			// 验证链式调用返回自身
			if result != r {
				t.Error("SetData 应返回 *R 以支持链式调用")
			}

			// 验证数据被设置
			if r.Data != tt.data {
				t.Errorf("期望 Data=%v, 实际=%v", tt.data, r.Data)
			}
		})
	}
}

// TestR_SetCode 测试 SetCode 方法
// 测试场景：设置各种状态码
func TestR_SetCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{name: "成功状态码", code: CodeSuccess},
		{name: "失败状态码", code: CodeFailed},
		{name: "自定义状态码", code: 500},
		{name: "零状态码", code: 0},
		{name: "负数状态码", code: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewR()
			result := r.SetCode(tt.code)

			// 验证链式调用返回自身
			if result != r {
				t.Error("SetCode 应返回 *R 以支持链式调用")
			}

			// 验证状态码被设置
			if r.Code != tt.code {
				t.Errorf("期望 Code=%d, 实际=%d", tt.code, r.Code)
			}
		})
	}
}

// TestR_SetMsg 测试 SetMsg 方法
// 测试场景：设置各种消息
func TestR_SetMsg(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{name: "成功消息", msg: "操作成功"},
		{name: "失败消息", msg: "操作失败"},
		{name: "空消息", msg: ""},
		{name: "特殊字符消息", msg: "消息包含特殊字符: !@#$%^&*()"},
		{name: "中文消息", msg: "操作成功完成"},
		{name: "长消息", msg: string(make([]byte, 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewR()
			result := r.SetMsg(tt.msg)

			// 验证链式调用返回自身
			if result != r {
				t.Error("SetMsg 应返回 *R 以支持链式调用")
			}

			// 验证消息被设置
			if r.Msg != tt.msg {
				t.Errorf("期望 Msg='%s', 实际='%s'", tt.msg, r.Msg)
			}
		})
	}
}

// TestR_SetFailed 测试 SetFailed 方法
// 测试场景：设置失败信息
func TestR_SetFailed(t *testing.T) {
	tests := []struct {
		name   string
		msg    string
		detail string
	}{
		{
			name:   "正常失败信息",
			msg:    "请求失败",
			detail: "参数验证错误",
		},
		{
			name:   "空消息和详情",
			msg:    "",
			detail: "",
		},
		{
			name:   "只有消息没有详情",
			msg:    "操作失败",
			detail: "",
		},
		{
			name:   "只有详情没有消息",
			msg:    "",
			detail: "详细错误信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewR()
			result := r.SetFailed(tt.msg, tt.detail)

			// 验证链式调用返回自身
			if result != r {
				t.Error("SetFailed 应返回 *R 以支持链式调用")
			}

			// 验证状态码被设置为失败
			if r.Code != CodeFailed {
				t.Errorf("期望 Code=%d, 实际=%d", CodeFailed, r.Code)
			}

			// 验证消息被设置
			if r.Msg != tt.msg {
				t.Errorf("期望 Msg='%s', 实际='%s'", tt.msg, r.Msg)
			}

			// 验证详情被设置
			if r.Detail != tt.detail {
				t.Errorf("期望 Detail='%s', 实际='%s'", tt.detail, r.Detail)
			}
		})
	}
}

// TestR_SetSuccess 测试 SetSuccess 方法
// 测试场景：设置成功信息
func TestR_SetSuccess(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{name: "正常成功消息", msg: "操作成功"},
		{name: "空成功消息", msg: ""},
		{name: "中文成功消息", msg: "创建成功"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewR()
			result := r.SetSuccess(tt.msg)

			// 验证链式调用返回自身
			if result != r {
				t.Error("SetSuccess 应返回 *R 以支持链式调用")
			}

			// 验证状态码被设置为成功
			if r.Code != CodeSuccess {
				t.Errorf("期望 Code=%d, 实际=%d", CodeSuccess, r.Code)
			}

			// 验证消息被设置
			if r.Msg != tt.msg {
				t.Errorf("期望 Msg='%s', 实际='%s'", tt.msg, r.Msg)
			}

			// 验证详情为空（SetSuccess 不设置 Detail）
			if r.Detail != "" {
				t.Errorf("期望 Detail='', 实际='%s'", r.Detail)
			}
		})
	}
}

// TestR_ChainedCalls 测试链式调用
// 测试场景：连续调用多个方法
func TestR_ChainedCalls(t *testing.T) {
	t.Run("完整链式调用", func(t *testing.T) {
		data := map[string]interface{}{"id": 123, "name": "test"}

		r := NewR().
			SetCode(CodeSuccess).
			SetMsg("操作成功").
			SetData(data)

		if r.Code != CodeSuccess {
			t.Errorf("期望 Code=%d, 实际=%d", CodeSuccess, r.Code)
		}
		if r.Msg != "操作成功" {
			t.Errorf("期望 Msg='操作成功', 实际='%s'", r.Msg)
		}
		if r.Data == nil {
			t.Error("期望 Data 被设置")
		}
	})

	t.Run("SetFailed 链式调用后继续设置数据", func(t *testing.T) {
		data := "error details"

		r := NewR().
			SetFailed("操作失败", "详细错误").
			SetData(data)

		if r.Code != CodeFailed {
			t.Errorf("期望 Code=%d, 实际=%d", CodeFailed, r.Code)
		}
		if r.Msg != "操作失败" {
			t.Errorf("期望 Msg='操作失败', 实际='%s'", r.Msg)
		}
		if r.Detail != "详细错误" {
			t.Errorf("期望 Detail='详细错误', 实际='%s'", r.Detail)
		}
		if r.Data != data {
			t.Error("期望 Data 被设置")
		}
	})
}

// TestR_JSONSerialization 测试 JSON 序列化
// 测试场景：将 R 序列化为 JSON
func TestR_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		r        *R
		expected string
	}{
		{
			name:     "空R",
			r:        NewR(),
			expected: `{"code":0,"data":null,"msg":"","detail":""}`,
		},
		{
			name:     "成功响应",
			r:        NewR().SetCode(CodeSuccess).SetMsg("success").SetData("data"),
			expected: `{"code":200,"data":"data","msg":"success","detail":""}`,
		},
		{
			name:     "失败响应",
			r:        NewR().SetCode(CodeFailed).SetMsg("error").SetDetail("error detail"),
			expected: `{"code":400,"data":null,"msg":"error","detail":"error detail"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.r)
			if err != nil {
				t.Fatalf("序列化失败: %v", err)
			}

			jsonStr := string(jsonBytes)
			if jsonStr != tt.expected {
				t.Errorf("期望 JSON='%s', 实际='%s'", tt.expected, jsonStr)
			}
		})
	}
}

// TestR_JSONDeserialization 测试 JSON 反序列化
// 测试场景：从 JSON 反序列化 R
func TestR_JSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		wantCode int
		wantMsg  string
		wantData interface{}
	}{
		{
			name:     "空响应",
			jsonStr:  `{"code":0,"data":null,"msg":"","detail":""}`,
			wantCode: 0,
			wantMsg:  "",
		},
		{
			name:     "成功响应",
			jsonStr:  `{"code":200,"data":"test data","msg":"success","detail":""}`,
			wantCode: 200,
			wantMsg:  "success",
			wantData: "test data",
		},
		{
			name:     "带对象数据的响应",
			jsonStr:  `{"code":200,"data":{"key":"value"},"msg":"ok","detail":""}`,
			wantCode: 200,
			wantMsg:  "ok",
			wantData: map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r R
			err := json.Unmarshal([]byte(tt.jsonStr), &r)
			if err != nil {
				t.Fatalf("反序列化失败: %v", err)
			}

			if r.Code != tt.wantCode {
				t.Errorf("期望 Code=%d, 实际=%d", tt.wantCode, r.Code)
			}
			if r.Msg != tt.wantMsg {
				t.Errorf("期望 Msg='%s', 实际='%s'", tt.wantMsg, r.Msg)
			}
			if tt.wantData != nil && r.Data != tt.wantData {
				t.Errorf("期望 Data=%v, 实际=%v", tt.wantData, r.Data)
			}
		})
	}
}

// TestCodeConstants 测试状态码常量
func TestCodeConstants(t *testing.T) {
	t.Run("验证常量值", func(t *testing.T) {
		if CodeSuccess != 200 {
			t.Errorf("期望 CodeSuccess=200, 实际=%d", CodeSuccess)
		}
		if CodeFailed != 400 {
			t.Errorf("期望 CodeFailed=400, 实际=%d", CodeFailed)
		}
	})
}

// BenchmarkNewR 性能基准测试
func BenchmarkNewR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewR()
	}
}

// BenchmarkR_SetData 性能基准测试
func BenchmarkR_SetData(b *testing.B) {
	r := NewR()
	data := "test data"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.SetData(data)
	}
}

// BenchmarkR_ChainedCalls 性能基准测试
func BenchmarkR_ChainedCalls(b *testing.B) {
	data := map[string]interface{}{"key": "value"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewR().
			SetCode(CodeSuccess).
			SetMsg("success").
			SetData(data)
	}
}

// ExampleNewR 示例函数
func ExampleNewR() {
	// 创建一个空的响应对象
	r := NewR()

	_ = r
}

// ExampleR_SetSuccess 示例函数
func ExampleR_SetSuccess() {
	// 创建成功响应
	r := NewR().
		SetSuccess("操作成功").
		SetData(map[string]interface{}{"id": 123})

	_ = r
}

// ExampleR_SetFailed 示例函数
func ExampleR_SetFailed() {
	// 创建失败响应
	r := NewR().
		SetFailed("操作失败", "参数验证错误")

	_ = r
}

// ExampleR_Chained 示例函数 - 链式调用
func ExampleR_Chained() {
	// 使用链式调用创建完整响应
	r := NewR().
		SetCode(CodeSuccess).
		SetMsg("创建成功").
		SetData(map[string]interface{}{
			"id":   123,
			"name": "测试用户",
		})

	_ = r
}
