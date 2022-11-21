package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/sid-sun/ioctl-api/config"
	"golang.org/x/crypto/argon2"
)

// Output hash length
const hashLength = 32

func HashID(id []byte) []byte {
	return argon2.IDKey(id,
		config.Cfg.Crypto.Salt,
		config.Cfg.Crypto.ARGON2ID.Rounds,
		config.Cfg.Crypto.ARGON2ID.Memory,
		config.Cfg.Crypto.ARGON2ID.Parallelism,
		hashLength)
}

func GenSalt() *[32]byte {
	salt := new([32]byte)
	io.ReadFull(rand.Reader, salt[:])
	return salt
}

func GenKey(key, salt []byte) []byte {
	return argon2.IDKey(key, salt,
		config.Cfg.Crypto.ARGON2Key.Rounds,
		config.Cfg.Crypto.ARGON2Key.Memory,
		config.Cfg.Crypto.ARGON2Key.Parallelism,
		hashLength)
}

func Encrypt(data []byte, key []byte, salt *[32]byte) (ciphertext []byte, iv []byte, keysalt []byte) {
	c, _ := aes.NewCipher(key)

	ciphertext, iv = gcmEncrypt(data, c)

	return ciphertext, iv, (*salt)[:]
}

func Decrypt(data []byte, salt []byte, iv []byte, key []byte) []byte {
	// Derive Key for decryption from ID using Argon2
	key = GenKey(key, salt)

	// Generate cipher
	c, _ := aes.NewCipher(key)

	// returned data does not have salt
	return gcmDecrypt(data, iv, c)
}

func gcmDecrypt(data []byte, nonce []byte, blockCipher cipher.Block) []byte {
	gcm, err := cipher.NewGCMWithNonceSize(blockCipher, 32)
	if err != nil {
		panic(err)
	}
	// Create variable for storing decrypted note of shorter length taking into account IV
	// decrypted := make([]byte, len(data)-gcm.NonceSize())
	_, err = gcm.Open(data[:0], nonce, data, nil)
	if err != nil {
		panic(err)
	}

	// Remove 16 bits of authentication data from the end
	return data[:len(data)-16]
}

func gcmEncrypt(data []byte, blockCipher cipher.Block) (ciphertext []byte, nonce []byte) {
	// Create CFB Decrypter with cipher, instantiating with IV (first blockSize blocks of data)
	gcm, err := cipher.NewGCMWithNonceSize(blockCipher, 32)
	if err != nil {
		panic(err)
	}

	nonce = make([]byte, 32)
	// Read random values from crypto/rand for CFB initialization vector
	// Error can be safely ignored
	io.ReadFull(rand.Reader, nonce)

	ciphertext = gcm.Seal(data[:0], nonce, data, nil)
	if err != nil {
		panic(err)
	}

	// Output contains 16 bits of authentication data
	return ciphertext, nonce
}
