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

var httpCmd = &cobra.Command{
	Use:   "html",
	Short: "Start html server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("HTML Server started")
		//create controller service
		controllerService, _ := service.NewControllerService("api.yaml", "")
		//setup html handler
		htmlHandler := service.NewHTTPServer(controllerService, "static")
		srv := &http.Server{
			Addr:         "0.0.0.0:" + port,
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 30,
			IdleTimeout:  time.Second * 60,
			Handler:      htmlHandler,
		}

		go func() {
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
	serveCmd.AddCommand(httpCmd)
}
