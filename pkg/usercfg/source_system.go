package usercfg

//go:generate stringer -type=SourceSystemType --trimprefix Source

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type SourceSystemType int

const (
	SourceConfluence SourceSystemType = iota
	SourceHTTP
)

type SourceSystem struct {
	// gorm.Model in order to avoid name clashes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CollectionID uint             `gorm:"column:collection_id"`
	Name         string           `gorm:"column:name"`
	Type         SourceSystemType `gorm:"column:src_sys_type"`
	URL          string           `gorm:"column:url"`
	Key          string           `gorm:"column:key"`
	Parts        string           `gorm:"column:parts"`
}

func (s SourceSystem) splitParts() []string {
	parts := strings.Split(s.Parts, ",")
	if len(parts) < 1 {
		parts = strings.Split(s.Parts, " ")
	}
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
