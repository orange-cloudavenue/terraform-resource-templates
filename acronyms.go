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
	"strings"
	"sync"
)

var upperAcronyms = sync.Map{}

func fixUpperCase(s string) string {
	upperAcronyms.Range(func(key, value interface{}) bool {
		// acronym registered has Key as upper case and value as lower case
		// s is the string to be fixed. The string has a CamelCase format
		// Example
		// s = "HttpStatusCode"
		// key = "HTTP"
		// value = "http"
		// we need to replace "Http" with "HTTP" in s

		if strings.Contains(s, strings.Title(value.(string))) {
			// Replace "Http" or "http" with "HTTP" in s
			s = strings.ReplaceAll(s, strings.Title(value.(string)), key.(string))
		}

		return true
	})
	return s
}
