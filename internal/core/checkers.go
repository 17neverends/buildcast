package core

import (
	"fmt"
	cfg "github.com/17neverends/buildcast/internal/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
)

func CheckServers(servers []cfg.Server) []cfg.Server {
	var available []cfg.Server
	for _, server := range servers {
		if err := checkSFTPConnection(server); err != nil {
			log.Printf("[ERROR] Server %s unreachable: %v", server.IP, err)
			continue
		}
		log.Printf("[OK] Server %s is reachable", server.IP)
		available = append(available, server)
	}
	return available
}

func checkSFTPConnection(server cfg.Server) error {
	config := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.IP, server.SFTPPort), config)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	return nil
}
