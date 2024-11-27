package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/hirochachacha/go-smb2"
	"github.com/spf13/viper"
)

// SMBClient defines the structure for SMB operations
type SMBClient struct {
	Connection net.Conn
	Session    *smb2.Session
}

// SMBConfig holds the configuration for SMB connection
type SMBConfig struct {
	Address  string // e.g., "192.168.1.10"
	Port     int    // default is 445
	Username string
	Password string
	Domain   string // optional, use "" for most cases
	Share    string // e.g., "SharedFolder"
}

// NewSMBClient creates and returns a new SMBClient
func NewSMBClient(config SMBConfig) (*SMBClient, error) {

	// Establish a connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Address, config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMB server: %w", err)
	}
	defer conn.Close()

	// Establish SMB dialer
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.Username,
			Password: config.Password,
			Domain:   config.Domain,
		},
	}
	// Connect to the shared folder
	session, err := dialer.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer session.Logoff()

	return &SMBClient{
		Connection: conn,
		Session:    session,
	}, nil
}

// Close cleans up the SMB session and share
func (c *SMBClient) Close() {
	if c.Session != nil {
		c.Session.Logoff()
	}
	if c.Connection != nil {
		c.Connection.Close()
	}
}

func LoadSMBFromViper(env string) (*SMBClient, error) {
	smbEnv := viper.GetStringMapString("SMB." + env)
	config := SMBConfig{
		Address:  viper.GetString(fmt.Sprintf("SMB.%s.address", smbEnv)),
		Port:     viper.GetInt(fmt.Sprintf("SMB.%s.port", smbEnv)),
		Username: viper.GetString(fmt.Sprintf("SMB.%s.username", smbEnv)),
		Password: viper.GetString(fmt.Sprintf("SMB.%s.password", smbEnv)),
	}
	smbClient, err := NewSMBClient(config)
	if err != nil {
		return nil, err
	}
	return smbClient, nil
}

// UploadFile uploads a file to the SMB share
func (c *SMBClient) UploadFile(localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	dstFile, err := fs.Create(remoteFileName)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from the SMB share
func (c *SMBClient) DownloadFile(remotePath, localPath string) error {
	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	srcFile, err := fs.Open(remoteFileName)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer srcFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}

func (c *SMBClient) WriteFile(remotePath, content string) error {
	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	dstFile, err := fs.Create(remoteFileName)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer dstFile.Close()

	_, err = dstFile.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to remote file: %w", err)
	}

	return nil
}

func (c *SMBClient) WriteBinaryFile(remotePath string, content []byte) error {
	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	dstFile, err := fs.Create(remoteFileName)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer dstFile.Close()

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to remote file: %w", err)
	}

	return nil
}

func (c *SMBClient) ReadFile(remotePath string) (string, error) {
	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return "", fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	srcFile, err := fs.Open(remoteFileName)
	if err != nil {
		return "", fmt.Errorf("failed to open remote file: %w", err)
	}
	defer srcFile.Close()

	content, err := io.ReadAll(srcFile)
	if err != nil {
		return "", fmt.Errorf("failed to read remote file: %w", err)
	}

	return string(content), nil
}

func (c *SMBClient) ReadBinaryFile(remotePath string) ([]byte, error) {
	remoteFileName := filepath.Base(remotePath)
	remoteDir := filepath.Dir(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	srcFile, err := fs.Open(remoteFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open remote file: %w", err)
	}
	defer srcFile.Close()

	content, err := io.ReadAll(srcFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote file: %w", err)
	}

	return content, nil
}

// DeleteFile deletes a file from the SMB share
func (c *SMBClient) DeleteFile(remotePath string) error {
	remoteDir := filepath.Dir(remotePath)
	remoteFileName := filepath.Base(remotePath)

	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	err = fs.Remove(remoteFileName)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// ListFiles lists files in a remote directory
func (c *SMBClient) ListFiles(remoteDir string) ([]string, error) {
	fs, err := c.Session.Mount(remoteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to mount share: %w", err)
	}
	defer fs.Umount()

	dir, err := fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, file := range dir {
		files = append(files, file.Name())
	}

	return files, nil
}
