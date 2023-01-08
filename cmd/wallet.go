/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Blockchain_Go/wallet"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Manages accounts, keys, cryptography.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	walletCmd.AddCommand(walletNewAccountCmd())

	return walletCmd
}

func walletNewAccountCmd() *cobra.Command {
	var walletAccountCmd = &cobra.Command{
		Use:   "new-account",
		Short: "Creates a new account with a new set of a elliptic-curve Private + Public keys.",
		Run: func(cmd *cobra.Command, args []string) {
			password := getPassPhrase("Please enter a password to encrypt the new wallet:", true)

			dataDir := getDataDirFromCmd(cmd)
			ks := keystore.NewKeyStore(
				wallet.GetKeystoreDirPath(dataDir),
				keystore.StandardScryptN,
				keystore.StandardScryptP,
			)
			acc, err := ks.NewAccount(password)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("New account created: %s\n", acc.Address.Hex())
		},
	}

	addDefaultRequiredFlags(walletAccountCmd)

	return walletAccountCmd
}

func getPassPhrase(prompt string, confirmation bool) string {
	fmt.Println(prompt)
	password, err := console.Stdin.PromptPassword("Password: ")
	if err != nil {
		utils.Fatalf("Failed to read password: %v", err)
	}

	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat password: ")
		if err != nil {
			utils.Fatalf("Failed to read password confirmation: %v", err)
		}
		if password != confirm {
			utils.Fatalf("Passwords do not match")
		}
	}

	return password
}

func init() {
	rootCmd.AddCommand(walletCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
