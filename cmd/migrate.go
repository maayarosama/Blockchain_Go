// /*
// Copyright © 2022 NAME HERE <EMAIL ADDRESS>
// */
package cmd

// import (
// 	"Blockchain_Go/database"
// 	"Blockchain_Go/node"
// 	"Blockchain_Go/wallet"
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/spf13/cobra"
// )

// var migrateCmd = func() *cobra.Command {
// 	var migrateCmd = &cobra.Command{
// 		Use:   "migrate",
// 		Short: "Migrates the blockchain database according to new business rules.",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			miner, _ := cmd.Flags().GetString(flagMiner)
// 			ip, _ := cmd.Flags().GetString(flagIP)
// 			port, _ := cmd.Flags().GetUint64(flagPort)

// 			andrej := database.NewAccount(wallet.AndrejAccount)
// 			babayaga := database.NewAccount(wallet.BabaYagaAccount)
// 			caesar := database.NewAccount(wallet.CaesarAccount)

// 			peer := node.NewPeerNode(
// 				"127.0.0.1",
// 				8080,
// 				true,
// 				andrej,
// 				false,
// 			)

// 			n := node.New(getDataDirFromCmd(cmd), ip, port, database.NewAccount(miner), peer)

// 			// n.AddPendingTX(database.NewTx(andrej, andrej, 3, ""), peer)
// 			// n.AddPendingTX(database.NewTx(andrej, babayaga, 2000, ""), peer)
// 			// n.AddPendingTX(database.NewTx(babayaga, andrej, 1, ""), peer)
// 			// n.AddPendingTX(database.NewTx(babayaga, caesar, 1000, ""), peer)
// 			// n.AddPendingTX(database.NewTx(babayaga, andrej, 50, ""), peer)

// 			ctx, closeNode := context.WithTimeout(context.Background(), time.Minute*15)

// 			go func() {
// 				ticker := time.NewTicker(time.Second * 10)

// 				for {
// 					select {
// 					case <-ticker.C:
// 						if !n.LatestBlockHash().IsEmpty() {
// 							closeNode()
// 							return
// 						}
// 					}
// 				}
// 			}()

// 			err := n.Run(ctx)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		},
// 	}

// 	addDefaultRequiredFlags(migrateCmd)
// 	migrateCmd.Flags().String(flagMiner, node.DefaultMiner, "miner account of this node to receive block rewards")
// 	migrateCmd.Flags().String(flagIP, node.DefaultIP, "exposed IP for communication with peers")
// 	migrateCmd.Flags().Uint64(flagPort, node.DefaultHTTPort, "exposed HTTP port for communication with peers")

// 	return migrateCmd
// }

// func init() {
// 	rootCmd.AddCommand(migrateCmd())

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
