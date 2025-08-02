package util

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func SaveAvatarCampaign(file io.Reader, filename string) (string, error) {
	savePath := filepath.Join("images", filename)
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		log.Printf("Folder 'images' tidak ada. Membuat folder...")
		err = os.MkdirAll("images", os.ModePerm)
		if err != nil {
			log.Printf("Gagal membuat folder: %v", err)
			return "", err
		}
	}

	out, err := os.Create(savePath)
	if err != nil {
		log.Printf("Gagal membuat file: %v", err)
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Gagal menyalin file: %v", err)
		return "", err
	}

	log.Printf("File berhasil disimpan di: %s", savePath)
	return savePath, nil
}
