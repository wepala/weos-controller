package cmd

import (
	"bitbucket.org/wepala/weos-controller/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var httpMockCmd = &cobra.Command{
	Use:   "http-mock",
	Short: "Start mock html server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//create controller service
		controllerService, _ := service.NewControllerService(apiYaml, configYaml, nil)
		//setup html handler
		htmlHandler := service.NewHTTPServer(controllerService, "static")
		srv := &http.Server{
			Addr:         args[0],
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 30,
			IdleTimeout:  time.Second * 60,
			Handler:      htmlHandler,
		}

		go func() {
			log.Infof("Mock HTML Server started on %s", args[0])
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
	},
}

func init() {
	//rootCmd.AddCommand(httpMockCmd)
	serveCmd.AddCommand(httpMockCmd)
}