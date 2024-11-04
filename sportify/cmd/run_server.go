package cmd

import (
	"context"

	"github.com/TheVovchenskiy/sportify-backend/server"

	"github.com/spf13/cobra"
)

var runServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Runs sportify backend server.",
	Long:  "Use this command to run sportify http server.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		configPaths, err := cmd.Flags().GetStringSlice("config-path")
		if err != nil {
			return err
		}

		srv := server.Server{}
		baseCtx := context.Background()

		if err := srv.Run(baseCtx, configPaths); err != nil {
			panic(err)
		}
		return nil
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(runServerCmd)

	//nolint:lll
	runServerCmd.Flags().StringSliceP("config-path", "c", []string{}, "Path to config file dir to search in for config. Can be accepted multiple times.")
}
