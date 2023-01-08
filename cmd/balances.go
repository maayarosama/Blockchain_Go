/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"Blockchain_Go/database"

	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd())

	return balancesCmd
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk(getDataDirFromCmd(cmd))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Accounts balances at %x:\n", state.LatestBlockHash())
			fmt.Println("__________________")
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account.String(), balance))
			}
			fmt.Println("")
			fmt.Printf("Accounts nonces:")
			fmt.Println("")
			fmt.Println("__________________")
			fmt.Println("")
			for account, nonce := range state.Account2Nonce {
				fmt.Println(fmt.Sprintf("%s: %d", account.String(), nonce))
			}
		},
	}

	addDefaultRequiredFlags(balancesListCmd)
	return balancesListCmd
}

func init() {
	rootCmd.AddCommand(balancesCmd())

	rootCmd.AddCommand(balancesListCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// balancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// balancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
