package utils

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"strings"

	"github.com/leonklingele/passphrase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func GenerateID(n int) string {
	passphrase.Separator = "-"
	id, _ := passphrase.Generate(n)

	idFields := strings.Split(id, "-")
	for i := range idFields {
		idFields[i] = titleCaser.String(idFields[i])
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
