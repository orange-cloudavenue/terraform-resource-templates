package file

import (
	"io"
	"os"
)

// IsFileExists checks if a file exists and is not a directory. It returns false if the file does not exist or is a directory.
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ToString reads a file and returns its content as a string.
func ToString(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ReadFile reads a file and returns its content as a byte array.
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
