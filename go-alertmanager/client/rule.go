package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/rulefmt"
	"gopkg.in/yaml.v3"
)

// AddAlertRule 向指定路径的规则文件添加新的规则组
func AddAlertRule(filePath string, newGroup rulefmt.RuleGroup) error {
	existingRules := &rulefmt.RuleGroups{}
	// 1. 检查规则文件是否存在
	if _, err := os.Stat(filePath); os.IsExist(err) {
		// 2. 读取现有规则文件
		parsedRules, errs := rulefmt.ParseFile(filePath)
		if len(errs) > 0 {
			return fmt.Errorf("error parsing rules file: %v", errs)
		}
		existingRules = parsedRules
	}

	// 3. 添加新规则组
	existingRules.Groups = append(existingRules.Groups, newGroup)

	// 4. 验证规则
	if err := Validate(existingRules); err != nil {
		return err
	}

	// 5. 以读写模式打开文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open rules file: %w", err)
	}
	defer file.Close()

	// 6. 写回文件
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	defer encoder.Close() // 关闭 Encoder 以确保所有内容写入文件

	if err := encoder.Encode(existingRules); err != nil {
		return fmt.Errorf("failed to write rules: %w", err)
	}

	return nil
}

// Validate validates all rules in the rule groups.
func Validate(g *rulefmt.RuleGroups) error {
	set := map[string]struct{}{}
	var errs []error
	for _, g := range g.Groups {
		if g.Name == "" {
			errs = append(errs, errors.New("Groupname must not be empty"))
		}

		if _, ok := set[g.Name]; ok {
			errs = append(
				errs,
				fmt.Errorf(" groupname: \"%s\" is repeated in the same file", g.Name),
			)
		}

		set[g.Name] = struct{}{}

		for _, r := range g.Rules {
			errorArr := r.Validate()
			for _, wrappedError := range errorArr {
				errs = append(errs, errors.New(wrappedError.Error()))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

// GenerateRuleGroup 创建符合 rulefmt 规范的规则组
func GenerateRuleGroup(name string, interval model.Duration, rules ...rulefmt.RuleNode) rulefmt.RuleGroup {
	rg := rulefmt.RuleGroup{
		Name:  name,
		Rules: rules,
	}

	// 只有当interval非零时才设置
	if interval != 0 {
		rg.Interval = interval
	}

	return rg
}

// CreateAlertingRule 创建符合 rulefmt 规范的告警规则节点
func CreateAlertingRule(
	alert string,
	expr string,
	labels map[string]string,
	annotations map[string]string,
	duration model.Duration,
) rulefmt.RuleNode {
	// 必须使用 yaml.Node 包装字符串字段
	return rulefmt.RuleNode{
		Alert: yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: alert,
		},
		Expr: yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: expr,
		},
		For:         duration,
		Labels:      labels,
		Annotations: annotations,
	}
}
