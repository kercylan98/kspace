package cryptography

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
)

const (
	// RsaAlgorithmSign RSA算法符号
	RsaAlgorithmSign = crypto.SHA256
)

type RSA struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// CreateKeys 生成密钥对
func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}
	return nil
}

// NewRsaWithFile 通过公钥、私钥本地文件存储路径构建 RSA
func NewRsaWithFile(publicKeyPath, privateKeyPath string) (*RSA, error) {
	var err error
	var publicKey, privateKey []byte
	if publicKey, err = os.ReadFile(publicKeyPath); err != nil {
		return nil, err
	}
	if privateKey, err = os.ReadFile(privateKeyPath); err != nil {
		return nil, err
	}
	return NewRsa(publicKey, privateKey)
}

// NewRsa 通过公钥、私钥构建 RSA
func NewRsa(publicKey []byte, privateKey []byte) (*RSA, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pri, ok := priv.(*rsa.PrivateKey)
	if ok {
		return &RSA{
			publicKey:  pub,
			privateKey: pri,
		}, nil
	} else {
		return nil, errors.New("private key not supported")
	}
}

// PublicEncrypt 公钥加密
func (r *RSA) PublicEncrypt(data []byte) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(encryptPKCS1v15)
	}
	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// PrivateDecrypt 私钥解密
func (r *RSA) PrivateDecrypt(encrypted string) ([]byte, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split(raw, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), err
}

// Sign 数据加签
func (r *RSA) Sign(data []byte) (string, error) {
	h := RsaAlgorithmSign.New()
	h.Write(data)
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, RsaAlgorithmSign, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// Verify 数据验签
func (r *RSA) Verify(data string, sign string) error {
	h := RsaAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.publicKey, RsaAlgorithmSign, hashed, decodedSign)
}

// MarshalPKCS8PrivateKey 封送PKCS8私钥
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
