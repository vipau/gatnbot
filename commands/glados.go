package commands

import (
	"log"
	"math/rand"
	"os"
)

func GetGladosVoiceline() string {
	files, err := os.ReadDir(GladosDir)
	if err != nil {
		log.Fatal(err)
	}

	return files[rand.Intn(len(files))].Name()
}
