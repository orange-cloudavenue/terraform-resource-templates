/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/orange-cloudavenue/terraform-resource-templates/internal/terraform"
	"github.com/orange-cloudavenue/terraform-resource-templates/pkg/file"

	_ "embed"
)

var version = "v0.0.0-dev"

func main() {
	fileName := flag.String("filename", "", "filename")
	update := flag.Bool("update", false, "update the program")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *update {
		err := selfUpdate()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	checkVersion()

	if *fileName == "" {
		log.Fatal().Msg("filename is required")
	}

	if !file.IsFileExists(*fileName) {
		log.Fatal().Msgf("file %s not found", *fileName)
	}

	extractVarNaming()

	// Get the absolute path of the file
	absPath, err := filepath.Abs(*fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error getting absolute path")
	}

	// Set test dir
	testDir := filepath.Join(absPath, "../../../testsacc/")

	// Determine the schema Path with the absolute path
	schemaDir := filepath.Join(absPath, "../")

	log.Info().Msgf("using schema dir %s", schemaDir)

	// test if filedir exists and is a directory
	dir, err := os.Stat(testDir)
	if err != nil || !dir.IsDir() {
		// create the directory
		if err := os.MkdirAll(testDir, 0o755); err != nil {
			log.Fatal().Err(err).Msgf("error creating directory %s", testDir)
		}
	}

	log.Info().Msgf("using file %s", *fileName)

	f, err := file.ToString(*fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading file")
	}

	tfTypes := terraform.GetTFTypes(*fileName)
	if tfTypes == "" {
		log.Fatal().Msgf("tf types not found. The filename must be like `my_tf_name_datasource.go` or `my_tf_name_resource.go`")
	}
	log.Info().Msgf("tf type : %s", tfTypes)

	packageName := terraform.GetPackageName(f)
	if packageName == "" {
		log.Fatal().Msg("package name not found")
	}
	log.Info().Msgf("package name : %s", packageName)

	categoryName, resourceName := terraform.GetTFName(f)
	if categoryName == "" {
		log.Fatal().Msg("tfname not found. Please add a comment like `// tfname: category_resource_name")
	}
	log.Info().Msgf("categoryName : %s", categoryName)
	log.Info().Msgf("resourceName : %s", resourceName)

	templateDocDir := filepath.Join(absPath, "../../../../templates/")
	dir, err = os.Stat(templateDocDir)
	if err != nil || !dir.IsDir() {
		log.Info().Msgf("creating directory %s", templateDocDir)
		// create the directory
		if err := os.MkdirAll(templateDocDir, 0o755); err != nil {
			log.Fatal().Err(err).Msgf("error creating directory %s", templateDocDir)
		}
	}

	examplesDir := filepath.Join(absPath, "../../../../examples/")
	dir, err = os.Stat(examplesDir)
	if err != nil || !dir.IsDir() {
		log.Info().Msgf("creating directory %s", templateDocDir)
		// create the directory
		if err := os.MkdirAll(templateDocDir, 0o755); err != nil {
			log.Fatal().Err(err).Msgf("error creating directory %s", templateDocDir)
		}
	}

	ressOrData := "resources" //nolint:goconst
	if strings.Contains(*fileName, "datasource") {
		ressOrData = "data-sources" //nolint:goconst
	}

	examplesDirSubDir := filepath.Join(examplesDir, ressOrData, "cloudavenue_"+packageName+"_"+resourceName)
	subDir, err := os.Stat(examplesDirSubDir)
	if err != nil || !subDir.IsDir() {
		log.Info().Msgf("creating directory %s", examplesDirSubDir)
		// create the directory
		if err := os.MkdirAll(examplesDirSubDir, 0o755); err != nil {
			log.Fatal().Err(err).Msgf("error creating directory %s", examplesDirSubDir)
		}
	}

	t := genTemplateConf(categoryName, resourceName, packageName, testDir, *fileName, schemaDir, templateDocDir, examplesDir)
	if err := t.createTemplateFiles(tfTypes); err != nil {
		log.Fatal().Err(err).Msg("error creating file")
	}

	log.Info().Msg("Run linter")

	if err := exec.Command("golangci-lint", "run", "--fix", "--config", ".golangci.yml").Run(); err != nil {
		log.Error().Err(err).Msg("error running linter")
	}

	log.Info().Msg("Done")
}
