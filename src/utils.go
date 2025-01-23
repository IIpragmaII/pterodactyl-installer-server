package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/melbahja/goph"
)

func runCommand(client *goph.Client, cmd string) error {
	out, err := client.Run(cmd)

	fmt.Println(string(out))

	if err != nil {
		return err
	}

	return nil
}

func getFileContent(file string, settings *settings) string {
	fileData, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

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
