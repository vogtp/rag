package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Space struct {
	ent.Schema
}

func (Space) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name"),
		field.String("SpaceKey"),
	}
}
