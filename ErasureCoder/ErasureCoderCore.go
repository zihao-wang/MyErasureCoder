package ErasureCoder

import (
	"MyErasureCoder/ErasureCoder/mathop"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var debug = true

func DPrint(format string, a ...interface{}) {
	if debug {
		format = format + "\n"
		fmt.Printf(format, a...)
	}
}

type ECoder struct {
	shardSize    int
	numData      int
	dataBlocks   [][]byte
	dataPath     map[int]string
	numParity    int
	parityBlocks [][]byte
	parityPath   map[int]string
	matrix       *mathop.GF8Matrix
}

func NewECoder(numData int, numParity int) *ECoder {
	ec := ECoder{numData: numData, numParity: numParity}
	for i := 0; i < numData; i++ {
		ec.dataBlocks = append(ec.dataBlocks, nil)
	}
	ec.dataPath = make(map[int]string)
	for i := 0; i < numParity; i++ {
		ec.parityBlocks = append(ec.parityBlocks, nil)
	}
	ec.parityPath = make(map[int]string)
	ec.matrix = mathop.NewEncodingMatrix(numData, numParity)
	return &ec
}

// Encoding ...
// generates the parity blocks from the input data block
func (e *ECoder) Encoding() error {
	if len(e.dataBlocks) != e.numData {
		return ErrNotMatchNumDataBlocks
	}
	var wg sync.WaitGroup
	for iParity := 0; iParity < e.numParity; iParity++ {
		wg.Add(1)
		go e.encode(iParity, &wg)
	}
	wg.Wait()
	return nil
}

func (e *ECoder) encode(whichParity int, group *sync.WaitGroup) {
	encodingRow := e.matrix.GetRow(e.numData + whichParity)
	dataMatrix := mathop.NewMatrixFromRows(e.dataBlocks)
	parityRow := encodingRow.Mul(dataMatrix)
	DPrint("parity block # %d encoded, length %d", whichParity, len(parityRow.Matrix[0]))
	e.parityBlocks[whichParity] = parityRow.Matrix[0]
	group.Done()
}

// Split ...
// split a datastream of byte array
// note that a byte is 8 bit, so the basic encoding is based on GF(2^8)
// control parameters are preserved by ECoder
func (e *ECoder) Split(data []byte) error {
	if len(data) == 0 {
		return ErrShortData
	}
	// Calculate number of bytes per data shard.
	perShard := (len(data) + e.numData - 1) / e.numData
	e.shardSize = perShard
	if cap(data) > len(data) {
		data = data[:cap(data)]
	}
	// Only allocate memory if necessary
	if len(data) < (e.numData * perShard) {
		// Pad data to r.Shards*perShard.
		padding := make([]byte, (e.numData*perShard)-len(data))
		data = append(data, padding...)
	}
	// Split into equal-length shards.
	for i := 0; i < e.numData; i++ {
		e.dataBlocks[i] = data[:perShard]
		data = data[perShard:]
	}
	return nil
}

func (e *ECoder) ClearCache() {
	for i := 0; i < e.numData; i++ {
		e.dataBlocks[i] = nil
	}
	for i := 0; i < e.numParity; i++ {
		e.parityBlocks[i] = nil
	}
}

// CheckAll ...
// this function check whether all the data and parity blocks are there, if not do reconstruction & regeneration
func (e *ECoder) CheckAll() error {
	e.matrix.ShowMatrix()
	var healthDataKey []int
	var healthParityKey []int
	healthDataPath := make(map[int]string)
	brokenDataPath := make(map[int]string)
	healthParityPath := make(map[int]string)
	brokenParityPath := make(map[int]string)
	for k := 0; k < len(e.dataPath); k++ {
		v := e.dataPath[k]
		if _, err := os.Stat(v); err != nil {
			DPrint("data file %v err: [%v]", v, err)
			brokenDataPath[k] = v
		} else {
			healthDataPath[k] = v
			healthDataKey = append(healthDataKey, k)
		}
	}

	for k := 0; k < len(e.parityPath); k++ {
		v := e.parityPath[k]
		if _, err := os.Stat(v); err != nil {
			DPrint("parity file %v err: [%v]", v, err)
			brokenParityPath[k] = v
		} else {
			healthParityPath[k] = v
			healthParityKey = append(healthParityKey, k)
		}
	}

	DPrint("Survive data blocks %v, survive parity blocks %v", healthDataKey, healthParityKey)

	// recover
	if len(brokenDataPath) > 0 {
		if len(healthDataPath)+len(healthParityPath) >= e.numData {
			var matRows [][]byte
			var dataRows [][]byte
			numRow := 0
			for i := 0; i < len(healthDataKey); i++ {
				k := healthDataKey[i]
				matRows = append(matRows, e.matrix.GetRow(k).Matrix[0])
				dataRows = append(dataRows, e.dataBlocks[k])
				numRow++
			}
			for i := 0; i < len(healthDataKey); i++ {
				k := healthParityKey[i]
				matRows = append(matRows, e.matrix.GetRow(k + e.numData).Matrix[0])
				dataRows = append(dataRows, e.parityBlocks[k])
				numRow++
				if numRow >= e.numData {
					break
				}
			}
			subEncodeMat := mathop.NewMatrixFromRows(matRows)
			subDataNParityMat := mathop.NewMatrixFromRows(matRows)
			subEncodeMat.ShowMatrix()
			decodeMat, err := subEncodeMat.Inv()
			decodeMat.ShowMatrix()
			if err != nil {
				log.Fatal(err)
			}

			wg := sync.WaitGroup{}
			for k, v := range brokenDataPath {
				wg.Add(1)
				go func(k int, v string, wg *sync.WaitGroup) {
					subDecodeRow := decodeMat.GetRow(k)
					recovered := subDecodeRow.Mul(subDataNParityMat).Matrix[0]
					err := ioutil.WriteFile(v, recovered, 0644)
					DPrint("recover block %v at %v", k, v)
					if err != nil {
						log.Fatal(err)
					}
					wg.Done()
				}(k, v, &wg)
			}
			wg.Wait()

		} else {
			return ErrNoSufficientBlocks
		}
	}
	return nil
}
