package storage

import (
	"io"
	"time"
)

type Store interface {
	// Get returns the content of the file at the given path.
	PathType(path string) PathType
	List(path string) ([]Entry, error)
	File(path string) *File
	Cover(path string) *File
	Thumbnail(path string) *File
}

type PathType string

const (
	PathTypeNotExists  PathType = "not_exists"
	PathTypeFile       PathType = "file"
	PathTypeNavigation PathType = "application/atom+xml;profile=opds-catalog;kind=navigation"
	PathTypeAquisition PathType = "application/atom+xml;profile=opds-catalog;kind=acquisition"
)

type File struct {
	ContentType   string
	ContentLength int64
	Reader        io.ReadCloser
}

type Entry struct {
	Name       string
	Type       string
	Aquisition string
	Updated    time.Time
	Metadata   Metadata
}

type Metadata interface {
	GetTitle() string
	GetLanguage() string
	GetIdentifier() string
	GetCreator() string
	GetContributor() string
	GetPublisher() string
	GetSubject() string
	GetDescription() string
	HasCover() bool
	HasThumbnail() bool
}

type NOOPMetadata struct{}

func (NOOPMetadata) GetTitle() string       { return "" }
func (NOOPMetadata) GetLanguage() string    { return "" }
func (NOOPMetadata) GetIdentifier() string  { return "" }
func (NOOPMetadata) GetCreator() string     { return "" }
func (NOOPMetadata) GetContributor() string { return "" }
func (NOOPMetadata) GetPublisher() string   { return "" }
func (NOOPMetadata) GetSubject() string     { return "" }
func (NOOPMetadata) GetDescription() string { return "" }
func (NOOPMetadata) HasCover() bool         { return false }
func (NOOPMetadata) HasThumbnail() bool     { return false }
