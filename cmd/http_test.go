package cmd_test

import (
	"bitbucket.org/wepala/weos-controller/cmd"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"runtime"
	"syscall"
	"testing"
)

//TODO test that a http server is running

func runCommand(ready chan bool, command *cobra.Command, args []string) {
	command.Run(command, args)
	ready <- true
}

func TestNewHTTPCmd(t *testing.T) {
	t.SkipNow()
	var url string

	if runtime.GOOS == "darwin" {
		url = "localhost:8080"
	} else {
		url = "localhost:80"
	}

	command, _ := cmd.NewHTTPCmd()
	defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	done := make(chan bool, 1)
	go runCommand(done, command, []string{url})

	newreq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		r    *http.Request
	}{
		{name: "Home Page", r: newreq("GET", "http://"+url+"/", nil)},
		{name: "About Page", r: newreq("GET", "http://"+url+"/about", nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.r)
			if resp != nil {
				if resp.StatusCode > 399 {
					t.Errorf("received a status code %d expected status code 2XX", resp.StatusCode)
				}
				defer resp.Body.Close()
			}
			if err != nil {
				t.Errorf("page '%s' had an error '%s'", tt.name, err.Error())
			}

		})
	}

}
