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
	"net/http"

	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
)

// InitRouter will setup required routes for the web API and listen on the address.
func InitRouter(laddr string) {
	r := generateRouter()
	log.Infof("Running web API on %s", laddr)
	if err := http.ListenAndServe(laddr, r); err != nil {
		log.Exit(err)
	}
}

func generateRouter() *mux.Router {
	r := mux.NewRouter()
	f := FileHandler{FileSys: afero.NewOsFs()}
	r.HandleFunc("/prompts", handlePromptsGet).Methods("GET")
	r.HandleFunc("/prompts/list", handlePromptsList).Methods("GET")
	r.HandleFunc("/prompts", handlePromptsSet).Methods("POST", "PUT", "OPTIONS")
	r.HandleFunc("/prompts/{name}", handlePromptsDelete).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/target", handleTargetsGet).Methods("GET")
	r.HandleFunc("/target/{name}", handleTargetGet).Methods("GET")
	r.HandleFunc("/target/{name}", handleTargetSet).Methods("POST", "PUT", "OPTIONS")
	r.HandleFunc("/target/{name}", handleTargetDelete).Methods("DELETE")
	r.HandleFunc("/file", f.handleFileUpload).Methods("POST")
	r.HandleFunc("/file/{file}", f.handleFileDelete).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/run", handleRun).Methods("POST", "OPTIONS")
	r.HandleFunc("/run/output", handleRunOutput).Methods("GET")
	r.HandleFunc("/test", handleTestsGet).Methods("GET")
	r.HandleFunc("/test/order", handleTestsOrderGet).Methods("GET")
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})
	return r
}

var logErr = func(head http.Header, err error) {
	log.Errorf("error occured in view %v: %v", head, err)
}
