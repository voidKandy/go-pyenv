package pyenv

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

// Darwin Executor
func (env *DarwinPyEnv) AddDependencies(requirementsPath string) error {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/bin/pip")
	err := dependencyHelper(env.EnvOptions, fp, requirementsPath)
	if err != nil {
		return err
	}
	return nil
}

func (env *DarwinPyEnv) ExecutePython(args ...string) (*exec.Cmd, error) {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/bin/python")
	cmd, err := executeHelper(env.EnvOptions, fp, args...)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// Linux Executor
func (env *LinuxPyEnv) AddDependencies(requirementsPath string) error {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/bin/pip")
	err := dependencyHelper(env.EnvOptions, fp, requirementsPath)
	if err != nil {
		return err
	}
	return nil
}

func (env *LinuxPyEnv) ExecutePython(args ...string) (*exec.Cmd, error) {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/bin/python")
	cmd, err := executeHelper(env.EnvOptions, fp, args...)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// Windows Executor
func (env *WindowsPyEnv) AddDependencies(requirementsPath string) error {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/Scripts/pip3.exe")
	err := dependencyHelper(env.EnvOptions, fp, requirementsPath)
	if err != nil {
		return err
	}
	return nil
}

func (env *WindowsPyEnv) ExecutePython(args ...string) (*exec.Cmd, error) {
	fp := filepath.Join(env.EnvOptions.ParentPath, "dist/python/install/python.exe")
	cmd, err := executeHelper(env.EnvOptions, fp, args...)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// helper functions

func dependencyHelper(env *PyEnvOptions, fp string, requirementsPath string) error {
	if env.Compressed {
		if err := env.DecompressDist(); err != nil {
			return err
		}
	}
	log.Println("installing python dependencies")
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error installing python dependencies: %v", err)
	}
	log.Println("installing python dependencies complete")
	return nil
}

func executeHelper(env *PyEnvOptions, fp string, args ...string) (*exec.Cmd, error) {
	if env.Compressed {
		return nil, fmt.Errorf("cannot execute python with a compressed dist")
	}
	cmd := exec.Command(fp, args...)
	return cmd, nil
}
