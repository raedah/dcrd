// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainhash

import (
	"encoding/hex"
	"fmt"
)

// HashSize of array used to store sha hashes.  See Hash.
const HashSize = 32

// MaxHashStringSize is the maximum length of a Hash hash string.
const MaxHashStringSize = HashSize * 2

// ErrHashStrSize describes an error that indicates the caller specified a hash
// string that has too many characters.
var ErrHashStrSize = fmt.Errorf("max hash string length is %v bytes", MaxHashStringSize)

// Hash is used in several of the bitcoin messages and common structures.  It
// typically represents the double sha256 of data.
type Hash [HashSize]byte

// String returns the Hash as the hexadecimal string of the byte-reversed
// hash.
func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

// Bytes returns the bytes which represent the hash as a byte slice.
//
// NOTE: This makes a copy of the bytes and should have probably been named
// CloneBytes.  It is generally cheaper to just slice the hash directly thereby
// reusing the same bytes rather than calling this method.
func (hash *Hash) Bytes() []byte {
	newHash := make([]byte, HashSize)
	copy(newHash, hash[:])

	return newHash
}

// SetBytes sets the bytes which represent the hash.  An error is returned if
// the number of bytes passed in is not HashSize.
func (hash *Hash) SetBytes(newHash []byte) error {
	nhlen := len(newHash)
	if nhlen != HashSize {
		return fmt.Errorf("invalid sha length of %v, want %v", nhlen,
			HashSize)
	}
	copy(hash[:], newHash)

	return nil
}

// IsEqual returns true if target is the same as hash.
func (hash *Hash) IsEqual(target *Hash) bool {
	if hash == nil && target == nil {
		return true
	}
	if hash == nil || target == nil {
		return false
	}
	return *hash == *target
}

// NewHash returns a new Hash from a byte slice.  An error is returned if
// the number of bytes passed in is not HashSize.
func NewHash(newHash []byte) (*Hash, error) {
	var sh Hash
	err := sh.SetBytes(newHash)
	if err != nil {
		return nil, err
	}
	return &sh, err
}

// NewHashFromStr creates a Hash from a hash string.  The string should be
// the hexadecimal string of a byte-reversed hash, but any missing characters
// result in zero padding at the end of the Hash.
func NewHashFromStr(hash string) (*Hash, error) {
	// Return error if hash string is too long.
	if len(hash) > MaxHashStringSize {
		return nil, ErrHashStrSize
	}

	// Hex decoder expects the hash to be a multiple of two.
	if len(hash)%2 != 0 {
		hash = "0" + hash
	}

	// Convert string hash to bytes.
	buf, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	// Un-reverse the decoded bytes, copying into in leading bytes of a
	// Hash.  There is no need to explicitly pad the result as any
	// missing (when len(buf) < HashSize) bytes from the decoded hex string
	// will remain zeros at the end of the Hash.
	var ret Hash
	blen := len(buf)
	mid := blen / 2
	if blen%2 != 0 {
		mid++
	}
	blen--
	for i, b := range buf[:mid] {
		ret[i], ret[blen-i] = buf[blen-i], b
	}
	return &ret, nil
}
