// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import (
	"context"
	"io"
	"path"
	"sort"
	"strings"

	"storj.io/storjpics/gallery"
	"storj.io/uplink"
)

// Backend for Storj DCS.
type Backend struct {
	project *uplink.Project
	bucket  string
}

// NewBackend creates new backend for Storj DCS for bucket using accessGrant.
func NewBackend(ctx context.Context, accessGrant, bucket string) (*Backend, error) {
	access, err := uplink.ParseAccess(accessGrant)
	if err != nil {
		return nil, err
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, err
	}

	return &Backend{
		project: project,
		bucket:  bucket,
	}, nil
}

// GetAlbums returns all albums of the gallery.
func (storj *Backend) GetAlbums(ctx context.Context) ([]gallery.Album, error) {
	var albums []gallery.Album

	iterator := storj.project.ListObjects(ctx, storj.bucket, &uplink.ListObjectsOptions{
		Prefix: "pics/original/",
	})

	for iterator.Next() {
		item := iterator.Item()

		// ignore files
		if !item.IsPrefix {
			continue
		}

		album := path.Base(item.Key)
		pictures, err := storj.GetPictures(ctx, album)
		if err != nil {
			return nil, err
		}

		// skip empty albums
		if len(pictures) == 0 {
			continue
		}

		albums = append(albums, gallery.Album{
			Name:       album,
			CoverImage: pictures[0],
			Pictures:   pictures,
		})
	}
	if iterator.Err() != nil {
		return nil, iterator.Err()
	}

	// ensure alphabetical sort order
	sort.Slice(albums, func(i, k int) bool {
		return albums[i].Name < albums[k].Name
	})

	return albums, nil
}

// GetPictures returns all picture names in album.
func (storj *Backend) GetPictures(ctx context.Context, album string) ([]string, error) {
	var pictures []string

	iterator := storj.project.ListObjects(ctx, storj.bucket, &uplink.ListObjectsOptions{
		Prefix: path.Join("pics/original", album) + "/",
	})

	for iterator.Next() {
		item := iterator.Item()

		// ignore subfolders
		if item.IsPrefix {
			continue
		}

		picture := path.Base(item.Key)

		// ignore hidden files
		if strings.HasPrefix(picture, ".") {
			continue
		}

		pictures = append(pictures, picture)
	}
	if iterator.Err() != nil {
		return nil, iterator.Err()
	}

	// ensure alphabetical sort order
	sort.Strings(pictures)

	return pictures, nil
}

// CreateFile creates a new file for writing.
func (storj *Backend) CreateFile(ctx context.Context, path string) (io.WriteCloser, error) {
	upload, err := storj.project.UploadObject(ctx, storj.bucket, path, nil)
	if err != nil {
		return nil, err
	}

	return &uploadCloser{
		upload: upload,
	}, nil
}

// OpenFile opens a file for reading.
func (storj *Backend) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return storj.project.DownloadObject(ctx, storj.bucket, path, nil)
}

// uploadCloser implements io.WriteCloser for uplink.Upload.
type uploadCloser struct {
	upload *uplink.Upload
}

func (uc uploadCloser) Write(p []byte) (n int, err error) {
	return uc.upload.Write(p)
}

func (uc uploadCloser) Close() error {
	return uc.upload.Commit()
}
