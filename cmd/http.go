package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wepala/weos-controller/service"
	"golang.org/x/net/context"
	_ "golang.org/x/sys/unix"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func NewHTTPCmd() (*cobra.Command, *http.Server) {
	srv := &http.Server{
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	return &cobra.Command{
		Use:   "http",
		Short: "Start html server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
			//set log format to json WECON-3
			log.SetFormatter(&log.JSONFormatter{})
			//create controller service
			controllerService, err := service.NewControllerService(apiYaml, service.NewPluginLoader())
			if err != nil {
				log.Fatalf("error occurred setting up controller service: %s", err)
			}
			//setup html handler
			htmlHandler := service.NewHTTPServer(controllerService, serveStatic, staticPath)
			srv := &http.Server{
				Addr:         args[0],
				WriteTimeout: time.Second * 30,
				ReadTimeout:  time.Second * 30,
				IdleTimeout:  time.Second * 60,
				Handler:      htmlHandler,
			}

			go func() {
				log.Infof("HTML Server started on %s", args[0])
				if err := srv.ListenAndServe(); err != nil {
					log.Fatal("error setting up server: " + err.Error())
				}
			}() //what does this mean? It means to invoke the function

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			err = srv.Shutdown(ctx)
			log.Infof("server shutdown: %s", err)

			os.Exit(0)
		}}, srv
}

func init() {
	command, _ := NewHTTPCmd()
	serveCmd.AddCommand(command)
}