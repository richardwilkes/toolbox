// Copyright (c) 2021-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package uti

import (
	"slices"
	"strings"
	"sync"
)

var (
	lock        sync.RWMutex
	byUTI       = make(map[string]*DataType)
	byMimeType  = make(map[string][]*DataType)
	byExtension = make(map[string][]*DataType)
)

// Some well-known Uniform Type Identifiers
var (
	Content = Register(&DataType{UTI: "public.content"})
	Item    = Register(&DataType{UTI: "public.item"})
	Data    = Register(&DataType{
		UTI:       "public.data",
		Parents:   []*DataType{Item},
		MimeTypes: []string{"application/octet-stream"},
	})
	Image = Register(&DataType{
		UTI:     "public.image",
		Parents: []*DataType{Data, Content},
	})
	Text = Register(&DataType{
		UTI:     "public.text",
		Parents: []*DataType{Data, Content},
	})
	PlainText = Register(&DataType{
		UTI:        "public.plain-text",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"text/plain"},
		Extensions: []string{".txt", ".text"},
	})
	UTF8PlainText = Register(&DataType{
		UTI:       "public.utf8-plain-text",
		Parents:   []*DataType{PlainText},
		MimeTypes: []string{"text/plain;charset=utf-8", `text/plain;charset="utf-8"`},
	})
	JSON = Register(&DataType{
		UTI:        "public.json",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"application/json"},
		Extensions: []string{".json"},
	})
	XML = Register(&DataType{
		UTI:        "public.xml",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"text/xml", "application/xml"},
		Extensions: []string{".xml"},
	})
	YAML = Register(&DataType{
		UTI:        "public.yaml",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"application/x-yaml"},
		Extensions: []string{".yaml", ".yml"},
	})
	Markdown = Register(&DataType{
		UTI:        "net.daringfireball.markdown",
		Parents:    []*DataType{PlainText},
		MimeTypes:  []string{"text/markdown"},
		Extensions: []string{".md", ".markdown"},
	})
	PDF = Register(&DataType{
		UTI:        "com.adobe.pdf",
		Parents:    []*DataType{Data},
		MimeTypes:  []string{"application/pdf", "application/x-pdf"},
		Extensions: []string{".pdf"},
	})
	BMP = Register(&DataType{
		UTI:        "com.microsoft.bmp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/bmp", "image/x-bmp"},
		Extensions: []string{".bmp", ".dib"},
	})
	GIF = Register(&DataType{
		UTI:        "com.compuserve.gif",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/gif"},
		Extensions: []string{".gif"},
	})
	ICNS = Register(&DataType{
		UTI:        "com.apple.icns",
		Parents:    []*DataType{Image},
		Extensions: []string{".icns"},
	})
	ICO = Register(&DataType{
		UTI:        "com.microsoft.ico",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/x-icon", "image/vnd.microsoft.icon"},
		Extensions: []string{".ico"},
	})
	JPEG = Register(&DataType{
		UTI:        "public.jpeg",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/jpeg", "image/jpg"},
		Extensions: []string{".jpeg", ".jpg", ".jpe", ".jif", ".jfif", ".jfi"},
	})
	PNG = Register(&DataType{
		UTI:        "public.png",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/png"},
		Extensions: []string{".png"},
	})
	SVG = Register(&DataType{
		UTI:        "public.svg-image",
		Parents:    []*DataType{Image, XML},
		MimeTypes:  []string{"image/svg+xml"},
		Extensions: []string{".svg", ".svgz"},
	})
	TIFF = Register(&DataType{
		UTI:        "public.tiff",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/tiff"},
		Extensions: []string{".tiff", ".tif"},
	})
	WBMP = Register(&DataType{
		UTI:        "com.adobe.wbmp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/vnd.wap.wbmp"},
		Extensions: []string{".wbmp"},
	})
	WEBP = Register(&DataType{
		UTI:        "org.webmproject.webp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/webp"},
		Extensions: []string{".webp"},
	})
	URL = Register(&DataType{
		UTI:     "public.url",
		Parents: []*DataType{Data},
	})
	FileURL = Register(&DataType{
		UTI:     "public.file-url",
		Parents: []*DataType{URL},
	})
)

// DataType holds information about a single Uniform Type Identifier (UTI) and its relationship to other UTI's, MIME
// types, and file extensions. Once registered, no fields within this struct should be modified without first
// unregistering it.
type DataType struct {
	UTI        string
	Parents    []*DataType
	MimeTypes  []string
	Extensions []string
}

// Register the DataType. If the UTI has been previously registered, the old one will first be unregistered and then the
// new one will be registered.
func Register(dataType *DataType) *DataType {
	uti := strings.ToLower(dataType.UTI)
	lock.Lock()
	defer lock.Unlock()
	unregister(uti)
	byUTI[uti] = dataType
	for _, mimeType := range dataType.MimeTypes {
		mimeType = strings.ToLower(mimeType)
		byMimeType[mimeType] = append(byMimeType[mimeType], dataType)
	}
	for _, extension := range dataType.Extensions {
		extension = strings.ToLower(extension)
		byExtension[extension] = append(byExtension[extension], dataType)
	}
	return dataType
}

// Unregister the DataType.
func Unregister(dataType *DataType) {
	uti := strings.ToLower(dataType.UTI)
	lock.Lock()
	defer lock.Unlock()
	unregister(uti)
}

func unregister(uti string) {
	if existing, ok := byUTI[uti]; ok {
		filter := func(one *DataType) bool {
			return one == existing
		}
		for _, mimeType := range existing.MimeTypes {
			mimeType = strings.ToLower(mimeType)
			byMimeType[mimeType] = slices.DeleteFunc(byMimeType[mimeType], filter)
		}
		for _, extension := range existing.Extensions {
			extension = strings.ToLower(extension)
			byExtension[extension] = slices.DeleteFunc(byExtension[extension], filter)
		}
		delete(byUTI, uti)
	}
}

// ByUTI looks up the DataType by its UTI.
func ByUTI(uti string) *DataType {
	lock.RLock()
	defer lock.RUnlock()
	return byUTI[strings.ToLower(uti)]
}

// ByMimeType looks up DataTypes that use the given MIME type.
func ByMimeType(mimeType string) []*DataType {
	lock.RLock()
	defer lock.RUnlock()
	return byMimeType[strings.ToLower(mimeType)]
}

// ByExtension looks up DataTypes that use the given file extension.
func ByExtension(extension string) []*DataType {
	lock.RLock()
	defer lock.RUnlock()
	return byExtension[strings.ToLower(extension)]
}

// ConformsTo returns true if this DataType is the same as or a descendant of the target DataType.
func (dt *DataType) ConformsTo(target *DataType) bool {
	if dt.UTI == target.UTI {
		return true
	}
	for _, parent := range dt.Parents {
		if parent.ConformsTo(target) {
			return true
		}
	}
	return false
}
