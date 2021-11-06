package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "snocode",
}

func init() {
	rootCmd.AddCommand(genCmd, minifyCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
