package ent

func (_c *CollectionCreate) SetCollection(input *Collection) *CollectionCreate {
	_c.SetName(input.Name)
	_c.SetAPIKey(input.APIKey)
	return _c
}

func (_u *CollectionUpdate) SetCollection(input *Collection) *CollectionUpdate {
	_u.SetName(input.Name)
	_u.SetAPIKey(input.APIKey)
	return _u
}

func (_u *CollectionUpdateOne) SetCollection(input *Collection) *CollectionUpdateOne {
	_u.SetName(input.Name)
	_u.SetAPIKey(input.APIKey)
	return _u
}

func (_c *SourceSystemCreate) SetSourceSystem(input *SourceSystem) *SourceSystemCreate {
	_c.SetName(input.Name)
	_c.SetType(input.Type)
	_c.SetURL(input.URL)
	_c.SetKey(input.Key)
	_c.SetParts(input.Parts)
	return _c
}

func (_u *SourceSystemUpdate) SetSourceSystem(input *SourceSystem) *SourceSystemUpdate {
	_u.SetName(input.Name)
	_u.SetType(input.Type)
	_u.SetURL(input.URL)
	_u.SetKey(input.Key)
	_u.SetParts(input.Parts)
	return _u
}

func (_u *SourceSystemUpdateOne) SetSourceSystem(input *SourceSystem) *SourceSystemUpdateOne {
	_u.SetName(input.Name)
	_u.SetType(input.Type)
	_u.SetURL(input.URL)
	_u.SetKey(input.Key)
	_u.SetParts(input.Parts)
	return _u
}

func (_c *UserCreate) SetUser(input *User) *UserCreate {
	_c.SetName(input.Name)
	_c.SetAPIKey(input.APIKey)
	return _c
}

func (_u *UserUpdate) SetUser(input *User) *UserUpdate {
	_u.SetName(input.Name)
	_u.SetAPIKey(input.APIKey)
	return _u
}

func (_u *UserUpdateOne) SetUser(input *User) *UserUpdateOne {
	_u.SetName(input.Name)
	_u.SetAPIKey(input.APIKey)
	return _u
}
