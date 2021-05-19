// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package gallery

// Album represents one album in the gallery
type Album struct {
	Name            string
	PictureFileName string
	Pictures        []string
}
