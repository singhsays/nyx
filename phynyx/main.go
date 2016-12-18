package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// generatekeyCmd
	generatekeyCmd := &cobra.Command{
		Use:   "generatekey",
		Short: "generates a new ecryption key.",
		Long:  `Generates a base64 encoded 32 bit random key, for use with AES 265 bit encryption.`,
	}
	generatekeyCmd.Run = func(cmd *cobra.Command, args []string) {
		// Generate key
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(base64.StdEncoding.EncodeToString([]byte(key)))
	}

	// rootCmd
	RootCmd := &cobra.Command{
		Use:   "phynyx",
		Short: "phynyx is the command-line tool for nyx.",
		Long:  `phynyx is the command-line tool for nyx.`,
	}
	RootCmd.AddCommand(generatekeyCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
