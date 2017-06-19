package main

import (
	"github.com/FriggaHel/dockerhubbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
	"os"
)

func main() {
	err := beat.Run("dockerhubbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
