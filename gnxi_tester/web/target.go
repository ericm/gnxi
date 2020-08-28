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
	"path"

	"errors"
	"net/http"

	"github.com/google/gnxi/gnxi_tester/config"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func getNameParam(w http.ResponseWriter, r *http.Request) string {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		logErr(r.Header, errors.New("name param not set"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return ""
	}
	return name
}

func handleTargetsGet(w http.ResponseWriter, r *http.Request) {
	targets := config.GetDevices()
	response, err := json.Marshal(targets)
	if err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Write(response)
}

func handleTargetGet(w http.ResponseWriter, r *http.Request) {
	name := getNameParam(w, r)
	if name == "" {
		return
	}
	devices := config.GetDevices()
	device, ok := devices[name]
	if !ok {
		logErr(r.Header, errors.New("device not found"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(device); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleTargetSet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	name := getNameParam(w, r)
	if name == "" {
		return
	}
	devices := config.GetDevices()
	if devices == nil {
		devices = map[string]config.Device{}
	}
	device := config.Device{}
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if device.Ca == "" || device.CaKey == "" {
		logErr(r.Header, errors.New("ca or ca key not set"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}
	if !path.IsAbs(device.Ca) {
		device.Ca = path.Join(filesDir(), device.Ca)
	}
	if !path.IsAbs(device.CaKey) {
		device.CaKey = path.Join(filesDir(), device.CaKey)
	}
	if err := config.SetTarget(name, device.Address, device.Ca, device.CaKey, false); err != nil {
		logErr(r.Header, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// handleTargetDelete deletes a target.
func handleTargetDelete(w http.ResponseWriter, r *http.Request) {
	name := getNameParam(w, r)
	targets := config.GetDevices()
	delete(targets, name)
	viper.Set("targets.devices", targets)
}
