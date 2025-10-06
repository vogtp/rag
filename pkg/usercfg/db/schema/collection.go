package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Collection struct {
	ent.Schema
}

func (Collection) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name"),
		field.String("APIKey").Optional(),
		field.String("collectionName").Optional(), // internal name of the vecDB collection
	}
}

func (Collection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Sources", SourceSystem.Type),
	}
}
