// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"ruixuego/bytepool"
)

var (
	bytePools, _ = bytepool.NewMultiRatedBytePool(4, 10, 1024)
)

// AESData AES 密钥缓存结构
type AESData struct {
	Block cipher.Block
	IV    []byte
}

func (d *AESData) Encrypter() cipher.BlockMode {
	return cipher.NewCBCEncrypter(d.Block, d.IV)
}

func (d *AESData) Decrypter() cipher.BlockMode {
	return cipher.NewCBCDecrypter(d.Block, d.IV)
}

// NewAESData 生成加密密钥数据块
func NewAESData(key []byte) (*AESData, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &AESData{
		block,
		key[:block.BlockSize()],
	}, nil
}

// NewAESDataWityHex 基于 64 长度的 16 进制密钥生成加密密钥数据块
func NewAESDataWityHex(key string) (*AESData, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return NewAESData(k)
}

// AesEncrypt AES-CBC加密
func AesEncrypt(content []byte, key *AESData) []byte {
	content = PKCS7Padding(content, key.Block.BlockSize())
	// 向量 (key[:blockSize]) 是密钥的前 blockSize (16) 个字节
	ciphertext := make([]byte, len(content))
	key.Encrypter().CryptBlocks(ciphertext, content)
	return ciphertext
}

// AesDecrypt AES-CBC解密
func AesDecrypt(ciphertext []byte, key *AESData) []byte {
	b := make([]byte, len(ciphertext))
	key.Decrypter().CryptBlocks(b, ciphertext)
	b = PKCS7UnPadding(b)
	return b
}

// AesEncryptWithPool AES-CBC加密
func AesEncryptWithPool(content []byte, key *AESData) []byte {
	content = PKCS7Padding(content, key.Block.BlockSize())
	l := len(content)
	b := bytePools.Get(l)
	key.Encrypter().CryptBlocks(b[:l], content)
	return b[:l]
}

// AesDecryptWithPool AES-CBC解密
func AesDecryptWithPool(ciphertext []byte, key *AESData) []byte {
	l := len(ciphertext)
	b := bytePools.Get(l)
	key.Decrypter().CryptBlocks(b[:l], ciphertext)
	b = PKCS7UnPadding(b[:l])
	return b
}

// AesEncryptBase64String AES-CBC加密字符串
func AesEncryptBase64String(content string, key *AESData) string {
	b := AesEncryptWithPool(StringToBytes(content), key)
	s := base64.StdEncoding.EncodeToString(b)
	bytePools.Put(b)
	return s
}

// AesDecryptBase64String AES-CBC解密字符串
func AesDecryptBase64String(content string, key *AESData) string {
	l := base64.StdEncoding.DecodedLen(len(content))
	b64 := bytePools.Get(l)
	n, err := base64.StdEncoding.Decode(b64[:l], StringToBytes(content))
	if err != nil {
		return ""
	}
	b := AesDecryptWithPool(b64[:n], key)
	s := string(b)
	bytePools.Put(b)
	bytePools.Put(b64)
	return s
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(oriData []byte) []byte {
	length := len(oriData)
	if length == 0 {
		return []byte{}
	}
	unpadding := int(oriData[length-1])
	if length < unpadding {
		return []byte{}
	}
	return oriData[:(length - unpadding)]
}

// GenerateAESKey 生成 AES-256 密钥
func GenerateAESKey() (string, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", key), err
}
