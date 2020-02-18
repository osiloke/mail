package main

import (
	"fmt"

	"github.com/jaffee/commandeer"
	"github.com/osiloke/mail/worker"
)

var BUILD string = ""

func main() {
	err := commandeer.Run(worker.NewWorker(BUILD))
	if err != nil {
		fmt.Println(err)
	}
}
