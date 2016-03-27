package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func main() {
	ipfs_path := os.Getenv("IPFS_PATH")
	if ipfs_path == "" {
		ipfs_path = "~/.ipfs"
	}
	fmt.Println("IPFS_PATH:", ipfs_path)

	r, err := fsrepo.Open(ipfs_path)
	if err != nil {
		panic(err)
	}

	head_node_id_path := path.Join(ipfs_path, "head_node.id")
	head_node_file, err := os.Open(head_node_id_path)
	if err != nil {
		panic(err)
	}
	defer head_node_file.Close()
	head_node_id := ""
	data := make([]byte, 100)
	count, err := head_node_file.Read(data)
	if err != nil {
		panic(err)
	}
	head_node_id = string(data[:count])
	fmt.Println("head_node_id:", head_node_id)

	conf, err := r.Config()
	if err != nil {
		panic(err)
	}
	swarm_address_parts := strings.Split(conf.Addresses.Swarm[0], "/")
	swarm_address_port := swarm_address_parts[len(swarm_address_parts)-1]
	swarm_address_port_i, _ := strconv.Atoi(swarm_address_port)
	head_node_port_i := swarm_address_port_i - 10000
	head_node_port := strconv.Itoa(head_node_port_i)
	head_node_location := "/ip4/127.0.0.1/tcp/" + head_node_port + "/ipfs/" + head_node_id
	new_bootstrap_peer_strings := []string{head_node_location}
	new_bootstrap_peers, err := config.ParseBootstrapPeers(new_bootstrap_peer_strings)
	if err != nil {
		panic(err)
	}
	fmt.Println("boostrap peers:", new_bootstrap_peers)
	conf.SetBootstrapPeers(new_bootstrap_peers)
	r.SetConfig(conf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	nd, err := core.NewNode(ctx, cfg)
	if err != nil {
		panic(err)
	}

	list, err := corenet.Listen(nd, "/app/interplanetary-game-system")
	if err != nil {
		panic(err)
	}

	fmt.Printf("I'm a peer: %s\n", nd.Identity.Pretty())

	for {
		con, err := list.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer con.Close()

		fmt.Fprintln(con, "Hello, this is the interplanetary-game-system")
		fmt.Printf("Connection from: %s\n", con.Conn().RemotePeer())
	}
}
