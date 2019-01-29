package main

import (
	ec "MyErasureCoder/ErasureCoder"
	"fmt"
	"log"
	"os"
)

func main() {
	//sampleFileName := "Reed-Solomon-Error-Correction.pdf"
	sampleFileName := "1707.07345.pdf"
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
	e.DeleteData(6)
	e.FTReAssemble(recoverFile)
}
