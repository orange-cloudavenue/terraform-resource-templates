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
	"context"
	"fmt"
	"os"
	"runtime"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/rs/zerolog/log"
)

func getLatest() (*selfupdate.Release, bool) {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("orange-cloudavenue/terraform-resource-templates"))
	if err != nil {
		log.Error().Err(err).Msg("error occurred while detecting version")
		return nil, false
	}

	if !found {
		log.Error().Msgf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
		return nil, false
	}

	return latest, true
}

func checkVersion() {
	latest, found := getLatest()
	if !found {
		return
	}

	if !latest.LessOrEqual(version) {
		log.Warn().Msgf("Running outdated version %s, latest version is %s", version, latest.Version())
		log.Warn().Msgf("Run `terraform-resource-templates -update` to update the application")
		return
	}
}

func selfUpdate() error {
	latest, found := getLatest()
	if !found {
		return fmt.Errorf("could not find latest version")
	}

	exe, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("could not locate executable path")
		return err
	}

	log.Info().Msgf("Updating to version %s", latest.Version())
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		log.Error().Err(err).Msg("error occurred while updating binary")
		return err
	}
	log.Info().Msgf("Successfully updated to version %s", latest.Version())
	return nil
}
