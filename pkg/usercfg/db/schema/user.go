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
		field.String("OpenaiAPIkey").Optional(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Confluence", Confluence.Type),
		edge.To("Collections", Collection.Type),
		// edge.To("friends", User.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("Name").Unique(),
		//	index.Fields("OpenaiAPIkey"),
	}
}
