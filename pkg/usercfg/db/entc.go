//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {

	ex, err := entgql.NewExtension(entgql.WithSchemaGenerator(), entgql.WithSchemaPath("./schema/ent.graphql"), entgql.WithConfigPath("./graph/gqlgen.yml")) // Generate GQL schema from the Ent's schema.
	//	entgql.WithSchemaGenerator(),
	// Generate the filters to a separate schema
	// file and load it in the gqlgen.yml config.
	// entgql.WithSchemaPath("../api/schema/ent.graphql"),
	//entgql.WithWhereFilters(true),

	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	err = entc.Generate("./schema", &gen.Config{
		Schema:  "github.com/vogtp/rag/pkg/usercfg/db/schema",
		Target:  "ent/",
		Package: "github.com/vogtp/rag/pkg/usercfg/db/ent",
		Features: []gen.Feature{
			// gen.FeaturePrivacy,
			gen.FeatureUpsert,
		},
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("struct methods").Parse(structUpdate)),
			// 	gen.MustParse(gen.NewTemplate("static").
			// 		Funcs(template.FuncMap{"title": strings.ToTitle}).
			// 		ParseFiles("template/static.tmpl")),
		},
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

var structUpdate = `package ent

{{ range $n := $.Nodes }}
    {{ $builder := $n.CreateName }}
    {{ $receiver := $n.CreateReceiver }}

    func ({{ $receiver }} *{{ $builder }}) Set{{ $n.Name }}(input *{{ $n.Name }}) *{{ $builder }} {
        {{- range $f := $n.Fields }}
            {{- $setter := print "Set" $f.StructField }}
            {{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
        {{- end }}
        return {{ $receiver }}
    }
		{{ $builder := $n.UpdateName }}
    {{ $receiver := $n.UpdateReceiver }}

    func ({{ $receiver }} *{{ $builder }}) Set{{ $n.Name }}(input *{{ $n.Name }}) *{{ $builder }} {
        {{- range $f := $n.Fields }}
            {{- $setter := print "Set" $f.StructField }}
            {{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
        {{- end }}
        return {{ $receiver }}
    }
		{{ $builder := $n.UpdateOneName }}
    {{ $receiver := $n.UpdateOneReceiver }}

    func ({{ $receiver }} *{{ $builder }}) Set{{ $n.Name }}(input *{{ $n.Name }}) *{{ $builder }} {
        {{- range $f := $n.Fields }}
            {{- $setter := print "Set" $f.StructField }}
            {{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
        {{- end }}
        return {{ $receiver }}
    }
{{ end }}
 
`
