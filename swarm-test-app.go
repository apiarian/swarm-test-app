package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/net/context"

	ma "gx/ipfs/QmcobAGsCjYt5DXoq9et9L8yR8er7o7Cu3DTvpaq12jYSz/go-multiaddr"

	"gx/ipfs/QmNefBbWHR9JEiP3KDVqZsBLQVRmH3GBG2D2Ke24SsFqfW/go-libp2p/p2p/metrics"
	"gx/ipfs/QmNefBbWHR9JEiP3KDVqZsBLQVRmH3GBG2D2Ke24SsFqfW/go-libp2p/p2p/net"
	"gx/ipfs/QmNefBbWHR9JEiP3KDVqZsBLQVRmH3GBG2D2Ke24SsFqfW/go-libp2p/p2p/net/swarm"
	"gx/ipfs/QmNefBbWHR9JEiP3KDVqZsBLQVRmH3GBG2D2Ke24SsFqfW/go-libp2p/p2p/peer"

	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func main() {
	ipfs_path := os.Getenv("IPFS_PATH")
	if ipfs_path == "" {
		ipfs_path = "~/.ipfs"
	}
	fmt.Println("IPFS_PATH:", ipfs_path)

	conf, err := fsrepo.ConfigAt(ipfs_path)
	if err != nil {
		panic(err)
	}
	fmt.Println("got config:", conf)

	peer_id, err := peer.IDB58Decode(conf.Identity.PeerID)
	if err != nil {
		panic(err)
	}

	to_call_peer_id, err := peer.IDB58Decode("QmYLynW5fGEhcNrxW5M9AT2z868feQ7tTFGFuF22FyEYJu")
	if err != nil {
		panic(err)
	}
	to_call_addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/50000")

	private_key, err := conf.Identity.DecodePrivateKey("")
	if err != nil {
		panic(err)
	}

	fmt.Println("peer_id", peer_id.Pretty())

	fmt.Println("swarm addresses", conf.Addresses.Swarm)

	multiaddrs := make([]ma.Multiaddr, len(conf.Addresses.Swarm))
	for i, v := range conf.Addresses.Swarm {
		multiaddr, err := ma.NewMultiaddr(v)
		if err != nil {
			panic(err)
		}
		multiaddrs[i] = multiaddr
	}
	fmt.Println("swarm multiaddrs", multiaddrs)

	pstore := peer.NewPeerstore()
	err = pstore.AddPrivKey(peer_id, private_key)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// s, err := swarm.NewSwarm(ctx, multiaddrs, peer_id, pstore, metrics.NewBandwidthCounter())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("got a swarm?", s)
	sn, err := swarm.NewNetwork(ctx, multiaddrs, peer_id, pstore, metrics.NewBandwidthCounter())
	if err != nil {
		panic(err)
	}
	fmt.Println("got a network?", sn)

	sn.SetStreamHandler(func(st net.Stream) {
		fmt.Println("handling some stream")
		out, err := ioutil.ReadAll(st)
		if err != nil {
			panic(err)
		}
		fmt.Println("received", string(out))
	})

	for {
		fmt.Println("conns:", sn.Conns())
		fmt.Println("peerstore:", sn.Peerstore().Peers())
		fmt.Println("listening on:", sn.Swarm().ListenAddresses())
		time.Sleep(2 * time.Second)
		if to_call_peer_id.Pretty() != peer_id.Pretty() {
			pstore.AddAddr(to_call_peer_id, to_call_addr, peer.PermanentAddrTTL)
			fmt.Println("connections to peer:", sn.ConnsToPeer(to_call_peer_id))
			str, err := sn.NewStream(context.Background(), to_call_peer_id)
			if err != nil {
				panic(err)
			}
			fmt.Println("connections to peer:", sn.ConnsToPeer(to_call_peer_id))
			fmt.Println("about to send the message...")
			n, err := fmt.Fprintln(str, "Hello peer!")
			if err != nil {
				panic(err)
			} else {
				fmt.Println("maybe sent bytes:", n)
			}
			fmt.Println("sent the message.")
			str.Close()
			fmt.Println("closing the stream")
		} else {
		}
	}
}
