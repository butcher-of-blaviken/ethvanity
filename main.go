package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	command         = flag.String("cmd", "generate", "command to run. valid options: generate, verify")
	numWorkers      = flag.Int("num-workers", runtime.NumCPU(), "number of workers to use. for generate only.")
	desiredPattern  = flag.String("desired-pattern", "", "desired prefix or suffix pattern.  for generate only.")
	patternPosition = flag.String("pattern-position", "prefix", "whether the pattern should be prefix or suffix.  for generate only.")
	privHex         = flag.String("priv-hex", "", "private key in hex, to verify the output address")
)

func main() {
	flag.Parse()
	switch *command {
	case "generate":
		out := make(chan out, 1_000)
		done := make(chan struct{})
		for i := 0; i < *numWorkers; i++ {
			go producer(out, done)
		}
		o := consumer(out, done, *desiredPattern, *patternPosition)
		fmt.Println("address:", o.address.Hex())
		fmt.Println("private key:", hexutil.Encode(o.priv))
	case "verify":
		privKey, err := crypto.HexToECDSA(*privHex)
		if err != nil {
			log.Fatal(err)
		}
		address := crypto.PubkeyToAddress(privKey.PublicKey)
		fmt.Println("address:", address.Hex())
	}

}

type out struct {
	address common.Address
	priv    []byte
}

func producer(o chan out, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			output := newAddress()
			o <- output
		}
	}
}

func consumer(o chan out, done chan struct{}, desiredPattern, patternPosition string) out {
	step := 0
	for {
		select {
		case data := <-o:
			if isValidVanityAddress(data.address, desiredPattern, patternPosition) {
				return data
			}
			if step%10000 == 0 {
				fmt.Println("step:", step, "latest addr:", data.address.Hex())
			}
			step++
		case <-done:
			return out{}
		}
	}
}

func isValidVanityAddress(a common.Address, desiredPattern, patternPosition string) bool {
	if patternPosition == "prefix" {
		return strings.HasPrefix(a.Hex()[2:], desiredPattern)
	} else {
		return strings.HasSuffix(a.Hex()[2:], desiredPattern)
	}
}

func newAddress() out {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return out{
		address: address,
		priv:    privateKeyBytes,
	}
}
