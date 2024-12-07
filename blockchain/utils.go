package blockchain

import (
	"io"
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

func ConnectToIPFS() {
	sh = shell.NewShell("localhost:5001")
	if !sh.IsUp() {
		log.Fatal("IPFS daemon is not running!")
	}
}

func AddFileToIPFS(data string) (string, error) {
	cid, err := sh.Add(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return cid, nil
}

func GetFileFromIPFS(cid string) (string, error) {
	content, err := sh.Cat(cid)
	if err != nil {
		return "", err
	}
	result, err := io.ReadAll(content)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
