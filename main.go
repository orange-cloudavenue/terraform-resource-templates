package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/orange-cloudavenue/terraform-resource-templates/internal/terraform"
	"github.com/orange-cloudavenue/terraform-resource-templates/pkg/file"

	_ "embed"
)

func main() {
	fileName := flag.String("filename", "", "filename")
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *fileName == "" {
		log.Fatal().Msg("filename is required")
	}

	if !file.IsFileExists(*fileName) {
		log.Fatal().Msgf("file %s not found", *fileName)
	}

	// if .golangci.yml exists, use it
	if !file.IsFileExists(".golangci.yml") {
		// read file ../.golangci.yml
		golangCIFile, err := file.ReadFile(".golangci.yml")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read file")
		}

		// parse file ../.golangci.yml
		golangCI := &golangCI{}
		if err := yaml.Unmarshal(golangCIFile, golangCI); err != nil {
			log.Fatal().Err(err).Msg("Failed to parse file")
		}

		varNaming := make([]string, 0)

		// get all var-naming rules
		for _, rule := range golangCI.LintersSettings.Revive.Rules {
			if rule.Name == "var-naming" {
				for _, arg := range rule.Arguments {
					varNaming = append(varNaming, arg[0].(string))
				}
			}
		}

		// configure var-naming rules
		for _, v := range varNaming {
			strcase.ConfigureAcronym(strings.ToUpper(v), strings.ToLower(v))
		}
	}
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

	t := genTemplateConf(categoryName, resourceName, packageName, testDir, *fileName, schemaDir)
	if err := t.createTemplateFiles(tfTypes); err != nil {
		log.Fatal().Err(err).Msg("error creating file")
	}

	log.Info().Msg("Run linter")

	if err := exec.Command("golangci-lint", "run", "--fix", "--config", ".golangci.yml").Run(); err != nil {
		log.Error().Err(err).Msg("error running linter")
	}

	log.Info().Msg("Done")
}

type golangCI struct {
	LintersSettings struct {
		Revive struct {
			Revive                interface{} `yaml:"revive"`
			IgnoreGeneratedHeader bool        `yaml:"ignore-generated-header"` //nolint:tagliatelle
			Severity              string      `yaml:"severity"`
			Rules                 []struct {
				Name      string          `yaml:"name"`
				Severity  string          `yaml:"severity"`
				Disabled  bool            `yaml:"disabled"`
				Arguments [][]interface{} `yaml:"arguments,omitempty"`
			} `yaml:"rules"`
		} `yaml:"revive"`
	} `yaml:"linters-settings"` //nolint:tagliatelle
}
