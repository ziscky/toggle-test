package main

import "os"

const (
	defaultAddr = ":8080"
	defaultDB   = ":memory:"
)

type options struct {
	addr       string
	dbFilePath string
}

func loadOptionsFromEnv() *options {
	opts := &options{
		addr:       defaultAddr,
		dbFilePath: defaultDB,
	}

	addr := os.Getenv("ADDR")
	db := os.Getenv("DB")

	if addr != "" {
		opts.addr = addr
	}
	if db != "" {
		opts.dbFilePath = db
	}
	return opts
}
