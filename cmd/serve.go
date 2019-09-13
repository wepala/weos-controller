// Copyright © 2019 Akeem Philbert <akeem.philbert@wepala.com>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port string
var apiConfigFlag string
var controllerConfigFlag string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve api ",
	Long:  `This command is used in conjunction with a sub command to serve apis of the service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called for real")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", viper.GetString("PORT"), "Port to open TCP connection. Defaults to 5000")
}
