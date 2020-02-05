package cmd

func init() {
	command, _ := NewHTTPCmd()
	serveCmd.AddCommand(command)
}
