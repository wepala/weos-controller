package handlers

import (
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type WeOSCommons struct {
	logger     log.Ext1FieldLogger
	HTTPClient *http.Client
	Store      sessions.Store
	Config     *Config
}

type RecordingConfig struct {
	BaseFolder string `json:"baseFolder"`
}

type Config struct {
	Recording *RecordingConfig `json:"recording"`
}

func (w *WeOSCommons) RecordRequest(rw http.ResponseWriter, r *http.Request) {
	count := 0
	name := strings.Replace(r.URL.Path, "/", "_", -1)
	baseFolder := w.Config.Recording.BaseFolder
	if baseFolder == "" {
		baseFolder = "testdata/http"
	}
	if count < 1 {
		log.Infof("Record request to %s", baseFolder+"/"+name+".input.http")
		count += 1
		reqf, err := os.Create(baseFolder + "/" + name + ".input.http")
		if err == nil {
			//record request
			requestBytes, _ := httputil.DumpRequest(r, true)
			_, err := reqf.Write(requestBytes)
			if err != nil {
				log.Errorf("error occurred during recording: %s", err)
			}
		} else {
			log.Errorf("error recording request because of error: %s", err)
		}

		defer func() {
			reqf.Close()
			if r := recover(); r != nil {
				log.Errorf("Recording failed with errors: %s", r)
			}
		}()
	}
}
