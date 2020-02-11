package service_test

import (
	"context"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestStartRecordingRequests(t *testing.T) {
	//t.SkipNow()
	startRecordingRequests("x_mock_status_code", "localhost:8080")
}

func startRecordingRequests(name string, address string) {

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Infof("Record request to %s", "testdata/html/http/"+name+".input.http")
			reqf, err := os.Create("testdata/html/http/" + name + ".input.http")
			if err == nil {
				//record request
				requestBytes, _ := httputil.DumpRequest(req, true)
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
		}),
	}

	go func() {
		log.Infof("Recording on %s", address)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("error setting up server: " + err.Error())
		}
	}() //what does this mean? It means to invoke the function

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	srv.Shutdown(ctx)

	os.Exit(0)
}
