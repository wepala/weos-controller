package cmd

import (
	"bitbucket.org/wepala/weos-controller/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var apiConfigFlag string
var controllerConfigFlag string

func NewHTTPCmd(apiConfig string, controllerConfig string) (*cobra.Command, *http.Server) {

	srv := &http.Server{
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	return &cobra.Command{
		Use:   "html",
		Short: "Start html server",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			controllerService, err := service.NewControllerService(apiConfig, controllerConfig, service.NewPluginLoader())
			if err != nil {
				log.Fatalf("error encountered setting up command '%s'", err.Error())
			}
			log.Info("HTML Server started")
			//setup html handler
			htmlHandler := service.NewHTTPServer(controllerService, "static")
			srv.Addr = args[0]
			srv.Handler = htmlHandler
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
	}, srv
}

func init() {
	serveCmd.Flags().StringVar(&apiConfigFlag, "apiConfig", "", "Source directory to read from")
	viper.BindPFlag("apiConfigFlag", serveCmd.LocalFlags().Lookup("apiConfig"))
	serveCmd.Flags().StringVar(&controllerConfigFlag, "controllerConfig", "", "Source directory to read from")
	viper.BindPFlag("controllerConfigFlag", serveCmd.LocalFlags().Lookup("controllerConfig"))
	command, _ := NewHTTPCmd(apiConfigFlag, controllerConfigFlag)
	serveCmd.AddCommand(command)
}
