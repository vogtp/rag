package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").Unique(),
		field.String("APIKey").Optional(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Collections", Collection.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("Name").Unique(),
	}
}
