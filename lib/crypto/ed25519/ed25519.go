// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package ed25519

import (
	ed25519 "crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/crypto"

	"github.com/ChainSafe/go-schnorrkel"
)

// PublicKeyLength is the fixed Public Key Length
const PublicKeyLength int = 32

// SeedLength is the fixed Seed Length
const SeedLength int = 32

// PrivateKeyLength is the fixed Private Key Length
const PrivateKeyLength int = 64

// SignatureLength is the fixed Signature Length
const SignatureLength int = 64

// Keypair is a ed25519 public-private keypair
type Keypair struct {
	public  *PublicKey
	private *PrivateKey
}

// PrivateKey is the ed25519 Private Key
type PrivateKey ed25519.PrivateKey

// PublicKey is the ed25519 Public Key
type PublicKey ed25519.PublicKey

// PublicKeyBytes is an encoded ed25519 public key
type PublicKeyBytes [PublicKeyLength]byte

// VerifySignature verifies a signature given a public key and a message
func VerifySignature(publicKey, signature, message []byte) error {
	pubKey, err := NewPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("ed25519: %w", err)
	}

	ok, err := pubKey.Verify(message, signature)
	if err != nil {
		return fmt.Errorf("ed25519: %w", err)
	} else if !ok {
		return fmt.Errorf("ed25519: %w: for message 0x%x, signature 0x%x and public key 0x%x",
			crypto.ErrSignatureVerificationFailed, message, signature, publicKey)
	}

	return nil
}

// String returns the PublicKeyBytes formatted as a hex string
func (b PublicKeyBytes) String() string {
	pk := [PublicKeyLength]byte(b)
	return common.BytesToHex(pk[:])
}

// Encode returns the SCALE encoding of PublicKeyBytes
func (b PublicKeyBytes) Encode() ([]byte, error) {
	return b[:], nil
}

// Decode returns the SCALE decoded PublicKeyBytes
func (b PublicKeyBytes) Decode(r io.Reader) ([PublicKeyLength]byte, error) {
	_, err := r.Read(b[:])
	return b, err
}

// SignatureBytes is a ed25519 signature
type SignatureBytes [SignatureLength]byte

// NewKeypair returns an Ed25519 keypair given a ed25519 private key
func NewKeypair(priv ed25519.PrivateKey) *Keypair {
	pubkey := PublicKey(priv.Public().(ed25519.PublicKey))
	privkey := PrivateKey(priv)
	return &Keypair{
		public:  &pubkey,
		private: &privkey,
	}
}

// NewKeypairFromPrivate returns a ed25519 Keypair given a *ed25519.PrivateKey
func NewKeypairFromPrivate(priv *PrivateKey) (*Keypair, error) {
	pub, err := priv.Public()
	if err != nil {
		return nil, err
	}
	return &Keypair{
		public:  pub.(*PublicKey),
		private: priv,
	}, nil
}

// NewKeypairFromSeed generates a new ed25519 keypair from a 32 byte seed
func NewKeypairFromSeed(seed []byte) (*Keypair, error) {
	if len(seed) != SeedLength {
		return nil, fmt.Errorf("cannot generate key from seed: seed is not 32 bytes long")
	}
	edpriv := ed25519.NewKeyFromSeed(seed)
	return NewKeypair(edpriv), nil
}

// NewKeypairFromPrivateKeyString returns a Keypair given a 0x prefixed private key string
func NewKeypairFromPrivateKeyString(in string) (*Keypair, error) {
	privBytes, err := common.HexToBytes(in)
	if err != nil {
		return nil, err
	}

	return NewKeypairFromSeed(privBytes)
}

// NewKeypairFromMnenomic returns a new Keypair using the given mnemonic and password.
func NewKeypairFromMnenomic(mnemonic, password string) (*Keypair, error) {
	seed, err := schnorrkel.SeedFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, err
	}
	return NewKeypairFromSeed(seed[:32])
}

// GenerateKeypair returns a new ed25519 keypair
func GenerateKeypair() (*Keypair, error) {
	buf := make([]byte, SeedLength)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	priv := ed25519.NewKeyFromSeed(buf)

	return NewKeypair(priv), nil
}

// NewPublicKey returns an ed25519 public key that consists of the input bytes
// Input length must be 32 bytes
func NewPublicKey(in []byte) (*PublicKey, error) {
	if len(in) != PublicKeyLength {
		return nil, fmt.Errorf("cannot create public key: input is not 32 bytes")
	}

	pub := PublicKey(ed25519.PublicKey(in))
	return &pub, nil
}

// NewPrivateKey returns an ed25519 private key that consists of the input bytes
// Input length must be 64 bytes
func NewPrivateKey(in []byte) (*PrivateKey, error) {
	if len(in) != PrivateKeyLength {
		return nil, fmt.Errorf("cannot create private key: input is not 64 bytes")
	}

	priv := PrivateKey(ed25519.PrivateKey(in))
	return &priv, nil
}

// Verify returns true if the signature is valid for the given message and public key, false otherwise
func Verify(pub *PublicKey, msg, sig []byte) (bool, error) {
	if len(sig) != SignatureLength {
		return false, errors.New("invalid signature length")
	}

	return ed25519.Verify(ed25519.PublicKey(*pub), msg, sig), nil
}

// Type returns Ed25519Type
func (kp *Keypair) Type() crypto.KeyType {
	return crypto.Ed25519Type
}

// Sign uses the keypair to sign the message using the ed25519 signature algorithm
func (kp *Keypair) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(ed25519.PrivateKey(*kp.private), msg), nil
}

// Public returns the keypair's public key
func (kp *Keypair) Public() crypto.PublicKey {
	return kp.public
}

// Private returns the keypair's private key
func (kp *Keypair) Private() crypto.PrivateKey {
	return kp.private
}

// Sign uses the ed25519 signature algorithm to sign the message
func (k *PrivateKey) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(ed25519.PrivateKey(*k), msg), nil
}

// Public returns the public key corresponding to the ed25519 private key
func (k *PrivateKey) Public() (crypto.PublicKey, error) {
	kp := NewKeypair(ed25519.PrivateKey(*k))
	return kp.Public(), nil
}

// Encode returns the bytes underlying the ed25519 PrivateKey
func (k *PrivateKey) Encode() []byte {
	return []byte(ed25519.PrivateKey(*k))
}

// Decode turns the input bytes into a ed25519 PrivateKey
// the input must be 64 bytes, or the function will return an error
func (k *PrivateKey) Decode(in []byte) error {
	priv, err := NewPrivateKey(in)
	if err != nil {
		return err
	}
	*k = *priv
	return nil
}

// Hex will return PrivateKey Hex
func (k *PrivateKey) Hex() string {
	enc := k.Encode()
	h := hex.EncodeToString(enc)
	return "0x" + h
}

// Verify checks that Ed25519PublicKey was used to create the signature for the message
func (k *PublicKey) Verify(msg, sig []byte) (bool, error) {
	if len(sig) != SignatureLength {
		return false, errors.New("invalid signature length")
	}
	return ed25519.Verify(ed25519.PublicKey(*k), msg, sig), nil
}

// Encode returns the encoding of the ed25519 PublicKey
func (k *PublicKey) Encode() []byte {
	return []byte(ed25519.PublicKey(*k))
}

// Decode turns the input bytes into an ed25519 PublicKey
// the input must be 32 bytes, or the function will return and error
func (k *PublicKey) Decode(in []byte) error {
	pub, err := NewPublicKey(in)
	if err != nil {
		return err
	}
	*k = *pub
	return nil
}

// Address returns the ss58 address for this public key
func (k *PublicKey) Address() common.Address {
	return crypto.PublicKeyToAddress(k)
}

// Hex returns the public key as a '0x' prefixed hex string
func (k *PublicKey) Hex() string {
	enc := k.Encode()
	h := hex.EncodeToString(enc)
	return "0x" + h
}

// AsBytes returns the public key as PublicKeyBytes
func (k *PublicKey) AsBytes() PublicKeyBytes {
	b := [PublicKeyLength]byte{}
	copy(b[:], k.Encode())
	return b
}

// NewSignatureBytes returns a SignatureBytes given a byte array
func NewSignatureBytes(in []byte) SignatureBytes {
	sig := SignatureBytes{}
	copy(sig[:], in)
	return sig
}
