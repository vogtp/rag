
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
// struct2ts:github.com/vogtp/rag/pkg/usercfg/db/ent.SourceSystem
class SourceSystem {
	id: number;
	Name: string;
	Type: string;
	URL: string;
	key: string;
	parts: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.id = ('id' in d) ? d.id as number : 0;
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.Type = ('Type' in d) ? d.Type as string : '';
		this.URL = ('URL' in d) ? d.URL as string : '';
		this.key = ('key' in d) ? d.key as string : '';
		this.parts = ('parts' in d) ? d.parts as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		cfg.id = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg/db/ent.CollectionEdges
class CollectionEdges {
	Sources: SourceSystem[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Sources = Array.isArray(d.Sources) ? d.Sources.map((v: any) => new SourceSystem(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg/db/ent.Collection
class Collection {
	id: number;
	Name: string;
	APIKey: string;
	collectionName: string;
	edges: CollectionEdges;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.id = ('id' in d) ? d.id as number : 0;
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.APIKey = ('APIKey' in d) ? d.APIKey as string : '';
		this.collectionName = ('collectionName' in d) ? d.collectionName as string : '';
		this.edges = new CollectionEdges(d.edges);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.id = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg/db/ent.UserEdges
class UserEdges {
	Collections: Collection[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Collections = Array.isArray(d.Collections) ? d.Collections.map((v: any) => new Collection(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/vogtp/rag/pkg/usercfg/db/ent.User
class User {
	id: number;
	Name: string;
	APIKey: string;
	edges: UserEdges;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.id = ('id' in d) ? d.id as number : 0;
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.APIKey = ('APIKey' in d) ? d.APIKey as string : '';
		this.edges = new UserEdges(d.edges);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.id = 'number';
		return ToObject(this, cfg);
	}
}

// exports
export {
	SourceSystem,
	CollectionEdges,
	Collection,
	UserEdges,
	User,
	ParseDate,
	ParseNumber,
	FromArray,
	ToObject,
};
