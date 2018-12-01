package core

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var (
	//StrKey 加密key
	StrKey = "Sl"
)

func getKey() []byte {

	keyLen := len(StrKey)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(StrKey)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//Encrypt 加密字符串
func Encrypt(strMesg string) (string, error) {
	key := getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
	str := base64.StdEncoding.EncodeToString(encrypted)
	return str, nil
}

//Decrypt 解密字符串
func Decrypt(src string) (strDesc string, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	srcByte, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	key := getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(srcByte))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, srcByte)
	return string(decrypted), nil
}
