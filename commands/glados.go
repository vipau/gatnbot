package commands 

import (
"log"
    "math/rand"
    "io/ioutil"
)

func GetGladosVoiceline() string {
    files, err := ioutil.ReadDir("glados")
    if err != nil {
        log.Fatal(err)
    }

    return files[rand.Intn(len(files))].Name()
}
