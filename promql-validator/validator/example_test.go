// Copyright 2025 The Prometheus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator_test

import (
	"fmt"

	"github.com/prometheus/promql-validator/validator"
)

func ExampleValidate() {
	result := validator.Validate("http_requests_total")

	fmt.Println(result.Valid)
	fmt.Println(result.ExprType)
	fmt.Println(result.Metrics)
	// Output:
	// true
	// vector
	// [http_requests_total]
}

func ExampleValidate_invalid() {
	result := validator.Validate("sum(http_requests_total")

	fmt.Println(result.Valid)
	// Output:
	// false
}

func ExampleValidate_withLabels() {
	result := validator.Validate(`http_requests_total{job="api",method="GET"}`)

	fmt.Println(result.Valid)
	fmt.Println(result.Metrics)
	// Output:
	// true
	// [http_requests_total]
}

func ExampleValidate_aggregation() {
	result := validator.Validate(`sum by (job) (rate(http_requests_total[5m]))`)

	fmt.Println(result.Valid)
	fmt.Println(result.Functions)
	// Output:
	// true
	// [rate]
}
