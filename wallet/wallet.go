package wallet

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"bytes"

	"golang.org/x/crypto/ripemd160"
	"github.com/michain/dotcoin/base58"
	"fmt"
)


// version pubkey for bitcoin, version = 0
const version = byte(0x00)
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// newWallet creates and returns a Wallet
func newWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// GetStringAddress returns wallet address string format
func (w Wallet) GetStringAddress() string{
	return fmt.Sprintf("%s", w.GetAddress())
}

// GetAddress returns wallet address
// 1.hashes public key - ripemd160(sha256(public key))
// 2.connect version to the PubKeyHash header.
// 3.get checksum，use first 4 bytes
// 4.connect checksum to the end of the hash data
// 5.base58 data
func (w Wallet) GetAddress() []byte {
	//1.hashes public key - ripemd160(sha256(public key))
	pubKeyHash := HashPublicKey(w.PublicKey)
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return []byte(address)
}


// HashPublicKey hashes public key
// 1.sha256 publick key
// 2.ripemd160(sha256(public key))
// 疑惑：
// RIPEMD160已被sha-256和SHA 512和它们的派生算法取代。
// SHA256和SHA512提供更好的安全和性能比RIPEMD160。
// 使用RIPEMD160仅为与旧的应用程序和数据的兼容性。
// 一种看法：在SAH256基础上进行RIPEMD160，一方面双重hash保护，一方面将生成的散列码由256位缩短为160位
func HashPublicKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
// 1.base58 decode
// 2.get checksum(4 byte)
// 3.get version
// 4.get pubKeyHash
// 5.get checksum，use first 4 bytes
func ValidateAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// GetPubKeyHashFromAddress get PublicKeyHash from address
// 1.base58 decode
// 2.remove version & checksum
func GetPubKeyHashFromAddress(address []byte) []byte{
	pubKeyHash := base58.Decode(string(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return pubKeyHash
}



// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// newKeyPair create private&public key with ecdsa and rand-key
// 1.随机数发生器生成一个私钥
// 2.P-256算法生成公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

