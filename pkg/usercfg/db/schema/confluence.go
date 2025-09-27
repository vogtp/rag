package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Confluence struct {
	ent.Schema
}

func (Confluence) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name"),
		field.String("URL"),
		field.String("ConfluenceAPIKey"),
	}
}

func (Confluence) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Spaces", Space.Type),
		// edge.To("friends", User.Type),
	}
}

// func (Confluence) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.Fields("Name").Unique(),
// 	}
// }
