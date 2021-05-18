// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package fs

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"storj.io/storjpics/gallery"
)

type Backend struct {
	rootDir string
}

func NewBackend(rootDir string) *Backend {
	return &Backend{
		rootDir: rootDir,
	}
}

func (fs *Backend) GetAlbums(ctx context.Context) ([]gallery.Album, error) {
	files, err := ioutil.ReadDir(filepath.Join(fs.rootDir, "pics", "original"))
	if err != nil {
		return nil, err
	}

	var albums []gallery.Album

	for _, file := range files {
		if file.IsDir() {
			album := file.Name()
			pictures, err := fs.GetPictures(ctx, album)
			if err != nil {
				return nil, err
			}
			// skip empty albums
			if len(pictures) == 0 {
				continue
			}
			albums = append(albums, gallery.Album{
				Name:            file.Name(),
				PictureFileName: pictures[0],
				Pictures:        pictures,
			})
		}
	}

	return albums, nil
}

func (fs *Backend) GetPictures(ctx context.Context, album string) ([]string, error) {
	files, err := ioutil.ReadDir(filepath.Join(fs.rootDir, "pics", "original", album))
	if err != nil {
		return nil, err
	}

	var pictures []string
	for _, file := range files {
		if !file.IsDir() {
			pictures = append(pictures, file.Name())
		}
	}

	return pictures, nil
}

func (fs *Backend) CreateFile(ctx context.Context, path string) (io.WriteCloser, error) {
	err := fs.ensureParentDir(path)
	if err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(fs.rootDir, path))
}

func (fs *Backend) ensureParentDir(path string) error {
	file := filepath.Join(fs.rootDir, path)
	dir := filepath.Dir(file)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *Backend) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fs.rootDir, path))
}
