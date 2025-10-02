package ent

func (cc *CollectionCreate) SetCollection(input *Collection) *CollectionCreate {
	cc.SetName(input.Name)
	cc.SetAPIKey(input.APIKey)
	return cc
}

func (cu *CollectionUpdate) SetCollection(input *Collection) *CollectionUpdate {
	cu.SetName(input.Name)
	cu.SetAPIKey(input.APIKey)
	return cu
}

func (cuo *CollectionUpdateOne) SetCollection(input *Collection) *CollectionUpdateOne {
	cuo.SetName(input.Name)
	cuo.SetAPIKey(input.APIKey)
	return cuo
}

func (ssc *SourceSystemCreate) SetSourceSystem(input *SourceSystem) *SourceSystemCreate {
	ssc.SetName(input.Name)
	ssc.SetType(input.Type)
	ssc.SetURL(input.URL)
	ssc.SetKey(input.Key)
	ssc.SetParts(input.Parts)
	return ssc
}

func (ssu *SourceSystemUpdate) SetSourceSystem(input *SourceSystem) *SourceSystemUpdate {
	ssu.SetName(input.Name)
	ssu.SetType(input.Type)
	ssu.SetURL(input.URL)
	ssu.SetKey(input.Key)
	ssu.SetParts(input.Parts)
	return ssu
}

func (ssuo *SourceSystemUpdateOne) SetSourceSystem(input *SourceSystem) *SourceSystemUpdateOne {
	ssuo.SetName(input.Name)
	ssuo.SetType(input.Type)
	ssuo.SetURL(input.URL)
	ssuo.SetKey(input.Key)
	ssuo.SetParts(input.Parts)
	return ssuo
}

func (uc *UserCreate) SetUser(input *User) *UserCreate {
	uc.SetName(input.Name)
	uc.SetAPIKey(input.APIKey)
	return uc
}

func (uu *UserUpdate) SetUser(input *User) *UserUpdate {
	uu.SetName(input.Name)
	uu.SetAPIKey(input.APIKey)
	return uu
}

func (uuo *UserUpdateOne) SetUser(input *User) *UserUpdateOne {
	uuo.SetName(input.Name)
	uuo.SetAPIKey(input.APIKey)
	return uuo
}
