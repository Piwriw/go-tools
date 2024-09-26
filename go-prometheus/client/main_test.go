package client

import (
	"log/slog"
	"testing"
	"time"

	"github.com/prometheus/prometheus/promql/parser"
)

func TestQuery(t *testing.T) {
	client, err := NewPrometheusClient("http://10.0.0.195:9002")
	if err != nil {
		t.Fatal(err)
	}
	query, err := client.Query("is_alive{category=\"database.gaussdb\"}", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(query)

}

func TestQueryRange(t *testing.T) {
	active := new(bool)
	t.Logf("%v", *active)
}

func TestValidted(t *testing.T) {
	// 要校验的 PromQL 查询
	query := "111"
	client, err := NewPrometheusClient("http://10.0.0.195:9002")
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Validate(query); err != nil {
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
