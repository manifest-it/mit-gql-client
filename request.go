package gql

import (
	"fmt"
	"strings"
)

type GraphQLOperation interface {
	String() string
}

type Query struct {
	Object string
	Where  Where
	Select Queries
	Wrap   bool
}

type Where map[string]any

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
		qParts = append(qParts, fmt.Sprintf("{%s}", strings.Join(subfields, " ")))
	}
	return WrapQuery(strings.Join(qParts, " "), q.Wrap)
}

// todo: refactor so that we don't need to hardcode specific fields for each provider
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
			// this is specific to monday.com
			if k == "column_values" {
				parts = append(parts, fmt.Sprintf("%s: %s", k, v))
			} else {
				parts = append(parts, fmt.Sprintf("%s:\"%s\"", k, v))
			}
		case []string:
			parts = append(parts, fmt.Sprintf("%s:[%s]", k, strings.Join(v, ", ")))
		case map[string]any:
			// this is specific to monday.com
			if rulesObj, ok := v["rules"]; ok {
				var rules []string
				for _, rule := range rulesObj.([]map[string]any) {
					var ruleParts []string
					for rk, rv := range rule {
						ruleParts = append(ruleParts, formatRule(rk, rv))
					}
					rules = append(rules, fmt.Sprintf("{%s}", strings.Join(ruleParts, ", ")))
				}
				parts = append(parts, fmt.Sprintf("%s: { rules: [%s] }", k, strings.Join(rules, ", ")))
			}
			// these are specific to panther.com
			if cursor, ok := v["cursor"]; ok {
				parts = append(parts, fmt.Sprintf("%s: { cursor: \"%s\" }", k, cursor.(string)))
			}
			if sql, ok := v["sql"]; ok {
				parts = append(parts, fmt.Sprintf("%s: { sql: \"%s\" }", k, sql.(string)))
			}
		default:
			parts = append(parts, fmt.Sprintf("%s:unknown_type", k))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
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

type Mutation struct {
	Name   string
	Input  Where
	Select Queries
	Wrap   bool
}

func (m Mutation) String() string {
	mParts := []string{m.Name}
	if len(m.Input) > 0 {
		mParts = append(mParts, m.Input.String())
	}
	if len(m.Select) == 0 {
		return WrapMutation(strings.Join(mParts, " "), m.Wrap)
	}
	subfields := []string{}
	for _, subfield := range m.Select {
		subfields = append(subfields, subfield.String())
	}
	if len(subfields) > 0 {
		mParts = append(mParts, fmt.Sprintf("{%s}", strings.Join(subfields, " ")))
	}
	return WrapMutation(strings.Join(mParts, " "), m.Wrap)
}

func WrapMutation(gql string, wrap bool) string {
	if wrap {
		return fmt.Sprintf("mutation {%s}", gql)
	}
	return gql
}
