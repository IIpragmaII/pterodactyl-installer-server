package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

type stepType string

const COPY, RUN stepType = "copy", "run"

type step struct {
	stepType    stepType
	file        string
	destination string
}

type setting struct {
	value       string
	placeholder string
}

type settings struct {
	DbPassword *setting
	Email      *setting
	Timezone   *setting
	Username   *setting
	FirstName  *setting
	LastName   *setting
	Password   *setting
	Url        *setting
}

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

var steps = [2]*step{
	//{stepType: RUN, file: "scripts/install_dependencies.sh"},
	//{stepType: RUN, file: "scripts/db_setup.sh"},
	// {stepType: COPY, file: "configs/environment", destination: "/var/www/pterodactyl/.env"},
	// {stepType: RUN, file: "scripts/env_configuration.sh"},
	// {stepType: COPY, file: "configs/pteroq.service", destination: "/etc/systemd/system/pteroq.service"},
	// {stepType: RUN, file: "scripts/queue_listener.sh"},
	{stepType: COPY, file: "configs/nginx.conf", destination: "/etc/nginx/sites-available/pterodactyl.conf"},
	{stepType: RUN, file: "scripts/nginx_setup.sh"},
}

func getFileContent(file string, settings *settings) string {
	fileData, _ := os.ReadFile(file)
	fileContent := string(fileData)

	fields := reflect.ValueOf(*settings)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fileContent = strings.Replace(fileContent, field.Interface().(*setting).placeholder, field.Interface().(*setting).value, -1)
	}

	return fileContent
}

func uploadFile(client *goph.Client, content string, destination string) error {
	sftp, err := client.NewSftp()

	if err != nil {
		return err
	}

	file, err := sftp.Create(destination)

	if err != nil {
		return err
	}

	file.Write([]byte(content))
	file.Close()
	return nil
}

func main() {
	generateCert(serverIp)
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

	cert, key := generateCert(serverIp)
	err = uploadFile(client, cert, "/etc/ssl/pterodactyl-cert.pem")
	err = uploadFile(client, key, "/etc/ssl/pterodactyl-key.pem")

	settings := &settings{
		DbPassword: &setting{value: os.Getenv("DB_PASSWORD"), placeholder: "{{db_password}}"},
		Email:      &setting{value: os.Getenv("EMAIL"), placeholder: "{{email}}"},
		Timezone:   &setting{value: os.Getenv("TIMEZONE"), placeholder: "{{timezone}}"},
		Username:   &setting{value: os.Getenv("USERNAME"), placeholder: "{{username}}"},
		FirstName:  &setting{value: os.Getenv("FIRST_NAME"), placeholder: "{{first_name}}"},
		LastName:   &setting{value: os.Getenv("LAST_NAME"), placeholder: "{{last_name}}"},
		Password:   &setting{value: os.Getenv("PASSWORD"), placeholder: "{{password}}"},
		Url:        &setting{value: serverIp, placeholder: "{{url}}"},
	}

	for _, step := range steps {
		fileData := getFileContent(step.file, settings)
		if step.stepType == COPY {
			err = uploadFile(client, fileData, step.destination)
		} else {
			err = runCommand(client, fileData)
		}
	}
}
