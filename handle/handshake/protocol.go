package handshake

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/xuyang-lee/wormhole/config/env"
	"io"
	"strings"
)

var ErrVersionMismatch = errors.New("version mismatch")

func Hello(rw io.ReadWriter) error {
	helloBuf := new(bytes.Buffer)
	if err := binary.Write(helloBuf, binary.BigEndian, env.ProtocolHelloSize); err != nil {
		return err
	}
	helloBuf.Write([]byte(env.ProtocolHelloMsg))

	// 1. 发送 HELLO
	if _, err := rw.Write(helloBuf.Bytes()); err != nil {
		return fmt.Errorf("handshake write failed: %w", err)
	}

	// 2. 等待 hi
	var hiSize uint16
	if err := binary.Read(rw, binary.BigEndian, &hiSize); err != nil {
		return err
	}

	// 读取 hi msg
	buf := make([]byte, hiSize)
	n, err := io.ReadFull(rw, buf)
	if err != nil {
		return fmt.Errorf("handshake read failed: %w", err)
	}

	args := strings.Split(string(buf[:n]), env.ProtocolSplitMark)
	if len(args) < 2 {
		return fmt.Errorf("handshake invalid message: %q", buf[:n])
	}
	serverName := args[0]
	version := args[1]

	if serverName != env.ServerName {
		return fmt.Errorf("handshake invalid server: %q", serverName)
	}
	if version != env.Version {
		return ErrVersionMismatch
	}

	return nil
}

func Hi(rw io.ReadWriter) error {

	var helloSize uint16
	if err := binary.Read(rw, binary.BigEndian, &helloSize); err != nil {
		return err
	}

	// 1. 读取 hello msg
	buf := make([]byte, helloSize)
	n, err := io.ReadFull(rw, buf)
	if err != nil {
		return fmt.Errorf("handshake read failed: %w", err)
	}

	// 拆分协议信息
	args := strings.Split(string(buf[:n]), env.ProtocolSplitMark)
	if len(args) < 2 {
		return fmt.Errorf("handshake invalid message: %q", buf[:n])
	}
	serverName := args[0]
	//version := args[1]

	if serverName != env.ServerName {
		return fmt.Errorf("handshake invalid server: %q", serverName)
	}

	// 2. 回复 OK
	hiBuf := new(bytes.Buffer)
	if err = binary.Write(hiBuf, binary.BigEndian, env.ProtocolHiSize); err != nil {
		return err
	}
	hiBuf.Write([]byte(env.ProtocolHelloMsg))
	if _, err = rw.Write(hiBuf.Bytes()); err != nil {
		return fmt.Errorf("handshake write failed: %w", err)
	}

	return nil
}
