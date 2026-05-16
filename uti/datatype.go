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

// Some well-known Uniform Type Identifiers
var (
	BMP           *DataType
	Content       *DataType
	Data          *DataType
	GIF           *DataType
	ICNS          *DataType
	ICO           *DataType
	Image         *DataType
	Item          *DataType
	JPEG          *DataType
	JSON          *DataType
	PlainText     *DataType
	PNG           *DataType
	SVG           *DataType
	Text          *DataType
	TIFF          *DataType
	UTF8PlainText *DataType
	WBMP          *DataType
	WEBP          *DataType
	XML           *DataType
	YAML          *DataType
)

var (
	lock        sync.RWMutex
	byUTI       = make(map[string]*DataType)
	byMimeType  = make(map[string][]*DataType)
	byExtension = make(map[string][]*DataType)
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

func init() {
	Content = register(&DataType{UTI: "public.content"})
	Image = register(&DataType{UTI: "public.image"})
	Item = register(&DataType{UTI: "public.item"})

	Data = register(&DataType{
		UTI:       "public.data",
		Parents:   []*DataType{Item},
		MimeTypes: []string{"application/octet-stream"},
	})

	Text = register(&DataType{
		UTI:     "public.text",
		Parents: []*DataType{Data, Content},
	})
	PlainText = register(&DataType{
		UTI:        "public.plain-text",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"text/plain"},
		Extensions: []string{".txt", ".text"},
	})
	UTF8PlainText = register(&DataType{
		UTI:       "public.utf8-plain-text",
		Parents:   []*DataType{PlainText},
		MimeTypes: []string{"text/plain;charset=utf-8", `text/plain;charset="utf-8"`},
	})
	JSON = register(&DataType{
		UTI:        "public.json",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"application/json"},
		Extensions: []string{".json"},
	})
	XML = register(&DataType{
		UTI:        "public.xml",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"text/xml", "application/xml"},
		Extensions: []string{".xml"},
	})
	YAML = register(&DataType{
		UTI:        "public.yaml",
		Parents:    []*DataType{Text},
		MimeTypes:  []string{"application/x-yaml"},
		Extensions: []string{".yaml", ".yml"},
	})

	BMP = register(&DataType{
		UTI:        "com.microsoft.bmp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/bmp", "image/x-bmp"},
		Extensions: []string{".bmp", ".dib"},
	})
	GIF = register(&DataType{
		UTI:        "com.compuserve.gif",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/gif"},
		Extensions: []string{".gif"},
	})
	ICNS = register(&DataType{
		UTI:        "com.apple.icns",
		Parents:    []*DataType{Image},
		Extensions: []string{".icns"},
	})
	ICO = register(&DataType{
		UTI:        "com.microsoft.ico",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/x-icon", "image/vnd.microsoft.icon"},
		Extensions: []string{".ico"},
	})
	JPEG = register(&DataType{
		UTI:        "public.jpeg",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/jpeg", "image/jpg"},
		Extensions: []string{"jpeg", ".jpg", "jpe", ".jif", ".jfif", ".jfi"},
	})
	PNG = register(&DataType{
		UTI:        "public.png",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/png"},
		Extensions: []string{".png"},
	})
	SVG = register(&DataType{
		UTI:        "public.svg-image",
		Parents:    []*DataType{Image, XML},
		MimeTypes:  []string{"image/svg+xml"},
		Extensions: []string{".svg", ".svgz"},
	})
	TIFF = register(&DataType{
		UTI:        "public.tiff",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/tiff"},
		Extensions: []string{".tiff", ".tif"},
	})
	WBMP = register(&DataType{
		UTI:        "com.adobe.wbmp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/vnd.wap.wbmp"},
		Extensions: []string{".wbmp"},
	})
	WEBP = register(&DataType{
		UTI:        "org.webmproject.webp",
		Parents:    []*DataType{Image},
		MimeTypes:  []string{"image/webp"},
		Extensions: []string{".webp"},
	})
}

func register(dataType *DataType) *DataType {
	Register(dataType)
	return dataType
}

// Register the DataType. If the UTI has been previously registered, the old one will first be unregistered and then the
// new one will be registered.
func Register(dataType *DataType) {
	uti := strings.ToLower(dataType.UTI)
	lock.Lock()
	defer lock.Unlock()
	unregister(uti)
	byUTI[uti] = dataType
	for _, mimeType := range dataType.MimeTypes {
		mimeType = strings.ToLower(mimeType)
		byMimeType[mimeType] = append(byMimeType[mimeType], dataType)
	}
	for _, extension := range dataType.MimeTypes {
		extension = strings.ToLower(extension)
		byExtension[extension] = append(byExtension[extension], dataType)
	}
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
		for _, extension := range existing.MimeTypes {
			extension = strings.ToLower(extension)
			byExtension[extension] = slices.DeleteFunc(byExtension[extension], filter)
		}
		delete(byUTI, uti)
	}
}

// ByUTI looks up the DataType by its UTI.
func ByUTI(uti string) *DataType {
	return byUTI[strings.ToLower(uti)]
}

// ByMimeType looks up DataTypes that use the given MIME type.
func ByMimeType(mimeType string) []*DataType {
	return byMimeType[strings.ToLower(mimeType)]
}

// ByExtension looks up DataTypes that use the given file extension.
func ByExtension(extension string) []*DataType {
	return byExtension[strings.ToLower(extension)]
}
