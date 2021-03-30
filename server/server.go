package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"

	golog "github.com/ipfs/go-log/v2"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

type config struct {
	listenPort int
	insecure   bool
	isServer   bool
	randseed   int64
}

const protocol = "/echo/1.0.0"

// Generate a key pair for this host. We will use it at least to obtain a valid host ID.
func generKeyPair(randseed int64) (crypto.PrivKey, crypto.PubKey, error) {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs

	if randseed != 0 {
		return crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, mrand.New(mrand.NewSource(randseed)))
	}

	return crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
}

func configOptions(cfg config) ([]libp2p.Option, error) {
	privateKey, _, errGen := generKeyPair(cfg.randseed)
	if errGen != nil {
		return []libp2p.Option{}, errors.WithMessage(errGen, "when generating Key Pair")
	}

	result := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", cfg.listenPort)),
		libp2p.Identity(privateKey),
		libp2p.DisableRelay(),
	}

	if cfg.insecure {
		result = append(result, libp2p.NoSecurity)
		return result, nil
	}

	return result, nil
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It won't encrypt the connection if insecure is true.
func makeBasicHost(cfg config) (host.Host, error) {
	options, errOpt := configOptions(cfg)
	if errOpt != nil {
		return nil, errors.WithMessage(errOpt, "when generating configuration options")
	}

	basicHost, err := libp2p.New(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	exePath, _ := os.Executable()
	exeName := exePath[strings.LastIndex(exePath, "/")+1:]

	log.Printf("I am %s. \n", fullAddr)

	if cfg.isServer {
		if cfg.insecure {
			log.Printf("Now run \" ./%s -l %d -d %s -insecure \" on a different terminal. \n", exeName, cfg.listenPort+1, fullAddr)
		} else {
			log.Printf("Now run \" ./%s -l %d -d %s \" on a different terminal. \n", exeName, cfg.listenPort+1, fullAddr)
		}
	}

	return basicHost, nil
}

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(golog.LevelInfo) // Change to INFO for extra info

	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	insecure := flag.Bool("insecure", false, "use an unencrypted connection")
	seed := flag.Int64("seed", 0, "set random seed for id generation")

	flag.Parse()

	if *listenF == 0 {
		*listenF = 5555
		log.Printf("Using port %d for incoming connections. \n", *listenF)
	}

	// Make a host that listens on the given multiaddress
	host, err := makeBasicHost(config{
		listenPort: *listenF,
		insecure:   *insecure,
		isServer:   len(*target) == 0,
		randseed:   *seed,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set a stream handler on host A. /echo/1.0.0 is a user-defined protocol name.
	host.SetStreamHandler(protocol, func(s network.Stream) {
		log.Println("Got a new stream!")

		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	if *target == "" {
		log.Println("Listening for connections")
		select {} // hang forever
	}
	/**** This is where the listener code ends ****/

	// The following code extracts target's peer ID from the given multiaddress.
	addressIPFS, err := ma.NewMultiaddr(*target)
	if err != nil {
		log.Fatalln(err)
	}

	hostID, err := addressIPFS.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerID, err := peer.IDB58Decode(hostID)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", hostID))
	targetAddr := addressIPFS.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	host.Peerstore().AddAddr(peerID, targetAddr, peerstore.PermanentAddrTTL)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	stream, err := host.NewStream(context.Background(), peerID, protocol)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = stream.Write([]byte("Hi, world!\n"))
	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("read reply: %q\n", out)
}

// doEcho reads a line of data from passed stream and writes it back.
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)

	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s\n", str)

	_, err = s.Write([]byte(str[:len(str)-1]))
	return err
}
