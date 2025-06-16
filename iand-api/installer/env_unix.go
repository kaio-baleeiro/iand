//go:build !windows

// installer/env_unix.go
package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func setupEnvironmentVariables(androidHome, executablesPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	profiles := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".zprofile"),
	}

	lines := []string{
		"# IAND ANDROID CONFIG - Adicionado automaticamente",
		fmt.Sprintf("export IAND_ANDROID_HOME=\"%s\"", androidHome),
		"export ANDROID_HOME=\"$IAND_ANDROID_HOME\"",
		fmt.Sprintf("export IAND_EXECUTABLES_PATH=\"%s\"", executablesPath),
		"export PATH=\"$IAND_EXECUTABLES_PATH:$PATH\"",
	}

	var profileFound bool
	for _, profilePath := range profiles {
		data, err := os.ReadFile(profilePath)
		content := string(data)
		var toAdd []string
		for _, line := range lines {
			if !strings.Contains(content, line) {
				toAdd = append(toAdd, line)
			}
		}
		if len(toAdd) > 0 {
			if err := appendToFile(profilePath, "\n"+strings.Join(toAdd, "\n")+"\n"); err != nil {
				return fmt.Errorf("falha ao escrever em %s: %w", profilePath, err)
			}
			profileFound = true
		} else if err == nil {
			profileFound = true // Já está tudo configurado
		}
	}

	if !profileFound {
		profilePath := filepath.Join(homeDir, ".bash_profile")
		if err := appendToFile(profilePath, "\n"+strings.Join(lines, "\n")+"\n"); err != nil {
			return fmt.Errorf("falha ao criar e escrever em %s: %w", profilePath, err)
		}
	}
	return nil
}

// appendToFile adiciona um conteúdo a um arquivo.
func appendToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
