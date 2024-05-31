package introspection

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SolaTyolo/httpclient"
	"github.com/wundergraph/graphql-go-tools/pkg/introspection"
)

type response struct {
	Data introspection.Data `json:"data"`
}

func Introspect(target string) (*Schema, error) {
	r, err := httpclient.Default().Post(context.Background(), target, map[string]interface{}{
		"operationName": "IntrospectionQuery",
		"query":         Query,
	})
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("response status is not 200 (OK)")
	}
	if r.RawBody == nil {
		return nil, errors.New("response body is empty")
	}
	var data response
	if err := json.Unmarshal(r.RawBody, &data); err != nil {
		return nil, err
	}

	return &Schema{data.Data.Schema}, nil
}
