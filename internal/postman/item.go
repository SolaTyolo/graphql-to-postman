package postman

import (
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
	intro_spec "github.com/wundergraph/graphql-go-tools/pkg/introspection"
)

func NewItemGroup(name string) ItemGroup {
	return ItemGroup{
		Name: lo.ToPtr(name),
		Item: make([]interface{}, 0),
	}
}

func (g *ItemGroup) AddItem(item Item) {
	g.Item = append(g.Item, item)
}

type GraphqlItem = Item

type GraphqlQuery struct {
	Query     string
	Variables map[string]interface{}
}

func (q GraphqlQuery) AsMap() map[string]string {
	b, _ := json.Marshal(q.Variables)
	return map[string]string{
		"query":     q.Query,
		"variables": string(b),
	}
}

func NewGraphqlItem(name, domain string, query GraphqlQuery) GraphqlItem {
	return GraphqlItem{
		Name: lo.ToPtr(name),
		Request: IRequest{
			Method: http.MethodPost,
			Body: &Body{
				Mode:    lo.ToPtr(Graphql),
				Graphql: query.AsMap(),
			},
			URL:    domain,
			Header: []Header{},
		},
		Response: []Response{},
	}
}

func GraphqlQueryFromOp(opType string, op intro_spec.Field, alltypes map[string]intro_spec.FullType) GraphqlQuery {
	q := GraphqlQuery{
		Query:     "",
		Variables: map[string]interface{}{},
	}

	var count int
	var line1, line2 string
	for _, v := range op.Args {
		count++

		line1 += `$` + v.Name + `: ` + v.Name
		if v.Type.Kind == intro_spec.NONNULL {
			line1 += `!`
		}

		line2 += v.Name + `: $` + v.Name

		// Not the last one? Add the delimiter
		if count != len(op.Args) {
			line1 += `, `
			line2 += `, `
		}
	}
	q.Query = opType + ` ` + op.Name + `(` + line1 + `) { ` + op.Name + `(` + line2 + `) { __typename } }`

	count = 0
	for _, v := range op.Args {
		count++

		// slice not supported
		if v.Type.Kind == intro_spec.LIST {
			continue
		}
		val := graphqlQueryVariablesArgs(v.Name, v.Type.Kind.String(), alltypes)
		if val != nil {
			q.Variables[v.Name] = val
		}
	}

	return q
}

func graphqlQueryVariablesArgs(name string, kind string, alltypes map[string]intro_spec.FullType) interface{} {
	t, ok := alltypes[name]
	if !ok {
		return nil
	}
	switch kind {
	case intro_spec.SCALAR.String():
		return graphqlScalarValue(t.Name)
	case intro_spec.INPUTOBJECT.String():
		dval := map[string]interface{}{}
		for _, v := range t.InputFields {
			dummyVal := graphqlQueryVariablesArgs(v.Name, v.Type.Kind.String(), alltypes)
			if v.Type.Kind == intro_spec.LIST {
				dval[v.Name] = []interface{}{dummyVal}
			} else {
				dval[v.Name] = dummyVal
			}
		}
		return dval
	case intro_spec.ENUM.String():
		// 从 alltype 中获取 enumvalues； 并且设置为第一个值
		if len(t.EnumValues) == 0 {
			return "NULL"
		}
		return t.EnumValues[0].Name
	default:
		return "NULL"
	}

}

func graphqlScalarValue(name string) interface{} {
	switch name {
	case "String":
		return "title test"
	case "Float":
		return 1.0
	case "Integer":
		return 1
	case "Boolean":
		return true
	case "ID":
		return ""
	default:
		return ""
	}
}
