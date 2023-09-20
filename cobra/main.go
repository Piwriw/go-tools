package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// base 命令
var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "a brief description of your applcation",
	Long:  "a long desc",
}

// cmd help
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "help for short desc",
	Long:  "help for long desc",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("do help")
		return nil
	},
}

// cmd export
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export for short desc",
	Long:  "export for long desc",
}

func main() {
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(exportCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
