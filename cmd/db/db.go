/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package db

import (
	"main/coldbrew/cmd/db/create"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "db",
	Short: "database actions",
	Long:  `database actions`,
}

func init() {
	Cmd.AddCommand(create.Cmd)
}
