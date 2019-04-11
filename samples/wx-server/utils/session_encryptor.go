package utils

import (
	"time"
	"bytes"
	"encoding/binary"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
)

const _key = "^abcdefghijklmnopqrstuvwxyz1234509876ABCDEFGHIJKLMNOPQRSTUVWXYZ$"

func init() {
	rand.Seed(time.Now().Unix())
}

func generateKey(idx int) string {
	switch idx {
	case 0:
		return _key
	default:
		return fmt.Sprintf("%s%s", _key[idx:], _key[:idx])
	}
}

func CreateSession(openId string, id int64) string {
	buf := bytes.Buffer{}
	timestamp := time.Now().Unix()
	length := len(openId) + 8 + 8 // len(openId) + sizeof(id) + sizeof(timestamp)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint16(b, uint16(length))
	buf.Write(b[:2])

	buf.WriteString(openId)

	binary.LittleEndian.PutUint64(b, uint64(id))
	buf.Write(b)

	binary.LittleEndian.PutUint64(b, uint64(timestamp))
	buf.Write(b)

	idx := rand.Intn(64)
	encText := base64.NewEncoding(generateKey(idx)).EncodeToString(buf.Bytes())

	// prefix: idx + timestamp
	buf.Reset()
	binary.LittleEndian.PutUint16(b, uint16(idx))
	buf.Write(b[:2])

	binary.LittleEndian.PutUint64(b, uint64(timestamp))
	buf.Write(b)

	return fmt.Sprintf("%s%s", hex.EncodeToString(buf.Bytes()), encText)
}

func ParseSession(session string, timeout int) (openId string, id int64, err error) {
	if len(session) <= 20 {
		err = fmt.Errorf("invalid session")
		return
	}

	prefix, e := hex.DecodeString(session[:20])
	if e != nil {
		err = e
		return
	}
	idx := int(binary.LittleEndian.Uint16(prefix[:2]))
	if idx < 0 || idx >= 64 {
		err = fmt.Errorf("bad key")
		return
	}
	timestamp := int64(binary.LittleEndian.Uint64(prefix[2:]))

	oriSession, e := base64.NewEncoding(generateKey(idx)).DecodeString(session[20:])
	if e != nil {
		err = e
		return
	}
	if len(oriSession) <= 2 {
		err = fmt.Errorf("bad format")
		return
	}
	length := int(binary.LittleEndian.Uint16(oriSession[:2]))
	oriSession = oriSession[2:]
	if length != len(oriSession) {
		err = fmt.Errorf("bad length")
		return
	}

	openId = string(oriSession[:length-2*8])
	oriSession = oriSession[length-2*8:]
	id = int64(binary.LittleEndian.Uint64(oriSession[:8]))
	ts := int64(binary.LittleEndian.Uint64(oriSession[8:]))
	if timestamp != ts {
		err = fmt.Errorf("bad timestamp")
		return
	}
	if timestamp + int64(timeout) < time.Now().Unix() {
		err = fmt.Errorf("session timeout")
		return
	}
	return
}
