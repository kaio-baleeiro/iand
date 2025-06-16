// logger/logger.go
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Init inicializa a configuração global de log para a aplicação.
// A saída será direcionada para o console (stdout) e para um arquivo
// em ~/.iand/logs/DD-MM-YYYY.log.txt
func Init() {
	// 1. Obter o diretório home do usuário de forma segura.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("ERRO CRÍTICO: Não foi possível obter o diretório home: %v", err)
	}

	// 2. Construir o caminho para o diretório de logs e criá-lo se não existir.
	logDir := filepath.Join(homeDir, ".iand", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("ERRO CRÍTICO: Não foi possível criar o diretório de log em %s: %v", logDir, err)
	}

	// 3. Formatar a data atual para o nome do arquivo (DD-MM-YYYY).
	// O layout "02-01-2006" é a forma padrão do Go para especificar esse formato.
	dateStr := time.Now().Format("02012006")
	fileName := fmt.Sprintf("%s.log.txt", dateStr)
	filePath := filepath.Join(logDir, fileName)

	// 4. Abrir o arquivo de log para adicionar conteúdo.
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("ERRO CRÍTICO: Não foi possível abrir o arquivo de log %s: %v", filePath, err)
	}

	// 5. Criar um MultiWriter para duplicar a saída de log para o terminal e para o arquivo.
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 6. Configurar o logger padrão do Go para usar nosso MultiWriter.
	log.SetOutput(multiWriter)

	// (Opcional) Adicionar prefixos úteis aos logs, como data, hora e arquivo de origem.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Log de confirmação que usará a nova configuração.
	log.Println("Logger inicializado com sucesso.")
	log.Printf("Saída de log configurada para o terminal e para o arquivo: %s", filePath)
}
