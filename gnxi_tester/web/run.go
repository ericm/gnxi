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
	"github.com/spf13/viper"
)

type runRequest struct {
	Prompts string   `json:"prompts"`
	Device  string   `json:"device"`
	Tests   []string `json:"tests"`
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	request := runRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	configPrompts := viper.GetStringMap("web.prompts")
	promptsGeneric, ok := configPrompts[request.Prompts]
	if !ok {
		logErr(r.Header, fmt.Errorf("%s prompts not found", request.Prompts))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	prompts, ok := promptsGeneric.(config.Prompts)
	if !ok {
		logErr(r.Header, fmt.Errorf("%s prompts invalid", request.Prompts))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	devices := config.GetDevices()
	device, ok := devices[request.Device]
	if !ok {
		logErr(r.Header, fmt.Errorf("%s device not found", request.Device))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
