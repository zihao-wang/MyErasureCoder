package main

import (
	"fmt"
	ec "MyErasureCoder/ErasureCoder"
)


func main() {
	sampleFileName := "Reed-Solomon-Error-Correction.pdf"
	fmt.Printf("test file %v\n", sampleFileName)
	e := ec.ECoder{10, 4}
	byteArrays, _ := e.LoadFile(sampleFileName)
	for i := range byteArrays {
		fmt.Println(byteArrays[i])
	}
}
