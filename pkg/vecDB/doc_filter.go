package vecdb

type Filter interface {
	// ShouldEmbedd returns true if the document should be embedded
	ShouldEmbedd(*EmbeddDocument) bool
	// ReqisterEmedded tels the Filter that a document has been embedded
	ReqisterEmedded(*EmbeddDocument)
}
 

