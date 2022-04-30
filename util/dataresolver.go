package util

import (
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"

	"github.com/pavel-v-chernykh/keystore-go/v4"
)

type nonRand struct{}

func (r nonRand) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 1
	}

	return len(p), nil
}

func zeroing(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func readPrivateKey(filepath string) []byte {
	pkPEM, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// privateKeyStr.Set(string(pkPEM))

	b, _ := pem.Decode(pkPEM)
	if b == nil {
		log.Fatal("should have at least one pem block")
	}

	if b.Type != "PRIVATE KEY" {
		log.Fatal("should be a private key")
	}

	return b.Bytes
}

func readCertificate(filepath string) []byte {
	pkPEM, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// certStr.Set(string(pkPEM))

	b, _ := pem.Decode(pkPEM)
	if b == nil {
		log.Fatal("should have at least one pem block")
	}

	if b.Type != "CERTIFICATE" {
		log.Fatal("should be a certificate")
	}

	return b.Bytes
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	err = ks.Store(f, password)
	if err != nil {
		panic(err)
	}
}
