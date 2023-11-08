package commands

import (
	"io/ioutil"
	"log"
	"math/rand"
)

func GetGladosVoiceline() string {
	files, err := ioutil.ReadDir("glados")
	if err != nil {
		log.Fatal(err)
	}

	return files[rand.Intn(len(files))].Name()
}
