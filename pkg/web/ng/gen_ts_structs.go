//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/OneOfOne/struct2ts"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
)

func main() {
	s2ts := struct2ts.New(&struct2ts.Options{
		// NoConstructor:    true,
		// NoToObject:       true,
		// NoHelpers:        true,
	})
	fileName := "./intrasearch/src/app/services/settings.service.structs.ts"
	fmt.Printf("Generation TS structs for User in %s\n", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	s2ts.Add(ent.User{})
	if err := s2ts.RenderTo(f); err != nil {
		fmt.Printf("Cannot render ts structs: %v", err)
	}
	//s2ts.RenderTo(os.Stdout)
}
