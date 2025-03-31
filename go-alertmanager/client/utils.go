package client

import (
	"github.com/prometheus/prometheus/model/rulefmt"
	"gopkg.in/yaml.v3"
)

func RulesToRuleNodes(rs ...rulefmt.Rule) []rulefmt.RuleNode {
	if len(rs) == 0 {
		return nil
	}
	nodes := make([]rulefmt.RuleNode, 0, len(rs))
	for _, r := range rs {
		nodes = append(nodes, RuleToRuleNode(r))
	}
	return nodes
}

// RuleToRuleNode 将 Rule 转换为 RuleNode
func RuleToRuleNode(r rulefmt.Rule) rulefmt.RuleNode {
	return rulefmt.RuleNode{
		Record:        stringToYamlNode(r.Record),
		Alert:         stringToYamlNode(r.Alert),
		Expr:          stringToYamlNode(r.Expr),
		For:           r.For,
		KeepFiringFor: r.KeepFiringFor,
		Labels:        r.Labels,
		Annotations:   r.Annotations,
	}
}

// RuleNodeToRule 将 RuleNode 转换为 Rule
func RuleNodeToRule(rn rulefmt.RuleNode) rulefmt.Rule {
	return rulefmt.Rule{
		Record:        rn.Record.Value,
		Alert:         rn.Alert.Value,
		Expr:          rn.Expr.Value,
		For:           rn.For,
		KeepFiringFor: rn.KeepFiringFor,
		Labels:        rn.Labels,
		Annotations:   rn.Annotations,
	}
}

// RuleNodesToRule 将 RuleNode 转换为 Rule
func RuleNodesToRule(rns ...rulefmt.RuleNode) []rulefmt.Rule {
	rules := make([]rulefmt.Rule, 0, len(rns))
	for _, rn := range rns {
		rules = append(rules, RuleNodeToRule(rn))
	}
	return rules
}

// stringToYamlNode 将字符串转换为 yaml.Node
func stringToYamlNode(s string) yaml.Node {
	if s == "" {
		return yaml.Node{Kind: yaml.ScalarNode}
	}
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: s,
	}
}
