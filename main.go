// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"embed"
	"log"

	"github.com/spf13/cobra"
	"storj.io/storjpics/generator"
	"storj.io/storjpics/storj"
)

//go:embed site-template/*
var website embed.FS

// Flags contains different flags for commands.
type Flags struct {
	Access string
	Bucket string
}

func main() {
	var flags Flags

	rootCmd := &cobra.Command{
		Use:   "storjpics",
		Short: "Photo Gallery Generator for Storj DCS",
	}

	rootCmd.PersistentFlags().StringVar(&flags.Access, "access", "", "access grant to project on Storj DCS")
	rootCmd.PersistentFlags().StringVar(&flags.Bucket, "bucket", "", "bucket on Storj DCS")

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates the photo gallery website from scratch",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			checkFlags(flags)
			ctx := context.Background()
			backend, err := storj.NewBackend(ctx, flags.Access, flags.Bucket)
			if err != nil {
				return err
			}
			return generator.New(website, backend).Generate(ctx)
		},
	}

	rootCmd.AddCommand(
		generateCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func checkFlags(flags Flags) {
	if len(flags.Access) == 0 {
		log.Fatal("--access flag is required")
	}

	if len(flags.Bucket) == 0 {
		log.Fatal("--bucket flag is required")
	}
}
