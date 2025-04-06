package main

import (
	"encoding/base64"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var configsLocation = os.Getenv("CONFIGS_LOCATION")
var scriptsLocation = os.Getenv("SCRIPTS_LOCATION")

var installPanel = [8]*step{
	{stepType: RUN, file: filepath.Join(scriptsLocation, "install_dependencies.sh")},
	{stepType: RUN, file: filepath.Join(scriptsLocation, "db_setup.sh")},
	{stepType: COPY, file: filepath.Join(configsLocation, "environment"), destination: "/var/www/pterodactyl/.env"},
	{stepType: RUN, file: filepath.Join(scriptsLocation, "env_configuration.sh")},
	{stepType: COPY, file: filepath.Join(configsLocation, "pteroq.service"), destination: "/etc/systemd/system/pteroq.service"},
	{stepType: RUN, file: filepath.Join(scriptsLocation, "queue_listener.sh")},
	{stepType: COPY, file: filepath.Join(configsLocation, "nginx.conf"), destination: "/etc/nginx/sites-available/pterodactyl.conf"},
	{stepType: RUN, file: filepath.Join(scriptsLocation, "nginx_setup.sh")},
}

var createNode = [1]*step{
	{stepType: RUN, file: filepath.Join(scriptsLocation, "create_node.sh")},
}

var installWings = [3]*step{
	{stepType: RUN, file: filepath.Join(scriptsLocation, "install_docker.sh")},
	{stepType: COPY, file: filepath.Join(configsLocation, "wings.service"), destination: "/etc/systemd/system/wings.service"},
	{stepType: RUN, file: filepath.Join(scriptsLocation, "install_wings.sh")},
}

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func runInstallation(c *gin.Context) {
	var settings settings
	err := c.ShouldBindJSON(&settings)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var auth goph.Auth
	if settings.Cert != "" {
		decodedCertBytes, _ := base64.StdEncoding.DecodeString(settings.Cert)
		decodedCert := string(decodedCertBytes)
		// Start new ssh connection with private key.
		auth, err = goph.RawKey(decodedCert, settings.Password)
	} else {
		auth = goph.Password(settings.Password)
	}
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
