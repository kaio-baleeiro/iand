// installer/unzipper.go
package installer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip extrai e organiza os arquivos.
func Unzip(src, dest string) error {
	// 1. Extrai o zip para um diretório temporário
	tempDest := dest + "_temp_unzip"
	defer os.RemoveAll(tempDest) // Garante a limpeza do diretório temporário

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(tempDest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(tempDest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(tempDest)+string(os.PathSeparator)) {
			return fmt.Errorf("caminho de arquivo inválido: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		cerr1 := outFile.Close()
		cerr2 := rc.Close()
		if cerr1 != nil {
			fmt.Fprintf(os.Stderr, "erro ao fechar arquivo %s: %v\n", f.Name, cerr1)
		}
		if cerr2 != nil {
			fmt.Fprintf(os.Stderr, "erro ao fechar arquivo zip: %v\n", cerr2)
		}

		if err != nil {
			return err
		}
	}

	// 2. Encontra a pasta "cmdline-tools" extraída e a move para o destino final
	extractedCmdlineToolsDir := filepath.Join(tempDest, "cmdline-tools")
	finalPath := filepath.Join(dest, "cmdline-tools", "latest")

	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return fmt.Errorf("falha ao criar o diretório pai final: %w", err)
	}

	if err := os.Rename(extractedCmdlineToolsDir, finalPath); err != nil {
		return fmt.Errorf("falha ao mover 'cmdline-tools' para '%s': %w", finalPath, err)
	}

	return nil
}
