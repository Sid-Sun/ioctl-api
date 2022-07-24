package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/fitant/xbin-api/config"
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
	// Generate cipher
	c, _ := aes.NewCipher(key)

	// use CFB to encrypt full data
	ciphertext, iv = gcmEncrypt(data, c)

	// prepend salt to the data
	return ciphertext, iv, (*salt)[:]
}

func Decrypt(data []byte, salt []byte, iv []byte, key []byte) []byte {
	// Derive Key for decryption from ID using Argon2
	key = GenKey(key, salt)

	// Generate cipher
	c, _ := aes.NewCipher(key)

	// Send IV and data bits to decrypt via CFB
	// returned data does not have salt
	return gcmDecrypt(data, iv, c)
}

func cfbEncrypt(data []byte, blockCipher cipher.Block) []byte {
	// Create dst with length of cipher blocksize + data length
	// And initialize first BlockSize bytes pseudorandom for IV
	dst := make([]byte, blockCipher.BlockSize()+len(data))

	// Read random values from crypto/rand for CFB initialization vector
	// Error can be safely ignored
	io.ReadFull(rand.Reader, dst[:blockCipher.BlockSize()])

	// dst from 0 to blockSize is the IV
	cfb := cipher.NewCFBEncrypter(blockCipher, dst[:blockCipher.BlockSize()])
	cfb.XORKeyStream(dst[blockCipher.BlockSize():], data)
	return dst
}

func gcmDecrypt(data []byte, nonce []byte, blockCipher cipher.Block) []byte {
	// Create CFB Decrypter with cipher, instantiating with IV (first blockSize blocks of data)
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
	return data
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

	// prepend generated nonce to data
	return ciphertext, nonce
}
