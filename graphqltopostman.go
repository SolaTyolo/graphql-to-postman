package graphqltopostman

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/SolaTyolo/graphqltopostman/internal/graphql/introspection"
	"github.com/samber/lo"
	intro_spec "github.com/wundergraph/graphql-go-tools/pkg/introspection"

	"github.com/SolaTyolo/graphqltopostman/internal/postman"
)

func ConvertToFile(name, url string, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return Convert(name, url, f)
}

func Convert(name, url string, w io.Writer) error {
	if len(url) == 0 {
		return errors.New("url is empty")
	}

	dataSchema, err := introspection.Introspect(url)
	if err != nil {
		return err
	}

	typeGroups := lo.SliceToMap(dataSchema.Types(), func(t intro_spec.FullType) (string, intro_spec.FullType) {
		return t.Name, t
	})

	c := postman.NewCollection(name)

	if len(dataSchema.GetQueryOperationsWithoutDeprecated()) != 0 {
		item, err := convertItem(
			url,
			// dataSchema.QueryType.Name,
			"query",
			dataSchema.GetQueryOperationsWithoutDeprecated(),
			typeGroups,
		)
		if err != nil {
			return err
		}
		c.AddItem(item)
	}

	if len(dataSchema.GetMutationOperationWithoutDeprecated()) != 0 {
		item, err := convertItem(
			url,
			// dataSchema.MutationType.Name,
			"mutation",
			dataSchema.GetMutationOperationWithoutDeprecated(),
			typeGroups,
		)
		if err != nil {
			return err
		}
		c.AddItem(item)
	}

	// TODO auth

	// TODO diretives

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func convertItem(domain, name string, op []intro_spec.Field, allTypes map[string]intro_spec.FullType) (postman.ItemGroup, error) {
	g := postman.NewItemGroup(name)
	for _, o := range op {
		item := postman.NewGraphqlItem(fmt.Sprintf("%s:%s", name, o.Name), domain, postman.GraphqlQueryFromOp(name, o, allTypes))
		item.Description = o.Description
		g.AddItem(item)
	}

	return g, nil
}
