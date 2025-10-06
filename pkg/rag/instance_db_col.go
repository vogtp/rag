package rag

var _ (Instance) = (*instanceDBCol)(nil)

type instanceDBCol struct {
	*instanceCfg
}