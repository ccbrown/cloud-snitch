package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/api"
	"github.com/ccbrown/cloud-snitch/backend/app"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "serves the api",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		go catchSignal(cancel)

		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		api := api.New(a, rootConfig.API)
		port, _ := cmd.Flags().GetInt("port")

		server := &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: api,
		}

		go func() {
			<-ctx.Done()
			if err := server.Shutdown(context.Background()); err != nil {
				zap.L().Error(err.Error())
			}
		}()

		zap.L().Info("listening", zap.String("url", fmt.Sprintf("http://127.0.0.1:%v", port)))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}

		return nil
	},
}

func init() {
	apiCmd.Flags().IntP("port", "p", 8080, "the port to listen on")

	rootCmd.AddCommand(apiCmd)
}
