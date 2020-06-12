package cmd

import "context"

// Execute runs the cli interface
func Execute() {
	rootCmd.ExecuteContext(context.Background())
}
