package pyenv

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PyEnvOptions struct {
	ParentPath string
	// distributions: darwin/amd64 darwin/arm64 linux/arm64 linux/amd64 windows/386 windows/amd64
	Distribution string
	Compressed   bool
}

type Installer interface {
	Install() error
}

type Executor interface {
	AddDependencies(string) error
	ExecutePython(...string) (*exec.Cmd, error)
}

type PyEnv struct {
	EnvOptions PyEnvOptions
	Installer
	Executor
}

type (
	DarwinPyEnv  struct{ EnvOptions *PyEnvOptions }
	LinuxPyEnv   struct{ EnvOptions *PyEnvOptions }
	WindowsPyEnv struct{ EnvOptions *PyEnvOptions }
)

func NewPyEnv(path string, dist string) (*PyEnv, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting $HOME directory: %v", err)
	}
	if path == homedir {
		err := fmt.Errorf("path cannot be homedir\npath given: %s\nhomedir: %s", path, homedir)
		return nil, err
	}
	pyEnv := PyEnv{
		EnvOptions: PyEnvOptions{
			ParentPath:   path,
			Distribution: dist,
			Compressed:   false,
		},
	}
	osArchError := fmt.Errorf("this os/arch distribution is not supported: %v", dist)
	switch {
	case strings.Contains(dist, "darwin"):
		if dist != "darwin/amd64" && dist != "darwin/arm64" {
			return nil, osArchError
		}
		pyEnv.Installer = &DarwinPyEnv{&pyEnv.EnvOptions}
		pyEnv.Executor = &DarwinPyEnv{&pyEnv.EnvOptions}
	case strings.Contains(dist, "linux"):
		if dist != "linux/amd64" && dist != "linux/arm64" {
			return nil, osArchError
		}
		pyEnv.Installer = &LinuxPyEnv{&pyEnv.EnvOptions}
		pyEnv.Executor = &LinuxPyEnv{&pyEnv.EnvOptions}
	case strings.Contains(dist, "windows"):
		if dist != "windows/amd64" && dist != "windows/386" {
			return nil, osArchError
		}
		pyEnv.Installer = &WindowsPyEnv{&pyEnv.EnvOptions}
		pyEnv.Executor = &WindowsPyEnv{&pyEnv.EnvOptions}
	default:
		return nil, osArchError
	}

	return &pyEnv, nil
}

func (env *PyEnvOptions) DistExists() (*bool, error) {
	t := true
	f := false
	_, err := os.Stat(DistDirPath(env))
	if err == nil {
		return &t, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Stat(DistZipPath(env))
		if err == nil {
			return &t, nil
		}
		if errors.Is(err, os.ErrNotExist) {
			return &f, nil
		}
		return nil, err

	}
	return nil, err
}

func (env *PyEnvOptions) CompressDist() error {
	if env.Compressed {
		return fmt.Errorf("dist is already compressed")
	}

	if err := compressDir(DistDirPath(env), DistZipPath(env)); err != nil {
		return fmt.Errorf("error compressing python environment: %v", err)
	}
	env.Compressed = true

	if err := os.RemoveAll(DistDirPath(env)); err != nil {
		return fmt.Errorf("error removing old uncompressed evironment: %v", err)
	}
	log.Printf("removed %v\n", DistDirPath(env))
	return nil
}

func (env *PyEnvOptions) DecompressDist() error {
	if !env.Compressed {
		log.Println("dist is already decompressed")
		return nil
	}

	env.Compressed = false

	if err := unzipSource(DistZipPath(env), DistDirPath(env)); err != nil {
		return fmt.Errorf("error unzipping compressed evironment: %v", err)
	}
	if err := os.RemoveAll(DistZipPath(env)); err != nil {
		return fmt.Errorf("error removing old compressed evironment: %v", err)
	}
	log.Printf("removed %v\n", DistZipPath(env))
	return nil
}

func DistDirPath(env *PyEnvOptions) string {
	return filepath.Join(env.ParentPath, "dist")
}

func DistZipPath(env *PyEnvOptions) string {
	return DistDirPath(env) + ZIP_FILE_EXT
}
