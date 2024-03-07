package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"time"
)

const PROM_CPU_USAGE = "100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"

type PromClient v1.API

func main() {
	cli, err := NewPromCli("http://10.10.101.159:30009")
	if err != nil {
		panic(err)
	}
	res, err := ExecutePromQLQuery(context.TODO(), cli, PROM_CPU_USAGE)
	//QuerySeries(context.TODO(), cli)
	fmt.Println(res)
}
func NewPromCli(promURL string) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address: promURL,
	})
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
		return nil, err
	}

	apiClient := v1.NewAPI(client)
	return apiClient, nil
}

func ExecutePromQLQuery(ctx context.Context, cliProm v1.API, query string) (model.Value, error) {
	result, _, err := cliProm.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error executing PromQL query: %w", err)
	}
	return result, nil
}

func QuerySeries(ctx context.Context, cliProm v1.API) {
	lbls, _, err := cliProm.Series(ctx, []string{
		"{__name__=~\".+\", job=\"node-exporter\"}",
	}, time.Now().Add(-time.Hour), time.Now())
	//fmt.Println(warnings)
	fmt.Println(lbls)
	if err != nil {

	}
}
