package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var errPKCS7Padding = errors.New("pkcs7pad: bad padding")

// Pad appends PKCS#7 padding to the given buffer such that the resulting slice
// of bytes has a length divisible by the given size. If you are using this
// function to pad a plaintext before encrypting it with a block cipher, the
// size should be equal to the block size of the cipher (e.g., aes.BlockSize).
func Pad(buf []byte, size int) []byte {
	if size < 1 || size > 255 {
		panic(fmt.Sprintf("pkcs7pad: inappropriate block size %d", size))
	}
	i := size - (len(buf) % size)
	return append(buf, bytes.Repeat([]byte{byte(i)}, i)...)
}

// Unpad returns a subslice of the input buffer with trailing PKCS#7 padding
// removed. It checks the correctness of the padding bytes in constant time, and
// returns an error if the padding bytes are malformed.
func Unpad(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return nil, errPKCS7Padding
	}

	// Here be dragons. We're attempting to check the padding in constant
	// time. The only piece of information here which is public is len(buf).
	// This code is modeled loosely after tls1_cbc_remove_padding from
	// OpenSSL.
	padLen := buf[len(buf)-1]
	toCheck := 255
	good := 1
	if toCheck > len(buf) {
		toCheck = len(buf)
	}
	for i := 0; i < toCheck; i++ {
		b := buf[len(buf)-1-i]

		outOfRange := subtle.ConstantTimeLessOrEq(int(padLen), i)
		equal := subtle.ConstantTimeByteEq(padLen, b)
		good &= subtle.ConstantTimeSelect(outOfRange, 1, equal)
	}

	good &= subtle.ConstantTimeLessOrEq(1, int(padLen))
	good &= subtle.ConstantTimeLessOrEq(int(padLen), len(buf))

	if good != 1 {
		return nil, errPKCS7Padding
	}

	return buf[:len(buf)-int(padLen)], nil
}

func CBCEncrypt(key string, content []byte) string {
	keyBytes := []byte(key)

	// create a new iv
	iv := make([]byte, aes.BlockSize)

	// read random bytes into iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// create a new cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	// create a new cbc
	cbc := cipher.NewCBCEncrypter(block, iv)

	// pad content
	paddedContent := Pad(content, aes.BlockSize)

	// create a new array
	ciphertext := make([]byte, len(paddedContent))

	// encrypt content
	cbc.CryptBlocks(ciphertext, paddedContent)

	// ciphertext to base64
	encodedtext := base64.StdEncoding.EncodeToString(ciphertext)

	// append iv to encoddedtext
	encodedtext = base64.StdEncoding.EncodeToString(iv) + ":" + encodedtext

	return encodedtext
}

func CBCEncryptWithString(key string, content string) string {
	return CBCEncrypt(key, []byte(content))
}

func CBCDecrypt(key string, content string) string {
	// convert key to bytes
	keyBytes := []byte(key)

	// split content by :
	splitContent := bytes.Split([]byte(content), []byte(":"))

	// decode iv
	iv, err := base64.StdEncoding.DecodeString(string(splitContent[0]))
	if err != nil {
		panic(err)
	}

	// decode ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(string(splitContent[1]))
	if err != nil {
		panic(err)
	}

	// create a new cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	// create a new cbc
	cbc := cipher.NewCBCDecrypter(block, iv)

	// create a new array
	plaintext := make([]byte, len(ciphertext))

	// decrypt content
	cbc.CryptBlocks(plaintext, ciphertext)

	// unpad content
	unpaddedContent, err := Unpad(plaintext)
	if err != nil {
		panic(err)
	}

	return string(unpaddedContent)
}

func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

// see if file content is identical to content, if not, replace it
func write_if_changes(filename string, content []byte) {
	// check if file exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("~~")

		// file does not exist, create it
		os.Create(filename)
	}

	// read file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// compare file content to content
	if bytes.Compare(fileContent, content) != 0 {
		// log that file content is different
		fmt.Println("[" + filename + "] File content is different, writing new content")

		// write content to file
		err = ioutil.WriteFile(filename, content, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		// log that file content is the same
		fmt.Println("[" + filename + "] File content is the same, not writing new content")
	}

}
