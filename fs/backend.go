// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package fs

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"storj.io/storjpics/gallery"
)

// Backend for the local file system.
type Backend struct {
	rootDir string
}

// NewBackend creates new backend for the local system rooted at rootDir.
func NewBackend(rootDir string) *Backend {
	return &Backend{
		rootDir: rootDir,
	}
}

// GetAlbums returns all albums of the gallery.
func (fs *Backend) GetAlbums(ctx context.Context) ([]gallery.Album, error) {
	var albums []gallery.Album

	files, err := ioutil.ReadDir(filepath.Join(fs.rootDir, "pics", "original"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// ignore files
		if !file.IsDir() {
			continue
		}

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
			Name:       album,
			CoverImage: pictures[0],
			Pictures:   pictures,
		})
	}

	return albums, nil
}

// GetPictures returns all picture names in album.
func (fs *Backend) GetPictures(ctx context.Context, album string) ([]string, error) {
	var pictures []string

	files, err := ioutil.ReadDir(filepath.Join(fs.rootDir, "pics", "original", album))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// ignore subfolders
		if file.IsDir() {
			continue
		}

		// ignore hidden files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		pictures = append(pictures, file.Name())
	}

	return pictures, nil
}

// CreateFile creates a new file for writing.
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

// OpenFile opens a file for reading.
func (fs *Backend) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fs.rootDir, path))
}
