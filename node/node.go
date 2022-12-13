package node

import (
	"Blockchain_Go/database"
	"fmt"
	"net/http"
)

const DefaultHTTPort = 8080

type Node struct {
	dataDir    string
	port       uint64
	state      *database.State
	knownPeers []PeerNode
}
type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	return &Node{
		dataDir:    dataDir,
		port:       port,
		knownPeers: []PeerNode{bootstrap},
	}
}
func (n *Node) Run() error {
	fmt.Println(fmt.Sprintf("Listening on HTTP port: %d", n.port))

	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})

	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, isActive}
}
