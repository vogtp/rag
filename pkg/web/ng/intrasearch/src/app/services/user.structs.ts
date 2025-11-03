// structs
// struct2ts:gorm.io/gorm.DeletedAt
class DeletedAt {
	Time: Date = new Date();
	Valid: boolean = false;
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.SourceSystem
class SourceSystem {
	ID: number = 0;
	CreatedAt: Date = new Date();
	UpdatedAt: Date = new Date();
	DeletedAt: DeletedAt = new DeletedAt();
	CollectionID: number = 0;
	Name: string = '';
	Type: number = 0;
	URL: string = '';
	Key: string = '';
	Parts: string = '';
	QueryRetryMax: number = 0;
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.Collection
class Collection {
	ID: number = 0;
	CreatedAt: Date = new Date();
	UpdatedAt: Date = new Date();
	DeletedAt: DeletedAt = new DeletedAt();
	UserID: number = 0;
	Displayname: string = '';
	Collectionname: string = '';
	APIKey: string = '';
	Genmodel: string = '';
	Embedmodel: string = '';
	DBUpdateIntervall: number = 0;
	NextDBUpdate: Date = new Date();
	StartDBUpdate: Date = new Date();
	Source: SourceSystem = new SourceSystem();
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.User
class User {
	ID: number = 0;
	CreatedAt: Date = new Date();
	UpdatedAt: Date = new Date();
	DeletedAt: DeletedAt = new DeletedAt();
	Name: string = '';
	APIKey: string = '';
	Collections: Collection[] | null = null;
	AdvancedUI: boolean = false;
}

// exports
export {
	DeletedAt,
	SourceSystem,
	Collection,
	User,
};
