package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcURL := getenv("RPC", "http://127.0.0.1:8545")
	datadir := getenv("DATADIR", "/tmp/geth1")
	pass := getenv("PASS", "") // empty password
	n := mustInt(getenv("N", "50"))

	ksDir := filepath.Join(datadir, "keystore")
	keyFile := pickFirstKeyfile(ksDir)

	keyjson, err := os.ReadFile(keyFile)
	must(err)

	key, err := keystore.DecryptKey(keyjson, pass)
	must(err)

	priv := key.PrivateKey
	from := key.Address

	client, err := ethclient.Dial(rpcURL)
	must(err)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chainID, err := client.ChainID(ctx)
	must(err)

	nonce, err := client.PendingNonceAt(ctx, from)
	must(err)

	// self-send (to = from)
	to := from
	value := big.NewInt(1)
	gas := uint64(21000)

	// EIP-1559 fee
	tip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		tip = big.NewInt(0)
	}
	fee, err := client.SuggestGasPrice(ctx)
	if err != nil {
		fee = big.NewInt(1_000_000_000) // 1 gwei fallback
	}

	fmt.Printf("from=%s chainID=%s nonce=%d fee=%s tip=%s\n", from.Hex(), chainID.String(), nonce, fee.String(), tip.String())

	signer := types.LatestSignerForChainID(chainID)

	for i := 0; i < n; i++ {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce + uint64(i),
			GasTipCap: tip,
			GasFeeCap: fee,
			Gas:       gas,
			To:        &to,
			Value:     value,
		})

		signed, err := types.SignTx(tx, signer, (*ecdsa.PrivateKey)(priv))
		must(err)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.SendTransaction(ctx2, signed)
		cancel2()

		if err != nil {
			log.Fatalf("send tx #%d failed: %v", i, err)
		}
		fmt.Printf("sent #%d hash=%s\n", i, signed.Hash())
	}
	fmt.Println("done.")
}

func pickFirstKeyfile(dir string) string {
	ents, err := os.ReadDir(dir)
	must(err)
	var files []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		// geth keystore file usually starts with "UTC--"
		if strings.HasPrefix(name, "UTC--") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	if len(files) == 0 {
		log.Fatalf("no keystore file found in %s", dir)
	}
	sort.Strings(files)
	return files[0]
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustInt(s string) int {
	i := new(big.Int)
	_, ok := i.SetString(s, 10)
	if !ok {
		log.Fatalf("bad int: %q", s)
	}
	return int(i.Int64())
}
