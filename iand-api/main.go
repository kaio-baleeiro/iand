// main.go
package main

import (
	"log"
	"net/http"

	"install-android/handlers"
	"install-android/logger" // <-- Importa o novo pacote
)

func main() {
	// Inicializa o logger como a primeira ação da aplicação.
	// Todo log a partir daqui será gerenciado pelo nosso pacote.
	logger.Init()

	// O resto do seu código continua normalmente...
	mux := http.NewServeMux()
	mux.HandleFunc("/setup-android-cli", handlers.SetupAndroidHandler)

	log.Println("🚀 Servidor da install-android iniciado em http://localhost:8080")
	log.Println("✅ Endpoint pronto para receber requisições: POST /setup-android-cli")
	log.Println("ℹ️  Para iniciar, execute: curl -X POST http://localhost:8080/setup-android-cli")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		// Este log também será capturado pelo nosso logger.
		log.Fatalf("❌ Erro fatal ao iniciar o servidor: %v", err)
	}
}
