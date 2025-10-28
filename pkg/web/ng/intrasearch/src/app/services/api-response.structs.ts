// structs
// struct2ts:github.com/vogtp/rag/pkg/web.QueryDocument
class QueryDocument {
	EmbedContent: string = '';
	Document: string = '';
	Modified: string = '';
	URL: string = '';
	Title: string = '';
	IDField: string = '';
	Distance: number = 0;
	UUID: number[] = [];
}

// struct2ts:github.com/vogtp/rag/pkg/web.CollectionSearchResponse
class CollectionSearchResponse {
	Collection: string = '';
	Query: string = '';
	Documents: QueryDocument[] | null = null;
}

// exports
export {
	QueryDocument,
	CollectionSearchResponse,
};
