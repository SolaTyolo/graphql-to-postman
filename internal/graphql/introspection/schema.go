package introspection

import (
	"github.com/samber/lo"
	"github.com/wundergraph/graphql-go-tools/pkg/introspection"
)

type Schema struct {
	introspection.Schema
}

func (s Schema) Types() []introspection.FullType {
	return s.Schema.Types
}

func (s Schema) GetQueryOperations() []introspection.Field {
	if s.Schema.QueryType == nil {
		return nil
	}
	for _, ty := range s.Schema.Types {
		if ty.Name == s.Schema.QueryType.Name {
			return ty.Fields
		}
	}
	return nil
}

func (s Schema) GetQueryOperationsWithoutDeprecated() []introspection.Field {
	return lo.Filter(s.GetQueryOperations(), func(f introspection.Field, _ int) bool {
		return !f.IsDeprecated
	})
}

func (s Schema) GetMutationOperation() []introspection.Field {
	if s.Schema.MutationType == nil {
		return nil
	}
	for _, ty := range s.Schema.Types {
		if ty.Name == s.Schema.MutationType.Name {
			return ty.Fields
		}
	}
	return nil
}

func (s Schema) GetMutationOperationWithoutDeprecated() []introspection.Field {
	return lo.Filter(s.GetMutationOperation(), func(f introspection.Field, _ int) bool {
		return !f.IsDeprecated
	})
}

func (s Schema) GetSubscribeOperation() []introspection.Field {
	if s.Schema.SubscriptionType == nil {
		return nil
	}
	for _, ty := range s.Schema.Types {
		if ty.Name == s.Schema.SubscriptionType.Name {
			return ty.Fields
		}
	}
	return nil
}

func (s Schema) GetSubscribeOperationWithoutDeprecated() []introspection.Field {
	return lo.Filter(s.GetSubscribeOperation(), func(f introspection.Field, _ int) bool {
		return !f.IsDeprecated
	})
}
