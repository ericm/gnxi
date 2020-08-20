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
	"net/http"

	"github.com/google/gnxi/gnxi_tester/config"
)

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

}

func handleTargetSet(w http.ResponseWriter, r *http.Request) {

}
