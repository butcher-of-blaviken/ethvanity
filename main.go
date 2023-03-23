package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/term"
)

var (
	command         = flag.String("cmd", "generate", "command to run. valid options: generate, verify")
	numWorkers      = flag.Int("num-workers", runtime.NumCPU(), "number of workers to use. for generate only.")
	desiredPattern  = flag.String("desired-pattern", "", "desired prefix or suffix pattern.  for generate only.")
	patternPosition = flag.String("pattern-position", "prefix", "whether the pattern should be prefix or suffix.  for generate only.")
	outFile         = flag.String("o", "out.txt", "output file to store info in")
	verbose         = flag.Bool("verbose", false, "whether to log progress or not")
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
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("address: %s\n", o.address.Hex()))
		f.WriteString(fmt.Sprintf("private key hex: %s\n", hexutil.Encode(o.priv)))
		f.Sync()
	case "verify":
		fmt.Print("Enter private key hex (without 0x): ")
		hexBytes, err := term.ReadPassword(0)
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		privKey, err := crypto.HexToECDSA(string(hexBytes))
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
			if *verbose && step%10000 == 0 {
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
