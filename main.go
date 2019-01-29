package main

import (
	ec "MyErasureCoder/ErasureCoder"
	"fmt"
	"log"
	"os"
)

func main() {
	sampleFileName := "Reed-Solomon-Error-Correction.pdf"
	outputFolder := "tmp"
	recoverFile := "recover.pdf"
	fmt.Printf("test file %v\n", sampleFileName)

	os.Remove(recoverFile)

	e := ec.NewECoder(10, 4)
	err := e.LoadFile(sampleFileName)
	if err != nil {
		log.Fatal(err)
	}
	e.Encoding()
	e.StoreAll(outputFolder)
	for k := 0; k < 10; k++ {
		e.DeleteData(1)
		e.FTReAssemble(recoverFile)
	}
}
