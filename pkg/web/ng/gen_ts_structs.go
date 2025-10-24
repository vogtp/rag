//go:build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/OneOfOne/struct2ts"
	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/web"
)

func main() {
	generate("./intrasearch/src/app/services/user.structs.ts", usercfg.User{})
	generate("./intrasearch/src/app/services/api-response.structs.ts", web.CollectionSearchResponse{})
	generate("./intrasearch/src/app/components/chat/interfaces/openai.structs.ts", openai.ChatCompletionResponse{}, openai.ChatCompletionRequest{})
	version := fmt.Sprintf("%v.%v.%v (%v)", cfg.VersionMajor, cfg.VersionMinor, cfg.VersionPatch, time.Now().Format("2006-01-02T15:04:05"))

	injectStrings("./intrasearch/src/app/go.transfer.ts", "version", version)
}
func generate(fileName string, types ...any) {
	s2ts := struct2ts.New(&struct2ts.Options{
		// NoConstructor:    true,
		// NoToObject:       true,
		NoHelpers:        true,
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
}

func injectStrings(fileName string, strs ...string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Printf("Inject go strings to TS in %s\n", fileName)
	for i := 0; i < len(strs); i += 2 {
		fmt.Printf("  %s\n", strs[i])
		fmt.Fprintf(f, `export const %s: string = "%s"`, strs[i], strs[i+1])
	}
}
