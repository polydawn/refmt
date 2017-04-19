package atlas

type AtlasBuilder struct {
}

func (ab *AtlasBuilder) Complete() (Atlas, error) {
	return Atlas{}, nil
}

func (ab *AtlasBuilder) AddStructMap() {}
