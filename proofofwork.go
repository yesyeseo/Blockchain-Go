package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

var (
	maxNonce = math.MaxInt64
)

const targetBits = 24

type ProofOfWork struct {
	// contatin block pointer and target pointer
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// initalize big.Int as 1
	target.Lsh(target, uint(256-targetBits))
	// left shift as 256-targetbits bits
	// 256 is bit length of SHA-256 hash that we would use

	pow := &ProofOfWork{b, target}
	return pow
}

// data for computing hash
// return fields of block, target, nonce value
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	// init variables
	var hashInt big.Int
	// integer value of hash
	var hash [32]byte
	nonce := 0
	// counter

	fmt.Printf("Mining the block contatining \"%s\"\n", pow.block.Data)

	// infinite loop: maximum maxNonce
	// set max <- to prevent overflow of nonce
	// the value is same with math.MaxInt64
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		// set data
		hash = sha256.Sum256(data)
		// hashing SHA-256
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		// convent hash value to Big Integer
		if hashInt.Cmp(pow.target) == -1 {
			// compare hash value and target value
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

// delete SetHash method of Block
// modify NewBlock function
