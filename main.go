package main

import (
	"github.com/little-cui/sc-compose/command"
	"github.com/little-cui/sc-compose/usage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sc-compose",
	Short: "A simple command line tool for ServiceCenter",
}

func init() {
	cobra.EnablePrefixMatching = true

	rootCmd.AddCommand(
		command.NewUpCommand(),
	)
}

func main() {
	Start()
}

func Start() {
	rootCmd.SetUsageFunc(usage.UsageFunc)

	// Make help just show the usage
	rootCmd.SetHelpTemplate(`{{.UsageString}}`)

	if err := rootCmd.Execute(); err != nil {
		command.ExitWithError(command.ErrExit, err.Error())
	}
}
