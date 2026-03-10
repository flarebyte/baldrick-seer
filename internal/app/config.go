package app

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrConfigPathDoesNotExist = errors.New("config path does not exist")
	ErrConfigPathIsDirectory  = errors.New("config path is a directory")
)

type ConfigRequest struct {
	Path string
}

type ConfigResult struct {
	Path string
}

func LoadConfig(req ConfigRequest) (ConfigResult, error) {
	info, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return ConfigResult{}, ErrConfigPathDoesNotExist
		}

		return ConfigResult{}, err
	}

	if info.IsDir() {
		return ConfigResult{}, ErrConfigPathIsDirectory
	}

	return ConfigResult{
		Path: filepath.Clean(req.Path),
	}, nil
}
