// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:11 下午 - UTC/GMT+08:00

package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"unsafe"
)

// []byte转string
func bytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) // nolint
}

// string转[]byte
func strToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // nolint
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h)) // nolint
}

// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
func encodeByBytes(key, str []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	str = pkcs5Padding(str, blockSize)
	// str = zeroPadding(str, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(str))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := str
	blockMode.CryptBlocks(crypted, str)
	return hex.EncodeToString(crypted), nil
}

func decodeByBytes(key, str []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(str))
	// origData := crypted
	blockMode.CryptBlocks(origData, str)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return bytesToStr(origData), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 解码得到storage id
func decodeStorageID(hexStr, key string) (string, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	// 将解码后的sid值解密成uuid
	sid, err := decodeByBytes(strToByte(key), b)
	if err != nil {
		return sid, err
	}

	return sid, nil
}
