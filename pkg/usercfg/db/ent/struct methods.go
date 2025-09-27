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

func (cc *ConfluenceCreate) SetConfluence(input *Confluence) *ConfluenceCreate {
	cc.SetName(input.Name)
	cc.SetURL(input.URL)
	cc.SetConfluenceAPIKey(input.ConfluenceAPIKey)
	return cc
}

func (cu *ConfluenceUpdate) SetConfluence(input *Confluence) *ConfluenceUpdate {
	cu.SetName(input.Name)
	cu.SetURL(input.URL)
	cu.SetConfluenceAPIKey(input.ConfluenceAPIKey)
	return cu
}

func (cuo *ConfluenceUpdateOne) SetConfluence(input *Confluence) *ConfluenceUpdateOne {
	cuo.SetName(input.Name)
	cuo.SetURL(input.URL)
	cuo.SetConfluenceAPIKey(input.ConfluenceAPIKey)
	return cuo
}

func (sc *SpaceCreate) SetSpace(input *Space) *SpaceCreate {
	sc.SetName(input.Name)
	sc.SetSpaceKey(input.SpaceKey)
	return sc
}

func (su *SpaceUpdate) SetSpace(input *Space) *SpaceUpdate {
	su.SetName(input.Name)
	su.SetSpaceKey(input.SpaceKey)
	return su
}

func (suo *SpaceUpdateOne) SetSpace(input *Space) *SpaceUpdateOne {
	suo.SetName(input.Name)
	suo.SetSpaceKey(input.SpaceKey)
	return suo
}

func (uc *UserCreate) SetUser(input *User) *UserCreate {
	uc.SetName(input.Name)
	uc.SetOpenaiAPIkey(input.OpenaiAPIkey)
	return uc
}

func (uu *UserUpdate) SetUser(input *User) *UserUpdate {
	uu.SetName(input.Name)
	uu.SetOpenaiAPIkey(input.OpenaiAPIkey)
	return uu
}

func (uuo *UserUpdateOne) SetUser(input *User) *UserUpdateOne {
	uuo.SetName(input.Name)
	uuo.SetOpenaiAPIkey(input.OpenaiAPIkey)
	return uuo
}
