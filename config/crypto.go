package config

type crypto struct {
	Salt      []byte
	ARGON2Key argon2Config
	ARGON2ID  argon2Config
}

type argon2Config struct {
	Parallelism uint8
	Memory      uint32
	Rounds      uint32
}
