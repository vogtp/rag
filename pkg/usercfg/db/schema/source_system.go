package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type SourceSystem struct {
	ent.Schema
}

func (SourceSystem) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name"),
		field.Enum("Type").Values("http", "confluence"),
		field.String("URL"),
		field.String("key").Optional(),
		field.String("parts").Optional(),
	}
}
