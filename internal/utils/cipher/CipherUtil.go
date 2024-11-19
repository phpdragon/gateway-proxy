package cipher

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func Md5(sourceData string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(sourceData)))
}

// RsaBase64Encrypt 加密
func RsaBase64Encrypt(publicKey string, sourceData string) (string, error) {
	encryptData, err := RsaEncrypt(publicKey, []byte(sourceData))
	if nil != err {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encryptData), nil
}

// RsaEncrypt 加密
func RsaEncrypt(publicKey string, sourceData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pubic := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pubic, sourceData)
}

// RsaBase64Decrypt 解密
func RsaBase64Decrypt(privateKey string, base64Str string) (string, error) {
	//base64解密
	encrypted, err := base64.RawURLEncoding.DecodeString(base64Str)
	if nil != err {
		return "", err
	}

	sourceData, err := RsaDecrypt(privateKey, encrypted)
	if nil != err {
		return "", err
	}
	return string(sourceData), nil
}

// RsaDecrypt 解密
func RsaDecrypt(privateKey string, encryptedData []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, private, encryptedData)
}
