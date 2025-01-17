package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var certPassword = os.Getenv("CERT_PASSWORD")
var serverIp = os.Getenv("SERVER_IP")
var certPath = os.Getenv("CERT_PATH")

// Disable host verification for now.
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func runCommand(client *goph.Client, cmd string) error {
	out, err := client.Run(cmd)

	fmt.Println(string(out))

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Start new ssh connection with private key.
	auth, err := goph.Key(certPath, certPassword)
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User: "root", Addr: serverIp, Port: 22, Auth: auth, Callback: VerifyHost,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing the network connection.
	defer client.Close()

	fileData, _ := os.ReadFile("scripts/install_dependencies.sh")

	runCommand(client, string(fileData))
}
