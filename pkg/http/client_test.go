package http

import (
	"encoding/json"
	"testing"
	"time"
)

// 定义数据结构
type AlertData struct {
	Version           string            `json:"Version"`
	Status            string            `json:"Status"`
	Receiver          string            `json:"Receiver"`
	GroupLabels       map[string]string `json:"GroupLabels"`
	CommonLabels      map[string]string `json:"CommonLabels"`
	CommonAnnotations map[string]string `json:"CommonAnnotations"`
	ExternalURL       string            `json:"ExternalURL"`
	Alerts            []Alert           `json:"Alerts"`
	Fingerprint       string            `json:"Fingerprint"`
}

type Alert struct {
	Status        string            `json:"Status"`
	Labels        map[string]string `json:"Labels"`
	Annotations   map[string]string `json:"Annotations"`
	StartsAt      string            `json:"StartsAt"`
	EndsAt        string            `json:"EndsAt"`
	GeneratorURL  string            `json:"GeneratorURL"`
	Fingerprint   string            `json:"Fingerprint"`
	AlertConfigID int               `json:"AlertConfigID"`
	InstanceID    int               `json:"InstanceID"`
	AlertExpr     string            `json:"AlertExpr"`
}

func TestPost(t *testing.T) {
	// 构造 body 数据
	body := AlertData{
		Version:  "4",
		Status:   "resolved",
		Receiver: "default-receiver",
		GroupLabels: map[string]string{
			"alertname": "AgileX 元数据库健康检查",
			"instance":  "instance-0",
			"severity":  "致命:#4A4A4A",
		},
		CommonLabels: map[string]string{
			"alertname": "AgileX 元数据库健康检查",
			"category":  "server.AgileX",
			"component": "agilex",
			"groupId":   "0",
			"instance":  "instance-0",
			"ip":        "postgresql",
			"job":       "pushgateway",
			"name":      "AgileX 元数据库",
			"severity":  "致命:#4A4A4A",
			"usefor":    "server.AgileX",
		},
		CommonAnnotations: map[string]string{
			"description": "AgileX 元数据库异常，AgileX 服务不可用",
			"suggestion":  "请检查 AgileX 元数据库状态，测试从 AgileX 后台服务器或对应的容器连接 AgileX 元数据库是否正常",
			"summary":     "AgileX 元数据库异常，请检查 AgileX 数据库状态",
		},
		ExternalURL: "http://5b777b30f66d:9003",
		Alerts: []Alert{
			{
				Status: "resolved",
				Labels: map[string]string{
					"alertname": "AgileX 元数据库健康检查",
					"category":  "server.AgileX",
					"component": "agilex",
					"groupId":   "0",
					"instance":  "instance-0",
					"ip":        "postgresql",
					"job":       "pushgateway",
					"name":      "AgileX 元数据库",
					"severity":  "致命:#4A4A4A",
					"usefor":    "server.AgileX",
				},
				Annotations: map[string]string{
					"description": "AgileX 元数据库异常，AgileX 服务不可用",
					"suggestion":  "请检查 AgileX 元数据库状态，测试从 AgileX 后台服务器或对应的容器连接 AgileX 元数据库是否正常",
					"summary":     "AgileX 元数据库异常，请检查 AgileX 数据库状态",
				},
				StartsAt:      "2024-11-18T15:40:11.610+08:00",
				EndsAt:        "2024-11-18T15:40:11.610+08:00",
				GeneratorURL:  "http://a790b901af21:9002/graph?g0.expr=metadb_up%7Bcomponent%3D%22agilex%22%7D%3D%3D0&g0.tab=1",
				Fingerprint:   "51335ea284d8a6e8",
				AlertConfigID: 0,
				InstanceID:    0,
				AlertExpr:     "",
			},
		},
		Fingerprint: "",
	}
	for i := 0; i < 7; i++ {
		body.Alerts = append(body.Alerts, body.Alerts...)
	}
	for i := 0; i < 10; i++ {
		marshal, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		post, err := NewHTTPClient().JSON().Post("http://10.0.0.194:9001/v1/write?configid=1", marshal)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)
		t.Log("i", i, string(post))
	}

}
