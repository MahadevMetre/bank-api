package ftp_server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPClient struct {
	Host     string
	User     string
	Password string
	Port     int
	conn     *ftp.ServerConn
}

func NewFTPClient(host string, port int, user, password string) *FTPClient {
	return &FTPClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (f *FTPClient) Connect() error {
	addr := fmt.Sprintf("%s:%d", f.Host, f.Port)
	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to FTP server: %v", err)
	}

	err = conn.Login(f.User, f.Password)
	if err != nil {
		conn.Quit()
		return fmt.Errorf("failed to login: %v", err)
	}

	f.conn = conn
	return nil
}

func (f *FTPClient) Disconnect() error {
	if f.conn != nil {
		return f.conn.Quit()
	}
	return nil
}

// UploadFile uploads a file to the FTP server
func (f *FTPClient) UploadFile(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	remoteDir := filepath.Dir(remotePath)
	if remoteDir != "." {
		err = f.createRemoteDir(remoteDir)
		if err != nil {
			return fmt.Errorf("failed to create remote directory: %v", err)
		}
	}

	err = f.conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from the FTP server
func (f *FTPClient) DownloadFile(localPath, remotePath string) error {
	resp, err := f.conn.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer resp.Close()

	// Create local directory if it doesn't exist
	localDir := filepath.Dir(localPath)
	if localDir != "." {
		err = os.MkdirAll(localDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create local directory: %v", err)
		}
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// ListFiles lists files in the specified directory
func (f *FTPClient) ListFiles(path string) ([]string, error) {
	entries, err := f.conn.List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name)
	}

	return files, nil
}

// createRemoteDir creates a directory and all necessary parent directories
func (f *FTPClient) createRemoteDir(path string) error {
	dirs := splitPath(path)
	currentPath := ""

	for _, dir := range dirs {
		currentPath = filepath.Join(currentPath, dir)
		err := f.conn.MakeDir(currentPath)
		if err != nil {
			// Ignore error if directory already exists
			if !isExistsError(err) {
				return err
			}
		}
	}

	return nil
}

// splitPath splits a path into its components
func splitPath(path string) []string {
	var parts []string
	dir := path

	for dir != "." && dir != "/" {
		parts = append([]string{filepath.Base(dir)}, parts...)
		dir = filepath.Dir(dir)
	}

	return parts
}

// isExistsError checks if the error is due to the directory already existing
func isExistsError(err error) bool {
	return err.Error() == "550 File exists" || err.Error() == "550 Directory exists"
}
