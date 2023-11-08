package commands

import (
	"log"
	"math/rand"
	"os"
)

func GetGladosVoiceline() string {
	files, err := os.ReadDir("glados")
	if err != nil {
		log.Fatal(err)
	}

	return files[rand.Intn(len(files))].Name()
}
