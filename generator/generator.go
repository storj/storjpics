// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package generator

import (
	"context"
	"embed"
	"fmt"
	"image"
	"io"
	"log"
	"path"
	"text/template"

	"github.com/disintegration/imaging"
	"storj.io/storjpics/gallery"
)

// Generator is a photo gallery generator.
type Generator struct {
	assets  embed.FS
	backend Backend
}

// Backend is a backend implementation to store the photo gallery assets.
type Backend interface {
	// GetAlbums returns all albums of the gallery.
	GetAlbums(ctx context.Context) ([]gallery.Album, error)
	// GetPictures returns all picture names in album.
	GetPictures(ctx context.Context, album string) ([]string, error)
	// CreateFile creates a new file for writing.
	CreateFile(ctx context.Context, path string) (io.WriteCloser, error)
	// OpenFile opens a file for reading.
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
}

// New creates a new generator.
func New(assets embed.FS, backend Backend) *Generator {
	return &Generator{
		assets:  assets,
		backend: backend,
	}
}

// Generate generates all photo gallery assets from pictures available at /pics/original.
func (generator *Generator) Generate(ctx context.Context) error {
	err := generator.copyAssetsDir(ctx, "site-template/homepage/assets", "assets/homepage")
	if err != nil {
		return err
	}

	err = generator.copyAssetsDir(ctx, "site-template/album/assets", "assets/album")
	if err != nil {
		return err
	}

	log.Print("Listing albums...")
	albums, err := generator.backend.GetAlbums(ctx)
	if err != nil {
		return err
	}

	page, err := generator.assets.ReadFile("site-template/album/index.html")
	if err != nil {
		return err
	}

	t, err := template.New("album").Parse(string(page))
	if err != nil {
		return err
	}

	for _, album := range albums {
		err = generator.createResizedImages(ctx, album)
		if err != nil {
			return err
		}

		data := struct {
			AlbumName string
			Pictures  []string
		}{
			AlbumName: album.Name,
			Pictures:  album.Pictures,
		}

		err = generator.copyHTMLFile(ctx, t, data, path.Join(album.Name, "index.html"))
		if err != nil {
			return err
		}
	}

	page, err = generator.assets.ReadFile("site-template/homepage/index.html")
	if err != nil {
		return err
	}

	t, err = template.New("homepage").Parse(string(page))
	if err != nil {
		return err
	}

	data := struct {
		Title  string
		Albums []gallery.Album
	}{
		// TODO: make the title configurable
		Title:  "Photo gallery hosted on Storj DCS",
		Albums: albums,
	}

	return generator.copyHTMLFile(ctx, t, data, "index.html")
}

func (generator *Generator) copyHTMLFile(ctx context.Context, t *template.Template, data interface{}, dest string) error {
	log.Printf("Copying HTML file %s...", dest)

	file, err := generator.backend.CreateFile(ctx, dest)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, data)
}

func (generator *Generator) copyAssetsDir(ctx context.Context, assetsDir, destDir string) error {
	log.Printf("Copying assets to %s...", destDir)

	files, err := generator.assets.ReadDir(assetsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = generator.copyAssetsDir(ctx, path.Join(assetsDir, file.Name()), path.Join(destDir, file.Name()))
			if err != nil {
				return err
			}
			continue
		}

		err = generator.copyAssetsFile(ctx, path.Join(assetsDir, file.Name()), path.Join(destDir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (generator *Generator) copyAssetsFile(ctx context.Context, src, dest string) error {
	log.Printf("Copying asset file %s...", dest)

	destFile, err := generator.backend.CreateFile(ctx, dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := generator.assets.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (generator *Generator) createResizedImages(ctx context.Context, album gallery.Album) error {
	for _, picture := range album.Pictures {
		reader, err := generator.backend.OpenFile(ctx, path.Join("pics", "original", album.Name, picture))
		if err != nil {
			return err
		}
		defer reader.Close()

		original, err := imaging.Decode(reader)
		if err != nil {
			return err
		}

		format, err := imaging.FormatFromFilename(picture)
		if err != nil {
			return err
		}

		small := imaging.Thumbnail(original, 360, 225, imaging.Lanczos)
		err = generator.copyResizedImage(ctx, small, format, album.Name, picture, 360, 225)
		if err != nil {
			return err
		}

		// set width to 0 to preserve aspect ratio
		big := imaging.Resize(original, 0, 750, imaging.Lanczos)
		err = generator.copyResizedImage(ctx, big, format, album.Name, picture, 1200, 750)
		if err != nil {
			return err
		}
	}

	return nil
}

func (generator *Generator) copyResizedImage(ctx context.Context, resized image.Image, format imaging.Format, album, picture string, width, height int) error {
	imagePath := path.Join("pics", "resized", fmt.Sprintf("%dx%d", width, height), album, picture)

	log.Printf("Copying image file to %s...", imagePath)

	writer, err := generator.backend.CreateFile(ctx, imagePath)
	if err != nil {
		return err
	}

	return imaging.Encode(writer, resized, format)
}
