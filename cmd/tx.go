/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// import (
// 	"Blockchain_Go/database"
// 	"fmt"
// 	"os"

// 	"github.com/spf13/cobra"
// )

// const flagFrom = "from"
// const flagTo = "to"
// const flagValue = "value"
// const flagData = "data"

// func txCmd() *cobra.Command {
// 	var txsCmd = &cobra.Command{
// 		Use:   "tx",
// 		Short: "Interact with txs (add...).",
// 		PreRunE: func(cmd *cobra.Command, args []string) error {
// 			return incorrectUsageErr()
// 		},
// 		Run: func(cmd *cobra.Command, args []string) {
// 		},
// 	}

// 	txsCmd.AddCommand(txAddCmd())

// 	return txsCmd
// }

// func txAddCmd() *cobra.Command {
// 	var cmd = &cobra.Command{
// 		Use:   "add",
// 		Short: "Adds new TX to database.",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			from, _ := cmd.Flags().GetString(flagFrom)
// 			to, _ := cmd.Flags().GetString(flagTo)
// 			value, _ := cmd.Flags().GetUint(flagValue)
// 			data, _ := cmd.Flags().GetString(flagData)

// 			tx := database.NewTx(database.NewAccount(from), database.NewAccount(to), value, data)

// 			dataDir, _ := cmd.Flags().GetString(flagDataDir)
// 			state, err := database.NewStateFromDisk(dataDir)
// 			if err != nil {
// 				fmt.Fprintln(os.Stderr, err)
// 				os.Exit(1)
// 			}
// 			defer state.Close()

// 			err = state.AddTx(tx)
// 			if err != nil {
// 				fmt.Fprintln(os.Stderr, err)
// 				os.Exit(1)
// 			}

// 			_, err = state.Persist()
// 			if err != nil {
// 				fmt.Fprintln(os.Stderr, err)
// 				os.Exit(1)
// 			}

// 			fmt.Println("TX successfully persisted to the ledger.")
// 		},
// 	}

// 	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
// 	cmd.MarkFlagRequired(flagFrom)

// 	cmd.Flags().String(flagTo, "", "To what account to send tokens")
// 	cmd.MarkFlagRequired(flagTo)

// 	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
// 	cmd.MarkFlagRequired(flagValue)

// 	cmd.Flags().String(flagData, "", "Possible values: 'reward'")

// 	return cmd
// }

// func init() {
// 	rootCmd.AddCommand(txCmd())

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
