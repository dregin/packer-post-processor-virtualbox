// https://gist.github.com/jedy/3357393
// https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
package main

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

const privateKey = `content of id_rsa`

type keychain struct {
	key *rsa.PrivateKey
}

func (k *keychain) Key(i int) (interface{}, error) {
	if i != 0 {
		return nil, nil
	}
	return &k.key.PublicKey, nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	hashFunc := crypto.SHA1
	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)
	return rsa.SignPKCS1v15(rand, k.key, hashFunc, digest)
}

func main() {
	block, _ := pem.Decode([]byte(privateKey))
	rsakey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	clientKey := &keychain{rsakey}
	clientConfig := &ssh.ClientConfig{
		User: "wuhao",
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(clientKey),
		},
	}
	client, err := ssh.Dial("tcp", "127.0.0.1:22", clientConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		content := "123456789\n"
		fmt.Fprintln(w, "C0644", len(content), "testfile")
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00") // 传输以\x00结束
	}()
	if err := session.Run("/usr/bin/scp -qrt ./"); err != nil {
		panic("Failed to run: " + err.Error())
	}
}
