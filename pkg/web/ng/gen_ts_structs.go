//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/OneOfOne/struct2ts"
	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
)

func main() {
	generate("./intrasearch/src/app/services/settings.service.structs.ts", ent.User{})
	generate("./intrasearch/src/app/components/chat/interfaces/openai.structs.ts", openai.ChatCompletionResponse{}, openai.ChatCompletionRequest{})
}
func generate(fileName string, types ...any) {
	s2ts := struct2ts.New(&struct2ts.Options{
		// NoConstructor:    true,
		// NoToObject:       true,
		// NoHelpers:        true,
	})
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Printf("Generation TS structs in %s\n", fileName)
	for _, t := range types {
		fmt.Printf("  %T\n", t)
		s2ts.Add(t)
	}
	if err := s2ts.RenderTo(f); err != nil {
		fmt.Printf("Cannot render ts structs: %v", err)
	}
	//s2ts.RenderTo(os.Stdout)
}
