package config

import (
	"time"
)

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}
