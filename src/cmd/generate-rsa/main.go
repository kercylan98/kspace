package main

import (
	"github.com/kercylan98/kspace/src/pkg/cryptography"
	"os"
	"path/filepath"
)

var (
	KeyBits = 256
)

func main() {
	var (
		err                           error
		keyDirPath                    string
		publicKeyFile, privateKeyFile *os.File
	)

	if keyDirPath, err = filepath.Abs("assets/security/rsa"); err != nil {
		panic(err)
	}

	if err = os.MkdirAll(keyDirPath, os.ModePerm); err != nil {
		panic(err)
	}

	if publicKeyFile, err = os.Create(filepath.Join(keyDirPath, "public.key")); err != nil {
		panic(err)
	}

	if privateKeyFile, err = os.Create(filepath.Join(keyDirPath, "private.key")); err != nil {
		panic(err)
	}

	if err = cryptography.CreateKeys(publicKeyFile, privateKeyFile, KeyBits); err != nil {
		panic(err)
	}
}
