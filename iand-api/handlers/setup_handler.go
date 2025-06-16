// handlers/setup_handler.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"install-android/installer"
)

func SetupAndroidHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Dispara a instalação em background para não bloquear a resposta
	go installer.Run()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Processo de instalação do Android CLI iniciado em background.",
		"details": "Verifique os logs do servidor para acompanhar o progresso.",
	})
	log.Println("Requisição recebida. Iniciando a instalação do Android CLI...")
}
