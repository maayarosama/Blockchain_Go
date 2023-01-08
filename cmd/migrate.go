/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Blockchain_Go/database"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk(getDataDirFromCmd(cmd))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			block0 := database.NewBlock(
				database.Hash{},
				0,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "andrej", 3, ""),
					database.NewTx("andrej", "andrej", 700, "reward"),
				},
			)

			block0hash, err := state.AddBlock(block0)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block1 := database.NewBlock(
				block0hash,
				1,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "babayaga", 2000, ""),
					database.NewTx("andrej", "andrej", 100, "reward"),
					database.NewTx("babayaga", "andrej", 1, ""),
					database.NewTx("babayaga", "caesar", 1000, ""),
					database.NewTx("babayaga", "andrej", 50, ""),
					database.NewTx("andrej", "andrej", 600, "reward"),
				},
			)

			block1hash, err := state.AddBlock(block1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block2 := database.NewBlock(
				block1hash,
				2,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "andrej", 24700, "reward"),
				},
			)

			_, err = state.AddBlock(block2)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
	addDefaultRequiredFlags(migrateCmd)

	return migrateCmd
}

func init() {
	rootCmd.AddCommand(migrateCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
