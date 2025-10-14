
// helpers
const maxUnixTSInSeconds = 9999999999;

function ParseDate(d: Date | number | string): Date {
	if (d instanceof Date) return d;
	if (typeof d === 'number') {
		if (d > maxUnixTSInSeconds) return new Date(d);
		return new Date(d * 1000); // go ts
	}
	return new Date(d);
}

function ParseNumber(v: number | string, isInt = false): number {
	if (!v) return 0;
	if (typeof v === 'number') return v;
	return (isInt ? parseInt(v) : parseFloat(v)) || 0;
}

function FromArray<T>(Ctor: { new (v: any): T }, data?: any[] | any, def = null): T[] | null {
	if (!data || !Object.keys(data).length) return def;
	const d = Array.isArray(data) ? data : [data];
	return d.map((v: any) => new Ctor(v));
}

function ToObject(o: any, typeOrCfg: any = {}, child = false): any {
	if (o == null) return null;
	if (typeof o.toObject === 'function' && child) return o.toObject();

	switch (typeof o) {
		case 'string':
			return typeOrCfg === 'number' ? ParseNumber(o) : o;
		case 'boolean':
		case 'number':
			return o;
	}

	if (o instanceof Date) {
		return typeOrCfg === 'string' ? o.toISOString() : Math.floor(o.getTime() / 1000);
	}

	if (Array.isArray(o)) return o.map((v: any) => ToObject(v, typeOrCfg, true));

	const d: any = {};

	for (const k of Object.keys(o)) {
		const v: any = o[k];
		if (v === undefined) continue;
		if (v === null) continue;
		d[k] = ToObject(v, typeOrCfg[k] || {}, true);
	}

	return d;
}

// structs
// struct2ts:gorm.io/gorm.DeletedAt
class DeletedAt {
	Time: Date;
	Valid: boolean;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Time = ('Time' in d) ? ParseDate(d.Time) : new Date();
		this.Valid = ('Valid' in d) ? d.Valid as boolean : false;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.Time = 'string';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.SourceSystem
class SourceSystem {
	ID: number;
	CreatedAt: Date;
	UpdatedAt: Date;
	DeletedAt: DeletedAt;
	CollectionID: number;
	Name: string;
	Type: number;
	URL: string;
	Key: string;
	Parts: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.ID = ('ID' in d) ? d.ID as number : 0;
		this.CreatedAt = ('CreatedAt' in d) ? ParseDate(d.CreatedAt) : new Date();
		this.UpdatedAt = ('UpdatedAt' in d) ? ParseDate(d.UpdatedAt) : new Date();
		this.DeletedAt = new DeletedAt(d.DeletedAt);
		this.CollectionID = ('CollectionID' in d) ? d.CollectionID as number : 0;
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.Type = ('Type' in d) ? d.Type as number : 0;
		this.URL = ('URL' in d) ? d.URL as string : '';
		this.Key = ('Key' in d) ? d.Key as string : '';
		this.Parts = ('Parts' in d) ? d.Parts as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		cfg.ID = 'number';
		cfg.CreatedAt = 'string';
		cfg.UpdatedAt = 'string';
		cfg.CollectionID = 'number';
		cfg.Type = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.Collection
class Collection {
	ID: number;
	CreatedAt: Date;
	UpdatedAt: Date;
	DeletedAt: DeletedAt;
	UserID: number;
	DisplayName: string;
	CollectionName: string;
	APIKey: string;
	Source: SourceSystem;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.ID = ('ID' in d) ? d.ID as number : 0;
		this.CreatedAt = ('CreatedAt' in d) ? ParseDate(d.CreatedAt) : new Date();
		this.UpdatedAt = ('UpdatedAt' in d) ? ParseDate(d.UpdatedAt) : new Date();
		this.DeletedAt = new DeletedAt(d.DeletedAt);
		this.UserID = ('UserID' in d) ? d.UserID as number : 0;
		this.DisplayName = ('DisplayName' in d) ? d.DisplayName as string : '';
		this.CollectionName = ('CollectionName' in d) ? d.CollectionName as string : '';
		this.APIKey = ('APIKey' in d) ? d.APIKey as string : '';
		this.Source = new SourceSystem(d.Source);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.ID = 'number';
		cfg.CreatedAt = 'string';
		cfg.UpdatedAt = 'string';
		cfg.UserID = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg.User
class User {
	ID: number;
	CreatedAt: Date;
	UpdatedAt: Date;
	DeletedAt: DeletedAt;
	Name: string;
	APIKey: string;
	Collections: Collection[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.ID = ('ID' in d) ? d.ID as number : 0;
		this.CreatedAt = ('CreatedAt' in d) ? ParseDate(d.CreatedAt) : new Date();
		this.UpdatedAt = ('UpdatedAt' in d) ? ParseDate(d.UpdatedAt) : new Date();
		this.DeletedAt = new DeletedAt(d.DeletedAt);
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.APIKey = ('APIKey' in d) ? d.APIKey as string : '';
		this.Collections = Array.isArray(d.Collections) ? d.Collections.map((v: any) => new Collection(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.ID = 'number';
		cfg.CreatedAt = 'string';
		cfg.UpdatedAt = 'string';
		return ToObject(this, cfg);
	}
}

// exports
export {
	DeletedAt,
	SourceSystem,
	Collection,
	User,
	ParseDate,
	ParseNumber,
	FromArray,
	ToObject,
};
