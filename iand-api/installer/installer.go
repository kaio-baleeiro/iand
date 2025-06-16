// installer/installer.go
package installer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// Defina um SDK recente e estável como padrão
	sdkVersion        = "34"
	buildToolsVersion = "34.0.0"
)

var (
	// Links de download fornecidos
	urls = map[string]string{
		"windows": "https://dl.google.com/android/repository/commandlinetools-win-13114758_latest.zip",
		"darwin":  "https://dl.google.com/android/repository/commandlinetools-mac-13114758_latest.zip",
		"linux":   "https://dl.google.com/android/repository/commandlinetools-linux-13114758_latest.zip",
	}
)

// Run orquestra todo o processo de instalação.
func Run() error {
	log.Println("--- Início do Setup do Android CLI ---")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível obter o diretório home do usuário: %w", err)
	}

	// Definição dos caminhos base
	iandDir := filepath.Join(homeDir, ".iand")
	softwaresDir := filepath.Join(iandDir, "softwares")
	androidCliDir := filepath.Join(softwaresDir, "android-cli")
	finalCmdlineToolsDir := filepath.Join(androidCliDir, "cmdline-tools", "latest")

	// 1. Criar diretórios necessários
	log.Printf("1/6. Criando diretório base em %s", androidCliDir)
	if err := os.MkdirAll(androidCliDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretórios: %w", err)
	}

	// 2. Baixar o arquivo ZIP correto para o SO
	osName := runtime.GOOS // "windows", "darwin", "linux"
	url, ok := urls[osName]
	if !ok {
		return fmt.Errorf("sistema operacional '%s' não suportado", osName)
	}

	zipPath := filepath.Join(softwaresDir, "android-cmdline-tools.zip")
	log.Printf("2/6. Baixando ferramentas para %s de %s", osName, url)
	if err := downloadFile(zipPath, url); err != nil {
		return fmt.Errorf("falha no download: %w", err)
	}
	defer os.Remove(zipPath) // Limpa o zip após o uso

	// 3. Extrair o conteúdo para o local correto
	log.Printf("3/6. Extraindo arquivo para %s", finalCmdlineToolsDir)
	if err := Unzip(zipPath, androidCliDir); err != nil {
		return fmt.Errorf("falha ao extrair o arquivo: %w", err)
	}

	// 4. Configurar Variáveis de Ambiente
	log.Println("4/6. Configurando variáveis de ambiente...")
	cmdlineToolsBin := filepath.Join(finalCmdlineToolsDir, "bin")
	platformToolsPath := filepath.Join(androidCliDir, "platform-tools")
	emulatorPath := filepath.Join(androidCliDir, "emulator")

	executablesPath := strings.Join([]string{cmdlineToolsBin, platformToolsPath, emulatorPath}, string(os.PathListSeparator))

	// Monta o novo PATH incluindo todos os binários relevantes
	newPath := os.Getenv("PATH") + string(os.PathListSeparator) + executablesPath

	if err := setupEnvironmentVariables(androidCliDir, executablesPath); err != nil {
		return fmt.Errorf("falha ao configurar variáveis de ambiente: %w", err)
	}
	log.Println("-> Variáveis de ambiente IAND_ANDROID_HOME, ANDROID_HOME e PATH configuradas.")
	// Força o carregamento das novas variáveis no processo atual
	os.Setenv("ANDROID_HOME", androidCliDir)
	os.Setenv("PATH", newPath)

	// 5. Aceitar licenças do SDK
	sdkManagerPath := filepath.Join(finalCmdlineToolsDir, "bin", "sdkmanager")
	if runtime.GOOS == "windows" {
		sdkManagerPath += ".bat"
	}

	log.Println("5/6. Aceitando licenças do SDK automaticamente...")
	licenseInput := strings.NewReader(strings.Repeat("y\n", 50))
	cmd := exec.Command(sdkManagerPath, "--licenses")
	cmd.Stdin = licenseInput
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("AVISO: Ocorreu um problema ao tentar aceitar as licenças, mas o processo continuará. Erro: %v", err)
	} else {
		log.Println("-> Licenças aceitas com sucesso.")
	}

	// 6. Instalar pacotes essenciais do SDK
	packagesToInstall := []string{
		"platform-tools",
		"emulator",
		fmt.Sprintf("platforms;android-%s", sdkVersion),
		fmt.Sprintf("build-tools;%s", buildToolsVersion),
		"system-images;android-34;google_apis;x86_64",
	}
	log.Printf("6/6. Instalando pacotes do SDK: %s", strings.Join(packagesToInstall, ", "))

	installCmd := exec.Command(sdkManagerPath, packagesToInstall...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("falha ao instalar pacotes do SDK: %w", err)
	}

	log.Println("-> Pacotes do SDK instalados com sucesso.")
	log.Println("--- ✅ Setup do Android CLI Concluído! ---")
	log.Println("AVISO IMPORTANTE: É necessário reiniciar o terminal (ou fazer logout/login) para que as variáveis de ambiente (PATH) sejam carregadas.")
	return nil
}

// downloadFile baixa um arquivo de uma URL para um caminho local.
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Fecha o body da resposta HTTP

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if cerr != nil {
			log.Printf("erro ao fechar arquivo %s: %v", filepath, cerr)
		}
	}() // Fecha o arquivo criado, com log de erro se necessário

	_, err = io.Copy(out, resp.Body)
	return err
}
