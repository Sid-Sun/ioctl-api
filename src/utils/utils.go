package utils

import (
	"bytes"
	"io/ioutil"
	"strings"
	"compress/zlib"

	"github.com/leonklingele/passphrase"
)

func GenerateID(n int) string {
	passphrase.Separator = "-"
	id, _ := passphrase.Generate(n)

	idFields := strings.Split(id, "-")
	for i := range idFields {
		idFields[i] = strings.Title(idFields[i])
	}
	return strings.Join(idFields, "")
}

func Defalte(data []byte) []byte {
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, 7)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func Inflate(data []byte) []byte {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	x, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return x
}
