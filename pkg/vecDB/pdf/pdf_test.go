package pdf_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/vogtp/langchaingo/documentloaders"
	"github.com/vogtp/langchaingo/textsplitter"
	"github.com/vogtp/rag/pkg/vecDB/pdf"
)

func Test_confluence_embeddPDF(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		link string
	}{
		{name: "pikett", link: "https://www.unibas.ch/dam/jcr:7fd58714-ee1d-4baa-b1cd-48793d0cd1f5/R_Pikett_Inkonvenienzentsch%C3%A4digung_00.pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := pdf.SplitFromLink(context.Background(), tt.link)
			if err != nil {
				t.Error(err)
			}
			if len(docs) < 10 {
				t.Errorf("not enought documents %v", len(docs))
			}
			for _, d := range docs {
				fmt.Printf("%v\n***\n", d.Document)
			}
			// t.Fail()
		})
	}
}

func Test_embedd_from_file(t *testing.T) {
	pdfName := "/tmp/3951918067R_Pikett_Inkonvenienzentsch%C3%A4digung_00.pdf"

	docs, err := pdf.SplitFromFile(context.Background(), pdfName)
	if err != nil {
		t.Error(err)
	}
	for _, d := range docs {
		fmt.Printf("%v\n***\n", d.Document)
	}
	if len(docs) < 10 {
		t.Errorf("not enought documents %v", len(docs))
	}
	// t.Fail()
}
func Test_pdf_docloader(t *testing.T) {
	pdfName := "/tmp/3951918067R_Pikett_Inkonvenienzentsch%C3%A4digung_00.pdf"
	f, err := os.Open(pdfName)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	pl := documentloaders.NewPDF(f, fi.Size())

	split := textsplitter.NewRecursiveCharacter()

	docs, err := pl.LoadAndSplit(context.Background(), split)
	if err != nil {
		t.Error(err)
	}
	if len(docs) < 10 {
		t.Errorf("not enought documents %v", len(docs))
	}
	for _, d := range docs {
		fmt.Printf("%v\n***\n", d.PageContent)
	}
	// t.Fail()
}
