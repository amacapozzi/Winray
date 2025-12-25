package models

type FileResult struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Kind           string `json:"kind,omitempty"`
	MetaLeft       string `json:"metaLeft,omitempty"`
	MetaRight      string `json:"metaRight,omitempty"`
	LastAccessTime int64  `json:"lastAccessTime,omitempty"`
}

type IndexedFile struct {
	Path           string
	LastAccessTime int64
}
