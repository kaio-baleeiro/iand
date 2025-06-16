// installer/env_windows.go
//go:build windows

package installer

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// setupEnvironmentVariables configura variáveis de ambiente permanentes no Windows.
func setupEnvironmentVariables(androidHome, executablesPath string) error {
	log.Println("-> Usando lógica de ambiente para Windows com comandos nativos (reg, setx).")

	// Usamos 'setx' para variáveis simples, pois é direto e eficaz.
	if err := setx("IAND_ANDROID_HOME", androidHome); err != nil {
		return fmt.Errorf("falha ao definir IAND_ANDROID_HOME: %w", err)
	}
	if err := setx("ANDROID_HOME", "%IAND_ANDROID_HOME%"); err != nil {
		return fmt.Errorf("falha ao definir ANDROID_HOME: %w", err)
	}
	if err := setx("IAND_EXECUTABLES_PATH", executablesPath); err != nil {
		return fmt.Errorf("falha ao definir IAND_EXECUTABLES_PATH: %w", err)
	}

	// Para o PATH, usamos 'reg' para evitar o bug de 1024 caracteres do 'setx'.
	return addToUserPathWithReg("%IAND_EXECUTABLES_PATH%")
}

// setx é um wrapper para o comando `setx` do Windows.
func setx(key, value string) error {
	cmd := exec.Command("setx", key, value)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Evita piscar uma janela de console.
	// Redireciona a saída para o log para depuração.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// addToUserPathWithReg adiciona um novo valor ao PATH do usuário usando o comando `reg.exe`.
func addToUserPathWithReg(newPath string) error {
	// 1. Consulta o valor atual do PATH no registro.
	queryCmd := exec.Command("reg", "query", `HKCU\Environment`, "/v", "Path")
	var out bytes.Buffer
	queryCmd.Stdout = &out
	err := queryCmd.Run()

	// Se a chave 'Path' não existir, o comando retorna um erro. Nós o tratamos como um path vazio.
	currentPath := ""
	if err == nil {
		// O output do 'reg query' é verboso. Precisamos extrair apenas o valor.
		// Ex: HKEY_CURRENT_USER\Environment\n    Path    REG_SZ    C:\path1;C:\path2
		output := out.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		lastLine := lines[len(lines)-1]
		parts := strings.Fields(lastLine)
		if len(parts) >= 3 {
			currentPath = strings.Join(parts[2:], " ")
		}
	} else {
		log.Println("-> Variável Path não encontrada no registro do usuário, será criada uma nova.")
	}

	// 2. Verifica se o caminho já existe para evitar duplicação.
	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		// Normaliza a comparação (e.g. %VAR% vs C:\path\to\var)
		if os.ExpandEnv(p) == os.ExpandEnv(newPath) {
			log.Println("-> Caminho já existe no PATH do usuário. Nenhuma alteração necessária.")
			return nil
		}
	}

	// 3. Adiciona o novo caminho e o define no registro.
	var newFullPath string
	if currentPath == "" {
		newFullPath = newPath
	} else {
		newFullPath = currentPath + ";" + newPath
	}

	log.Printf("-> Adicionando '%s' ao PATH do usuário.", newPath)
	// Usa /t REG_EXPAND_SZ para permitir variáveis como %SystemRoot% no PATH.
	addCmd := exec.Command("reg", "add", `HKCU\Environment`, "/v", "Path", "/t", "REG_EXPAND_SZ", "/d", newFullPath, "/f")
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr

	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("falha ao adicionar o caminho ao registro com 'reg add': %w", err)
	}

	return nil
}
