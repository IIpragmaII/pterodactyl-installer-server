package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melbahja/goph"
)

type settings struct {
	DbPassword   string `json:"dbPassword" placeholder:"{{db_password}}"`
	Email        string `json:"email" placeholder:"{{email}}"`
	Timezone     string `json:"timezone" placeholder:"{{timezone}}"`
	Username     string `json:"username" placeholder:"{{username}}"`
	FirstName    string `json:"firstName" placeholder:"{{first_name}}"`
	LastName     string `json:"lastName" placeholder:"{{last_name}}"`
	Password     string `json:"password" placeholder:"{{password}}"`
	Cert         string `json:"cert"`
	CertPassword string `json:"certPassword"`
	ServerIp     string `json:"serverIp" placeholder:"{{url}}"`
}

func runInstallation(c *gin.Context) {
	var settings settings
	err := c.ShouldBindJSON(&settings)

	fmt.Print(err)

	// Start new ssh connection with private key.
	auth, err := goph.RawKey(settings.Cert, settings.CertPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	client, err := goph.NewConn(&goph.Config{
		User: "root", Addr: settings.ServerIp, Port: 22, Auth: auth, Callback: VerifyHost,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Defer closing the network connection.
	defer client.Close()

	cert, key := generateCert(settings.ServerIp)
	err = uploadFile(client, cert, "/etc/ssl/pterodactyl-cert.pem")
	err = uploadFile(client, key, "/etc/ssl/pterodactyl-key.pem")

	runInstallSteps(client, installPanel[:], &settings)
	runInstallSteps(client, createNode[:], &settings)
	runInstallSteps(client, installWings[:], &settings)

}
