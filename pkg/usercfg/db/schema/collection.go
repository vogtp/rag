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
		field.String("APIKey"),
	}
}

func (Collection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Spaces", Space.Type),
		// edge.To("friends", User.Type),
	}
}

// func (Collection) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.Fields("Name").Unique(),
// 	}
// }
