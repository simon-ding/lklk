package main

import (
	"github.com/simon-ding/lklk/server"
	"github.com/sirupsen/logrus"
)

func main() {
	s := server.New()

	err := s.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}