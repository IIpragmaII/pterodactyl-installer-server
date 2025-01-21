package main

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

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
