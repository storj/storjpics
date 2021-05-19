// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package gallery

// Album represents one album in the gallery
type Album struct {
	// Name is the album's name.
	Name string
	// CoverImage is the file name of the picture to use as album's cover.
	CoverImage string
	// Pictures are the file names of all pictures in the album.
	Pictures []string
}
