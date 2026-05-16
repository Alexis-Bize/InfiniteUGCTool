// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package debug exposes a single global toggle for verbose diagnostic
// dumps (auth-flow HTML, regex context, response headers, ...). Enabled
// via the `--debug` CLI flag or the IUGC_DEBUG environment variable; off
// by default so normal sessions stay silent on success and terse on
// expected errors.
package debug

import "os"

var enabled bool

func init() {
	if os.Getenv("IUGC_DEBUG") != "" {
		enabled = true
	}
}

// Enable turns verbose diagnostic dumps on.
func Enable() { enabled = true }

// Enabled reports whether verbose diagnostic dumps are turned on.
func Enabled() bool { return enabled }
