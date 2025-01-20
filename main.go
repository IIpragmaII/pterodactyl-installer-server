package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
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

var installPanel = [8]*step{
	{stepType: RUN, file: "scripts/install_dependencies.sh"},
	{stepType: RUN, file: "scripts/db_setup.sh"},
	{stepType: COPY, file: "configs/environment", destination: "/var/www/pterodactyl/.env"},
	{stepType: RUN, file: "scripts/env_configuration.sh"},
	{stepType: COPY, file: "configs/pteroq.service", destination: "/etc/systemd/system/pteroq.service"},
	{stepType: RUN, file: "scripts/queue_listener.sh"},
	{stepType: COPY, file: "configs/nginx.conf", destination: "/etc/nginx/sites-available/pterodactyl.conf"},
	{stepType: RUN, file: "scripts/nginx_setup.sh"},
}

var createNode = [1]*step{
	{stepType: RUN, file: "scripts/create_node.sh"},
}

var installWings = [3]*step{
	{stepType: RUN, file: "scripts/install_docker.sh"},
	{stepType: COPY, file: "configs/wings.service", destination: "/etc/systemd/system/wings.service"},
	{stepType: RUN, file: "scripts/install_wings.sh"},
}

func getFileContent(file string, settings *settings) string {
	fileData, _ := os.ReadFile(file)
	fileContent := string(fileData)

	fields := reflect.ValueOf(*settings)

	for i := 0; i < fields.NumField(); i++ {
		placeholder := fields.Type().Field(i).Tag.Get("placeholder")
		if placeholder != "" {
			fileContent = strings.Replace(fileContent, placeholder, fields.Field(i).String(), -1)
		}
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

func runInstallSteps(client *goph.Client, steps []*step, settings *settings) {
	for _, step := range steps {
		fileData := getFileContent(step.file, settings)
		if step.stepType == COPY {
			err := uploadFile(client, fileData, step.destination)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := runCommand(client, fileData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {

	app := gin.Default()
	app.POST("/install", runInstallation)

	if address := os.Getenv("ADDRESS"); address != "" {
		port := os.Getenv("PORT")
		app.Run(address + ":" + port)
	}
	app.Run()
}
