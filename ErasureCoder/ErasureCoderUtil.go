package ErasureCoder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// LoadFile ...
// load one file and transfer the information into structure
func (e *ECoder) LoadFile(fileName string) error {
	inputFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = e.Split(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	DPrint("num of dataBlocks %d", len(e.dataBlocks))
	return err
}

func (e *ECoder) StoreAll(saveToFolder string) {
	DPrint("Store All ...")
	for iData := 0; iData < e.numData; iData++ {
		fileName := fmt.Sprintf("%v/data-%v.blk", saveToFolder, iData)
		e.dataPath[iData] = fileName
		DPrint("DataBlock   # %d,   %s, len=%d", iData, fileName, len(e.dataBlocks[iData]))
		err := ioutil.WriteFile(fileName, e.dataBlocks[iData], 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	for iParity := 0; iParity < e.numParity; iParity++ {
		fileName := fmt.Sprintf("%v/parity-%v.blk", saveToFolder, iParity)
		e.parityPath[iParity] = fileName
		DPrint("ParityBlock # %d, %s, len=%d", iParity, fileName, len(e.parityBlocks[iParity]))
		err := ioutil.WriteFile(fileName, e.parityBlocks[iParity], 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	DPrint("All Stored")
}

func (e *ECoder) DeleteData(delNum int) error {
	DPrint("delete data block # %d", delNum)
	err := os.Remove(e.dataPath[delNum])
	return err
}

func (e *ECoder) FTReAssemble(outName string) {
	e.CheckAll()
	e.ReAssemble(outName)
}

func (e *ECoder) ReAssemble(outName string) error {
	var LoadedDataBlocks [][]byte
	var originBlock []byte
	for i := 0; i < len(e.dataPath); i++ {
		LoadedDataBlocks = append(LoadedDataBlocks, nil)
	}
	for k, v := range e.dataPath {
		inputFile, err := ioutil.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		LoadedDataBlocks[k] = inputFile
	}
	for i := 0; i < len(LoadedDataBlocks); i++ {
		originBlock = append(originBlock, LoadedDataBlocks[i]...)
	}
	err := ioutil.WriteFile(outName, originBlock, 0644)
	return err
}
