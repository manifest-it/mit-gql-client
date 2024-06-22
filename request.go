package gql

import (
	"fmt"
	"strings"
)

type Query struct {
	Object string
	Where  Where
	Select Queries
	Wrap   bool
}

type Queries []Query

func (q Query) String() string {
	qParts := []string{q.Object}
	if len(q.Where) > 0 {
		qParts = append(qParts, q.Where.String())
	}
	if len(q.Select) == 0 {
		return WrapQuery(strings.Join(qParts, " "), q.Wrap)
	}
	subfields := []string{}
	for _, subfield := range q.Select {
		subfields = append(subfields, subfield.String())
	}
	if len(subfields) > 0 {
		qParts = append(qParts,
			"{"+strings.Join(subfields, " ")+"}")
	}
	return WrapQuery(strings.Join(qParts, " "), q.Wrap)
}

type Where map[string]any

func (w Where) String() string {
	if len(w) == 0 {
		return ""
	}
	parts := []string{}
	for k, v := range w {
		switch v := v.(type) {
		case int:
			parts = append(parts, fmt.Sprintf("%s: %d", k, v))
		case string:
			parts = append(parts, fmt.Sprintf("%s:\"%s\"", k, v))
		case []string:
			parts = append(parts, fmt.Sprintf("%s:[%s]", k, strings.Join(v, ", ")))
		case map[string]any:
			var rules []string
			for _, rule := range v["rules"].([]map[string]any) {
				var ruleParts []string
				for rk, rv := range rule {
					ruleParts = append(ruleParts, formatRule(rk, rv))
				}
				rules = append(rules, "{"+strings.Join(ruleParts, ", ")+"}")
			}
			parts = append(parts, fmt.Sprintf("%s: { rules: [%s] }", k, strings.Join(rules, ", ")))
		default:
			parts = append(parts, fmt.Sprintf("%s:unknown_type", k))
		}
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

func formatRule(key string, value interface{}) string {
	switch key {
	case "operator":
		return fmt.Sprintf("%s: %v", key, value)
	case "compare_value":
		var formattedValues []string
		for _, val := range value.([]string) {
			formattedValues = append(formattedValues, fmt.Sprintf("\"%s\"", val))
		}
		return fmt.Sprintf("%s: [%s]", key, strings.Join(formattedValues, ", "))
	default:
		return fmt.Sprintf("%s: \"%v\"", key, value)
	}
}

func WrapQuery(gql string, wrap bool) string {
	if wrap {
		return fmt.Sprintf("query {%s}", gql)
	}
	return gql
}
