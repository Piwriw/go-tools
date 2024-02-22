package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

func main() {

	var timeout int
	flag.IntVar(&timeout, "timeout", 32400, "Timeout in seconds")
	flag.Parse()

	buff := new(bytes.Buffer)

	fmt.Printf("Timeout: %d seconds\n", timeout)

	// 配置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 配置连接方式， "local" 表示，无论你后面 Inventory 选项怎么配，都是在执行 ansible 的本地执行
	// 要真正连接到 Inventory 配置的机器，注释掉 Connection 选项或者使用 "smart" 或 "ssh" 作为参数值
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
		// User:       "apenella",
	}

	// 资产清单和变量文件，也可以是一个 map 类型来作为变量，就无需引入文件
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: "10.10.114.40,",
		ExtraVarsFile: []string{
			"@vars-file1.yml",
		},
	}

	// 执行结果缓存
	measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
	)

	// 构造 ansible
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
		// Exec: executorTimeMeasurement,
		StdoutCallback: "json",
	}

	// 执行 playbook
	err := playbook.Run(ctx)
	if err != nil {
		panic(err)
	}

	// 输出结果
	res, err := results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
