package core

import (
	"fmt"
	cfg "github.com/17neverends/buildcast/internal/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RunCommand(initCmd string) error {
	parts := strings.Fields(initCmd)
	command := parts[0]
	args := parts[1:]

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeployToServer(server cfg.Server, remotePath, buildOutput string) error {
	if _, err := os.Stat("build"); os.IsNotExist(err) {
		return fmt.Errorf("local build directory does not exist")
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         60 * time.Second,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.IP, server.SFTPPort), sshConfig)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %v", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("SFTP client creation failed: %v", err)
	}
	defer client.Close()

	remoteBuildPath := filepath.ToSlash(filepath.Join(remotePath, "build"))

	log.Printf("Preparing to deploy from %s to %s", buildOutput, remoteBuildPath)

	localFiles, err := os.ReadDir(buildOutput)
	if err != nil {
		return fmt.Errorf("failed to read local build directory: %v", err)
	}
	if len(localFiles) == 0 {
		return fmt.Errorf("local build directory is empty")
	}

	if err := prepareRemoteDirectory(client, remoteBuildPath); err != nil {
		return fmt.Errorf("failed to prepare remote directory: %v", err)
	}

	if err := CopyFiles(client, buildOutput, remoteBuildPath); err != nil {
		return fmt.Errorf("failed to copy files: %v", err)
	}

	remoteFiles, err := client.ReadDir(remoteBuildPath)
	if err != nil {
		return fmt.Errorf("failed to verify remote files: %v", err)
	}

	log.Printf("Deployment successful. Copied %d files", len(remoteFiles))
	return nil
}
