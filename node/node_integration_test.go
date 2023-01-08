package node

import (
	"Blockchain_Go/database"
	"Blockchain_Go/fs"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getTestDataDirPath() string {
	return filepath.Join(os.TempDir(), ".tbb_test")
}

func TestNode_Run(t *testing.T) {
	// Remove the test directory if it already exists
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}
	// Construct a new Node instance
	n := New(
		datadir,
		"127.0.0.1",
		8085,
		database.NewAccount("andrej"),
		PeerNode{},
	)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx)
	if err.Error() != "http: Server closed" {
		// assert expected behaviour
		t.Fatal("node server was suppose to close after 5s")
	}
}
func TestNode_Mining(t *testing.T) {
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}
	nInfo := NewPeerNode("127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true)

	n := New(
		datadir,
		nInfo.IP,
		nInfo.Port,
		database.NewAccount("andrej"),
		nInfo,
	)
	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)

	go func() {
		time.Sleep(time.Second * miningIntervalSeconds / 3)
		tx := database.NewTx("andrej", "babayaga", 1, "")
		// Add it to the Mempool
		_ = n.AddPendingTX(tx, nInfo)
	}()
	go func() {
		time.Sleep(time.Second*miningIntervalSeconds + 2)
		tx := database.NewTx("andrej", "babayaga", 2, "")
		_ = n.AddPendingTX(tx, nInfo)
	}()
	go func() {
		// Periodically check if we mined the 2 blocks
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				// 2 blocks mined as expected? (height: 0 and 1)
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()
	_ = n.Run(ctx)
	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("2 pending TX not mined into 2 under 30m")
	}
}

func TestNode_MiningStopsOnNewSyncedBlock(t *testing.T) {
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}
	nInfo := NewPeerNode("127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)
	andrejAcc := database.NewAccount("andrej")
	babayagaAcc := database.NewAccount("babayaga")
	n := New(datadir, nInfo.IP, nInfo.Port, babayagaAcc, nInfo)
	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)
	tx1 := database.NewTx("andrej", "babayaga", 1, "")
	tx2 := database.NewTx("andrej", "babayaga", 2, "")
	tx2Hash, _ := tx2.Hash()

	validPreMinedPb := NewPendingBlock(
		database.Hash{},
		0,
		andrejAcc,
		[]database.Tx{tx1},
	)
	validSyncedBlock, err := Mine(ctx, validPreMinedPb)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds - 2))
		err := n.AddPendingTX(tx1, nInfo)
		if err != nil {
			t.Fatal(err)
		}
		err = n.AddPendingTX(tx2, nInfo)
		if err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining")
		}
		_, err := n.state.AddBlock(validSyncedBlock)
		if err != nil {
			t.Fatal(err)
		}
		// Mock the Andrej's block came from a network
		n.newSyncedBlocks <- validSyncedBlock
		time.Sleep(time.Second * 2)
		if n.isMining {
			t.Fatal("synced block should have canceled mining")
		}
		// Mined TX1 by Andrej should be removed from the Mempool
		_, onlyTX2IsPending := n.pendingTXs[tx2Hash.Hex()]
		if len(n.pendingTXs) != 1 && !onlyTX2IsPending {
			t.Fatal("TX1 should be still pending")
		}
		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should attempt to mine TX1 again")
		}
	}()

	go func() {
		// Regularly check whenever both TXs are now mined
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)
		// Take a snapshot of the DB balances
		// before the mining is finished and the 2 blocks
		// are created.
		startingAndrejBalance := n.state.Balances[andrejAcc]
		startingBabaYagaBalance := n.state.Balances[babayagaAcc]
		// Wait until the 30 mins timeout is reached or
		// the 2 blocks get mined and
		// the closeNode() is triggered
		<-ctx.Done()
		// Query balances again
		endAndrejBalance := n.state.Balances[andrejAcc]
		endBabaYagaBalance := n.state.Balances[babayagaAcc]
		// In TX1 Andrej transferred 1 TBB token to BabaYaga
		// In TX2 Andrej transferred 2 TBB tokens to BabaYaga
		expectedEndAndrejBalance := startingAndrejBalance - tx1.Value - tx2.Value + database.BlockReward
		expectedEndBabaYagaBalance :=
			startingBabaYagaBalance +
				tx1.Value +
				tx2.Value +
				database.BlockReward
		if endAndrejBalance != expectedEndAndrejBalance {
			t.Fatalf(
				"Andrej expected end balance is %d not %d",
				expectedEndAndrejBalance,
				endAndrejBalance,
			)
		}
		if endBabaYagaBalance != expectedEndBabaYagaBalance {
			t.Fatalf(
				"BabaYaga expected end balance is %d not %d",
				expectedEndBabaYagaBalance,
				endBabaYagaBalance,
			)
		}
		t.Logf("Before Andrej: %d TBB", startingAndrejBalance)
		t.Logf("Before BabaYaga: %d TBB", startingBabaYagaBalance)
		t.Logf("After Andrej: %d TBB", endAndrejBalance)
		t.Logf("After BabaYaga: %d TBB", endBabaYagaBalance)
	}()
	_ = n.Run(ctx)
	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("2 pending TX not mined into 2 blocks under 30m")
	}
	if len(n.pendingTXs) != 0 {
		t.Fatal("no pending TXs should be left to mine")
	}
}
