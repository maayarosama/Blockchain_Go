/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"Blockchain_Go/cmd"

	"fmt"
)

func main() {
	// state, err := database.NewStateFromDisk()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
	// defer state.Close()

	// block0 := database.NewBlock(
	// 	database.Hash{},
	// 	uint64(time.Now().Unix()),
	// 	[]database.Tx{
	// 		database.NewTx("andrej", "andrej", 3, ""),
	// 		database.NewTx("andrej", "andrej", 700, "reward"),
	// 	},
	// )

	// state.AddBlock(block0)
	// block0hash, _ := state.Persist()

	// block1 := database.NewBlock(
	// 	block0hash,
	// 	uint64(time.Now().Unix()),
	// 	[]database.Tx{
	// 		database.NewTx("andrej", "babayaga", 2000, ""),
	// 		database.NewTx("andrej", "andrej", 100, "reward"),
	// 		database.NewTx("babayaga", "andrej", 1, ""),
	// 		database.NewTx("babayaga", "caesar", 1000, ""),
	// 		database.NewTx("babayaga", "andrej", 50, ""),
	// 		database.NewTx("andrej", "andrej", 600, "reward"),
	// 	},
	// )

	// state.AddBlock(block1)
	// state.Persist()

	cmd.Execute()
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
