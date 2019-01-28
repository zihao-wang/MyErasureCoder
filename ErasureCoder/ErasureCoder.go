package ErasureCoder

import (
	"fmt"
	"io/ioutil"
	"log"
)

var debug = true

func DPrint(format string, a ... interface{}) {
	if debug {
		s := fmt.Sprintf(format, a)
		fmt.Printf("%s\n", s)
	}
}

type ECoder struct {
	NumData int
	dataBlocks []dataBuff
	NumParity int
	parityBlocks []parityBuff
}


type dataBuff struct {

}

type parityBuff struct {

}

func (e ECoder) Encoding(data) {

}

func (e ECoder) Split(data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, ErrShortData
	}
	// Calculate number of bytes per data shard.
	perShard := (len(data) + e.NumData - 1) / e.NumData

	if cap(data) > len(data) {
		data = data[:cap(data)]
	}
	// Only allocate memory if necessary
	if len(data) < (e.NumData * perShard) {
		// Pad data to r.Shards*perShard.
		padding := make([]byte, (e.NumData*perShard)-len(data))
		data = append(data, padding...)
	}

	// Split into equal-length shards.
	dst := make([][]byte, e.NumData)
	for i := 0; i < e.NumData; i++  {
		DPrint("slicing the %v-th shard", i)
		dst[i] = data[:perShard]
		data = data[perShard:]
	}

	return dst, nil
}


func (e ECoder) LoadFile(fileName string) ([][]byte, error) {
	inputFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	split, err := e.Split(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	return split, err
}
