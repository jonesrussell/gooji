package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gooji-admin",
		Short: "Gooji administration tool",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gooji-admin v2.0.0")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
