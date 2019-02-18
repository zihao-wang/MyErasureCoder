package main

import (
	ec "MyErasureCoder/ErasureCoder"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func NewFile(fileSize int, fileName string) (err error) {

	randSeq := make([]byte, fileSize)

	for i := 0; i < fileSize; i++ {
		randSeq[i] = byte(rand.Intn(16))
	}

	err = ioutil.WriteFile(fileName, randSeq, 0644)
	return
}

func SpeedTest(test_size int) {

	file_GB := float64(test_size) / 1024. / 1024. / 1024.

	e := ec.NewECoder(10, 4)
	test_file_name := "test_random_sequence"
	err := NewFile(test_size, test_file_name)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	e.LoadFile(test_file_name)
	elapse := time.Since(start).Seconds()
	fmt.Printf("load file time cost %v \n", elapse)

	start = time.Now()
	e.Encoding()
	elapse = time.Since(start).Seconds()
	fmt.Printf("encoding file time cost %v, speed %v GB/s \n", elapse, file_GB/elapse)

	start = time.Now()
	e.StoreAll("tmp")
	elapse = time.Since(start).Seconds()
	fmt.Printf("save file time cost %v \n", elapse)

	start = time.Now()
	e.DeleteData(1)
	e.FTReAssemble("recover")
	fmt.Printf("recover one block file time cost %v \n", elapse)
}

func main() {
	//sampleFileName := "Reed-Solomon-Error-Correction.pdf"
	//outputFolder := "tmp"
	//recoverFile := "recover.pdf"
	//fmt.Printf("test file %v\n", sampleFileName)
	//
	//os.Remove(recoverFile)
	//
	//e := ec.NewECoder(10, 4)
	//err := e.LoadFile(sampleFileName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//e.Encoding()
	//e.StoreAll(outputFolder)
	//e.DeleteData(3)
	//e.DeleteData(5)
	//e.DeleteData(7)
	//e.FTReAssemble(recoverFile)
	SpeedTest(1000000000)
}
