package node

import (
	"Blockchain_Go/database"
	"Blockchain_Go/fs"
	"Blockchain_Go/wallet"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const testKsAndrejAccount = "0x73e4B4d8A69b8820744988A0f7D098cAa6066352"
const testKsBabaYagaAccount = "0x640946ff71f8497420ac7db36d363aa873e36028"
const testKsAndrejFile = "test_andrej--73e4B4d8A69b8820744988A0f7D098cAa6066352"
const testKsBabaYagaFile = "test_babayaga--640946ff71f8497420ac7db36d363aa873e36028"
const testKsAccountsPwd = "mayar"

func TestNode_Run(t *testing.T) {
	datadir, err := getTestDataDirPath()
	if err != nil {
		t.Fatal(err)
	}
	err = fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	n := New(datadir, "127.0.0.1", 8085, database.NewAccount(DefaultMiner), PeerNode{})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx)
	if err.Error() != "http: Server closed" {
		t.Fatal("node server was suppose to close after 5s")
	}
}

func TestNode_Mining(t *testing.T) {
	dataDir, andrej, babaYaga, err := setupTestNodeDir()
	if err != nil {
		t.Error(err)
	}
	defer fs.RemoveDir(dataDir)

	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		babaYaga,
		true,
	)

	n := New(dataDir, nInfo.IP, nInfo.Port, andrej, nInfo)

	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)

	go func() {
		time.Sleep(time.Second * miningIntervalSeconds / 3)

		tx := database.NewTx(andrej, babaYaga, 1, 1, "")
		signedTx, err := wallet.SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
		if err != nil {
			t.Error(err)
			return
		}

		_ = n.AddPendingTX(signedTx, nInfo)
	}()

	go func() {
		time.Sleep(time.Second*miningIntervalSeconds + 2)

		tx := database.NewTx(andrej, babaYaga, 2, 2, "")
		signedTx, err := wallet.SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
		if err != nil {
			t.Error(err)
			return
		}

		_ = n.AddPendingTX(signedTx, nInfo)
	}()

	go func() {
		// Periodically check if we mined the 2 blocks
		ticker := time.NewTicker(10 * time.Second)

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

	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("2 pending TX not mined into 2 under 30m")
	}
}

func TestNode_ForgedTx(t *testing.T) {
	dataDir, andrej, babaYaga, err := setupTestNodeDir()
	if err != nil {
		t.Error(err)
	}
	defer fs.RemoveDir(dataDir)

	n := New(dataDir, "127.0.0.1", 8085, andrej, PeerNode{})
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*15)
	andrejPeerNode := NewPeerNode("127.0.0.1", 8085, false, andrej, true)

	txValue := uint(5)
	txNonce := uint(1)
	tx := database.NewTx(andrej, babaYaga, txValue, txNonce, "")

	signedTx, err := wallet.SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		time.Sleep(time.Second * 1)

		_ = n.AddPendingTX(signedTx, andrejPeerNode)
	}()

	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds + 1))

		// Attempt to replay the same TX but with modified time
		// Because the TX.time changed, the TX.signature will be considered forged
		forgedTx := database.NewTx(andrej, babaYaga, txValue, txNonce, "")
		forgedSignedTx := database.NewSignedTx(forgedTx, signedTx.Sig)

		_ = n.AddPendingTX(forgedSignedTx, andrejPeerNode)
	}()

	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 0 {
		t.Fatal("was suppose to mine only one TX. The second TX was forged")
	}
}

func TestNode_ReplayedTx(t *testing.T) {
	dataDir, andrej, babaYaga, err := setupTestNodeDir()
	if err != nil {
		t.Error(err)
	}
	defer fs.RemoveDir(dataDir)

	n := New(dataDir, "127.0.0.1", 8085, andrej, PeerNode{})
	ctx, closeNode := context.WithCancel(context.Background())
	andrejPeerNode := NewPeerNode("127.0.0.1", 8085, false, andrej, true)
	babaYagaPeerNode := NewPeerNode("127.0.0.1", 8086, false, babaYaga, true)

	txValue := uint(5)
	txNonce := uint(1)
	tx := database.NewTx(andrej, babaYaga, txValue, txNonce, "")

	signedTx, err := wallet.SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
	if err != nil {
		t.Error(err)
		return
	}

	_ = n.AddPendingTX(signedTx, andrejPeerNode)

	go func() {
		ticker := time.NewTicker(time.Second * (miningIntervalSeconds - 3))
		wasReplayedTxAdded := false

		for {
			select {
			case <-ticker.C:
				// The Andrej's original TX got mined.
				// Execute the attack by replaying the TX again!
				if n.state.LatestBlock().Header.Number == 0 {
					if wasReplayedTxAdded && !n.isMining {
						closeNode()
						return
					}

					// Simulate the TX was submitted to different node
					n.archivedTXs = make(map[string]database.SignedTx)
					// Execute the attack
					_ = n.AddPendingTX(signedTx, babaYagaPeerNode)
					wasReplayedTxAdded = true
				}

				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	_ = n.Run(ctx)

	if n.state.Balances[babaYaga] == txValue*2 {
		t.Errorf("replayed attack was successful :( Damn digital signatures!")
		return
	}

	if n.state.LatestBlock().Header.Number == 1 {
		t.Errorf("the second block was not suppose to be persisted because it contained a malicious TX")
		return
	}
}

// The test logic summary:
//   - BabaYaga runs the node
//   - BabaYaga tries to mine 2 TXs
//   - The mining gets interrupted because a new block from Andrej gets synced
//   - Andrej will get the block reward for this synced block
//   - The synced block contains 1 of the TXs BabaYaga tried to mine
//   - BabaYaga tries to mine 1 TX left
//   - BabaYaga succeeds and gets her block reward
func TestNode_MiningStopsOnNewSyncedBlock(t *testing.T) {
	babaYaga := database.NewAccount(testKsBabaYagaAccount)
	andrej := database.NewAccount(testKsAndrejAccount)

	dataDir, err := getTestDataDirPath()
	if err != nil {
		t.Fatal(err)
	}

	genesisBalances := make(map[common.Address]uint)
	genesisBalances[andrej] = 1000000
	genesis := database.Genesis{Balances: genesisBalances}
	genesisJson, err := json.Marshal(genesis)
	if err != nil {
		t.Fatal(err)
	}

	err = database.InitDataDirIfNotExists(dataDir, genesisJson)
	defer fs.RemoveDir(dataDir)

	err = copyKeystoreFilesIntoTestDataDirPath(dataDir)
	if err != nil {
		t.Fatal(err)
	}

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)

	n := New(dataDir, nInfo.IP, nInfo.Port, babaYaga, nInfo)

	// Allow the test to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(context.Background(), time.Minute*30)

	tx1 := database.NewTx(andrej, babaYaga, 1, 1, "")
	tx2 := database.NewTx(andrej, babaYaga, 2, 2, "")

	signedTx1, err := wallet.SignTxWithKeystoreAccount(tx1, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
	if err != nil {
		t.Error(err)
		return
	}

	signedTx2, err := wallet.SignTxWithKeystoreAccount(tx2, andrej, testKsAccountsPwd, wallet.GetKeystoreDirPath(dataDir))
	if err != nil {
		t.Error(err)
		return
	}
	tx2Hash, err := signedTx2.Hash()
	if err != nil {
		t.Error(err)
		return
	}

	// Pre-mine a valid block without running the `n.Run()`
	// with Andrej as a miner who will receive the block reward,
	// to simulate the block came on the fly from another peer
	validPreMinedPb := NewPendingBlock(database.Hash{}, 0, andrej, []database.SignedTx{signedTx1})
	validSyncedBlock, err := Mine(ctx, validPreMinedPb)
	if err != nil {
		t.Fatal(err)
	}

	// Add 2 new TXs into the BabaYaga's node, triggers mining
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds - 2))

		err := n.AddPendingTX(signedTx1, nInfo)
		if err != nil {
			t.Fatal(err)
		}

		err = n.AddPendingTX(signedTx2, nInfo)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Interrupt the previously started mining with a new synced block
	// BUT this block contains only 1 TX the previous mining activity tried to mine
	// which means the mining will start again for the one pending TX that is left and wasn't in
	// the synced block
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
			t.Fatal("synced block should have canceled mining of already mined TX")
		}

		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining again the 1 TX not included in synced block")
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
		startingAndrejBalance := n.state.Balances[andrej]
		startingBabaYagaBalance := n.state.Balances[babaYaga]

		// Wait until the 30 mins timeout is reached or
		// the 2 blocks got already mined and the closeNode() was triggered
		<-ctx.Done()

		endAndrejBalance := n.state.Balances[andrej]
		endBabaYagaBalance := n.state.Balances[babaYaga]

		// In TX1 Andrej transferred 1 TBB token to BabaYaga
		// In TX2 Andrej transferred 2 TBB tokens to BabaYaga
		expectedEndAndrejBalance := startingAndrejBalance - tx1.Value - tx2.Value + database.BlockReward
		expectedEndBabaYagaBalance := startingBabaYagaBalance + tx1.Value + tx2.Value + database.BlockReward

		if endAndrejBalance != expectedEndAndrejBalance {
			t.Fatalf("Andrej expected end balance is %d not %d", expectedEndAndrejBalance, endAndrejBalance)
		}

		if endBabaYagaBalance != expectedEndBabaYagaBalance {
			t.Fatalf("BabaYaga expected end balance is %d not %d", expectedEndBabaYagaBalance, endBabaYagaBalance)
		}

		t.Logf("Starting Andrej balance: %d", startingAndrejBalance)
		t.Logf("Starting BabaYaga balance: %d", startingBabaYagaBalance)
		t.Logf("Ending Andrej balance: %d", endAndrejBalance)
		t.Logf("Ending BabaYaga balance: %d", endBabaYagaBalance)
	}()

	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("was suppose to mine 2 pending TX into 2 valid blocks under 30m")
	}

	if len(n.pendingTXs) != 0 {
		t.Fatal("no pending TXs should be left to mine")
	}
}

// Creates dir like: "/tmp/tbb_test945924586"
func getTestDataDirPath() (string, error) {
	return ioutil.TempDir(os.TempDir(), "tbb_test")
}

// Copy the pre-generated, commited keystore files from this folder into the new testDataDirPath()
//
// Afterwards the test datadir path will look like:
//
//	"/tmp/tbb_test945924586/keystore/test_andrej--3eb92807f1f91a8d4d85bc908c7f86dcddb1df57"
//	"/tmp/tbb_test945924586/keystore/test_babayaga--6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8"
func copyKeystoreFilesIntoTestDataDirPath(dataDir string) error {
	andrejSrcKs, err := os.Open(testKsAndrejFile)
	if err != nil {
		return err
	}
	defer andrejSrcKs.Close()

	ksDir := filepath.Join(wallet.GetKeystoreDirPath(dataDir))

	err = os.Mkdir(ksDir, 0777)
	if err != nil {
		return err
	}

	andrejDstKs, err := os.Create(filepath.Join(ksDir, testKsAndrejFile))
	if err != nil {
		return err
	}
	defer andrejDstKs.Close()

	_, err = io.Copy(andrejDstKs, andrejSrcKs)
	if err != nil {
		return err
	}

	babayagaSrcKs, err := os.Open(testKsBabaYagaFile)
	if err != nil {
		return err
	}
	defer babayagaSrcKs.Close()

	babayagaDstKs, err := os.Create(filepath.Join(ksDir, testKsBabaYagaFile))
	if err != nil {
		return err
	}
	defer babayagaDstKs.Close()

	_, err = io.Copy(babayagaDstKs, babayagaSrcKs)
	if err != nil {
		return err
	}

	return nil
}

// setupTestNodeDir creates a default testing node directory with 2 keystore accounts
//
// Remember to remove the dir once test finishes: defer fs.RemoveDir(dataDir)
func setupTestNodeDir() (dataDir string, andrej, babaYaga common.Address, err error) {
	babaYaga = database.NewAccount(testKsBabaYagaAccount)
	andrej = database.NewAccount(testKsAndrejAccount)

	dataDir, err = getTestDataDirPath()
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	genesisBalances := make(map[common.Address]uint)
	genesisBalances[andrej] = 1000000
	genesis := database.Genesis{Balances: genesisBalances}
	genesisJson, err := json.Marshal(genesis)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	err = database.InitDataDirIfNotExists(dataDir, genesisJson)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	err = copyKeystoreFilesIntoTestDataDirPath(dataDir)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	return dataDir, andrej, babaYaga, nil
}
