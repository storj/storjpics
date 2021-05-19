// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"embed"
	"log"

	"github.com/spf13/cobra"
	"storj.io/storjpics/fs"
	"storj.io/storjpics/generator"
)

//go:embed site-template/*
var website embed.FS

func main() {
	rootCmd := &cobra.Command{
		Use:   "storjpics",
		Short: "Photo Gallery Generator for Storj DCS",
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generates the photo gallery website from scratch",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return generator.New(
				website,
				fs.NewBackend("/tmp/storjpics"), // TODO: replace with Storj backend
			).Generate(context.Background())
		},
	}

	rootCmd.AddCommand(
		generateCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
