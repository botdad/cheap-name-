package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	signature   = flag.String("selector", "", "selector given")
	prefix      = flag.String("prefix", "", "optional prefix before random string")
	charCount   = flag.Int("chars", 5, "number of random characters to use")
	bytesString = flag.String("bytes", "00000000", "hex bytes to match")
)

func IntPow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func main() {
	const chars = "abcdefghijklmnopqrstuvwxyz_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	cl := len(chars)

	flag.Parse()

	if *signature == "" {
		log.Fatal("no selector given")
	}

	sig := *signature
	pfx := *prefix
	numSearchChars := *charCount

	maxPerms := IntPow(cl, numSearchChars)
	chopped := maxPerms / 4

	wanted, err := hex.DecodeString(*bytesString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to find a match for " + pfx + "*****(" + sig + ")")
	fmt.Println("Kicking off 4 threads")

	started := time.Now()
	for i := 0; i < 4; i++ {
		jmin := i * chopped
		jmax := i*chopped + chopped
		fmt.Println("Thread", i, "range", jmin, jmax)
		go func() {
			for j := jmin; j < jmax; j++ {
				a := make([]byte, numSearchChars)
				for k := 0; k < numSearchChars; k++ {
					a[k] = chars[j/IntPow(cl, k)%cl]
				}

				combined := pfx + string(a) + "(" + sig + ")"

				b := crypto.Keccak256([]byte(combined))[:4]

				if bytes.Equal(wanted, b) {
					fmt.Println("FOUND exact match after",
						time.Since(started))
					fmt.Println(combined, "should match", *bytesString)
					os.Exit(0)
					return
				}
			}
			fmt.Println("Thread completed without match")
		}()
	}

	select {}
}
