package rsaoeap

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

// RSA implements the crypto.Cipher interface using RSA public/private key algorithm
// as specified in PKCS #1. Messages are encrypted with the public key and can only be
// decrypted using the private key. RSA objects must have a public key but the private
// key is only required for decryption.
type RSA struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

// New creates an RSA Crypto handler with the specified key pair. If the cipher is only
// being used for encryption, simply pass the public key: New(pub *rsa.PublicKey); If
// the cipher is being used for decryption, then pass the private key:
// New(key *rsa.PrivateKey).
func New(key interface{}) (_ *RSA, err error) {
	switch t := key.(type) {
	case *rsa.PublicKey:
		return &RSA{pub: t, priv: nil}, nil
	case *rsa.PrivateKey:
		return &RSA{pub: &t.PublicKey, priv: t}, nil
	default:
		return nil, fmt.Errorf("could not create RSA cipher from %T", t)
	}
}

// Encrypt the message using the public key.
func (c *RSA) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	hash := sha512.New()
	ciphertext, err = rsa.EncryptOAEP(hash, rand.Reader, c.pub, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Decrypt the message using the private key.
func (c *RSA) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if c.priv == nil {
		return nil, errors.New("private key required for decryption")
	}

	hash := sha512.New()
	plaintext, err = rsa.DecryptOAEP(hash, rand.Reader, c.priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// EncryptionAlgorithm returns the name of the algorithm for adding to the Transaction.
func (c *RSA) EncryptionAlgorithm() string {
	return "RSA-OAEP-SHA512"
}

// PublicKeySignature implements KeyIdentifier by computing a base64 encoded SHA-256
// hash of the public key serialized as a PKIX public key without PEM encoding. This is
// a prototype method of computing the public key signature and may not match other
// external signature computation methods.
// TODO: verify that this method matches openssl or GitHub public key identification.
func (c *RSA) PublicKeySignature() (_ string, err error) {
	var data []byte
	if data, err = x509.MarshalPKIXPublicKey(c.pub); err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return fmt.Sprintf("SHA256:%s", base64.RawStdEncoding.EncodeToString(sum[:])), nil
}
