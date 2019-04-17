/**
 * 微信信息加解密/签名/随机数
 * 1. HashStrings([]string)    -- 微信常用的签名算法 []string -> sort -> sha1 -> hex
 * 2. GetRandomBytes(int)      -- 获取指定长度的随机串，随机字符为 数字/小写字母/大写字母
 */
package wxmsg

import (
	"sort"
	"fmt"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"time"
	"bytes"
	"io"
	"github.com/rosbit/go-wx-api/conf"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func HashStrings(sl []string) string {
	sort.Strings(sl)
	h := sha1.New()
	for _, s := range sl {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func _PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func _PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func decryptMsg(body string, signature string, timestamp string, nonce string) ([]byte, error) {
	l := []string{wxconf.WxParams.Token, timestamp, nonce, body}
	hashcode := HashStrings(l)
	if hashcode != signature {
		return nil, fmt.Errorf("bad signature")
	}

	cryptedMsg, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}
	key := wxconf.WxParams.AesKey
	aesBlk, err := aes.NewCipher(key)
	blockSize := aesBlk.BlockSize()
	iv := key[:blockSize]
	blockMode := cipher.NewCBCDecrypter(aesBlk, iv)
	plainText := make([]byte, len(cryptedMsg))
	blockMode.CryptBlocks(plainText, cryptedMsg)
	plainText = _PKCS5UnPadding(plainText)

	xmlLen := binary.BigEndian.Uint32(plainText[16:20])
	return plainText[20:xmlLen+20], nil
}

var _rule = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
const _block_size = 32

func GetRandomBytes(n int) []byte {
	b := make([]byte, n)
	rc := len(_rule)
	for i:=0; i<n; i++ {
		b[i] = _rule[rand.Intn(rc)]
	}
	return b
}

func _msgToPad(msgLen int) []byte {
	paddingLen := _block_size - (msgLen % _block_size)
	if paddingLen == 0 {
		paddingLen = _block_size
	}
	padByte := byte(paddingLen)
	pad := make([]byte, paddingLen)
	for i:=0; i<paddingLen; i++ {
		pad[i] = padByte
	}
	return pad
}

func encryptMsg(msg []byte, timestamp, nonce string) (string, string) {
	randBytes := GetRandomBytes(16)
	msgLenInNet := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenInNet, uint32(len(msg)))

	origData := bytes.Buffer{}
	origData.Write(randBytes)
	origData.Write(msgLenInNet)
	origData.Write(msg)
	origData.WriteString(wxconf.WxParams.AppId)
	pad := _msgToPad(origData.Len())
	origData.Write(pad)

	key := wxconf.WxParams.AesKey
	block, _ := aes.NewCipher(key)
	iv := key[:block.BlockSize()]
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, origData.Len())
	blockMode.CryptBlocks(crypted, origData.Bytes())

	cryptedText := base64.StdEncoding.EncodeToString(crypted)
	return cryptedText, HashStrings([]string{wxconf.WxParams.Token, timestamp, nonce, cryptedText})
}
