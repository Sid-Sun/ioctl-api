package config

type Crypto struct {
	Salt      []byte
	ARGON2Key ARGON2Config
	ARGON2ID  ARGON2Config
}

type ARGON2Config struct {
	Parallelism uint8
	Memory      uint32
	Rounds      uint32
}
