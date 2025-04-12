package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var rootConfig Config

var rootCmd = &cobra.Command{
	Use:           filepath.Base(os.Args[0]),
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logConfig := zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		if term.IsTerminal(unix.Stdout) {
			logConfig.Encoding = "console"
			logConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		}

		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}

		logger, err := logConfig.Build()
		if err != nil {
			return fmt.Errorf("error initializing logger: %w", err)
		}
		zap.ReplaceGlobals(logger)

		mustFindConfigFile := false
		if config, _ := cmd.Flags().GetString("config"); config != "" {
			viper.SetConfigFile(config)
			mustFindConfigFile = true
		} else {
			viper.AddConfigPath(".")
			viper.SetConfigName("config")
		}

		if err := LoadConfigEnvVariables(); err != nil {
			return fmt.Errorf("error initializing config environment variables: %w", err)
		} else if err := viper.ReadInConfig(); err != nil && (mustFindConfigFile || !errors.As(err, &viper.ConfigFileNotFoundError{})) {
			return fmt.Errorf("error reading config: %w", err)
		} else if err := viper.Unmarshal(&rootConfig); err != nil {
			return fmt.Errorf("error unmarshaling config: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "make output more verbose")
}

func Execute() {
	defer zap.L().Sync()
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal(err.Error())
	}
}

func catchSignal(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	zap.L().Info("signal caught. shutting down...")
	cancel()
}
