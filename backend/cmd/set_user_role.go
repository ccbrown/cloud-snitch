package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

var setUserRoleCmd = &cobra.Command{
	Use:   "set-user-role",
	Short: "sets a user's role",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())

		go catchSignal(cancel)

		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		id, _ := cmd.Flags().GetString("id")

		roleString, _ := cmd.Flags().GetString("role")
		role := model.UserRole(roleString)
		if !role.IsValid() {
			return fmt.Errorf("invalid role")
		}

		patched, err := a.SetUserRole(ctx, model.Id(id), role)
		if err != nil {
			return err
		} else if patched == nil {
			return fmt.Errorf("no such user")
		}

		return nil
	},
}

func init() {
	setUserRoleCmd.Flags().String("id", "", "the user id")
	setUserRoleCmd.MarkFlagRequired("id")

	setUserRoleCmd.Flags().String("role", "", "the user's new role (example: \"administrator\")")
	setUserRoleCmd.MarkFlagRequired("role")

	rootCmd.AddCommand(setUserRoleCmd)
}
