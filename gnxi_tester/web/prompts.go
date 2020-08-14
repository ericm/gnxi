/* Copyright 2020 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/gnxi/gnxi_tester/config"
)

// handleConfigSet will set the prompt config variables in persistent storage.
func handlePromptsSet(w http.ResponseWriter, r *http.Request) {
	prompts := &config.Prompts{}
	if err := json.NewDecoder(r.Body).Decode(prompts); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	if err := prompts.Set(); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

// handleConfigGet will get all prompt config variables from persistent storage.
func handlePromptsGet(w http.ResponseWriter, r *http.Request) {
	prompts := config.GetPrompts()
	if err := json.NewEncoder(w).Encode(prompts); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// handlePromptsList will give all the required fields needed to be given by the client.
func handlePromptsList(w http.ResponseWriter, r *http.Request) {
	// TODO: reflect fields of prompts.
}
