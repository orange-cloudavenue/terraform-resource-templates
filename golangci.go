package main

import (
	"strings"

	"github.com/orange-cloudavenue/terraform-resource-templates/pkg/file"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

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

func extractVarNaming() {
	// if .golangci.yml exists, use it
	if file.IsFileExists(".golangci.yml") {
		log.Info().Msg("using .golangci.yml file")
		// read file ../.golangci.yml
		golangCIFile, err := file.ReadFile(".golangci.yml")
		if err != nil {
			return
		}

		// parse file ../.golangci.yml
		golangCI := &golangCI{}
		if err := yaml.Unmarshal(golangCIFile, golangCI); err != nil {
			log.Info().Msgf("error parsing .golangci.yml file: %s. Ignore it", err)
		}

		// get all var-naming rules
		for _, rule := range golangCI.LintersSettings.Revive.Rules {
			if rule.Name == "var-naming" {
				for _, arg := range rule.Arguments {
					for _, v := range arg {
						// check if arg is a string
						argS, ok := v.(string)
						if !ok {
							continue
						}
						// configure var-naming rules
						upperAcronyms.Store(strings.ToUpper(argS), strings.ToLower(argS))
					}
				}
			}
		}
	}
}
