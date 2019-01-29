package mathop

import (
	"errors"
	"fmt"
	"log"
)

const SanityCheck = false

type GF8Matrix struct {
	numRow int
	numCol int
	Matrix [][]byte
}

func NewMatrix(numRow, numCol int) *GF8Matrix {
	var m GF8Matrix
	m.numRow = numRow
	m.numCol = numCol
	for rIter := 0; rIter < numRow; rIter++ {
		m.Matrix = append(m.Matrix, make([]byte, numCol))
	}
	return &m
}

func NewMatrixFromRows(rows [][]byte) *GF8Matrix {
	m := NewMatrix(len(rows), len(rows[0]))
	m.Matrix = rows
	return m
}

func NewEncodingMatrix(numData, numParity int) *GF8Matrix {
	m := NewMatrix(numData+numParity, numData)
	for r := 0; r < numData; r++ {
		m.Matrix[r][r] = 1
	}
	for r := 0; r < numParity; r++ {
		for c := 0; c < numData; c++ {
			m.Matrix[r+numData][c] = inverseTbl[r+1^c+1]
		}
	}
	return m
}

func (m *GF8Matrix) Mul(n *GF8Matrix) *GF8Matrix {
	if m.numCol != n.numRow {
		log.Fatalf("GF8Matrix shape not compatible (%v, %v) != (%v, %v)", m.numRow, m.numRow, n.numRow, n.numCol)
	}
	outMat := NewMatrix(m.numRow, n.numCol)
	for lrIter := 0; lrIter < m.numRow; lrIter++ {
		for rcIter := 0; rcIter < n.numCol; rcIter++ {
			var out byte = 0
			for ii := 0; ii < m.numCol; ii++ {
				out ^= mulTbl[m.Matrix[lrIter][ii]][n.Matrix[ii][rcIter]]
			}
			outMat.Matrix[lrIter][rcIter] = out
		}
	}
	return outMat
}

// swapRow ...
// locally swap the two rows in a matrix
func (m *GF8Matrix) swapRow(i, j int) {
	for c := 0; c < m.numCol; c++ {
		m.Matrix[i][c], m.Matrix[j][c] = m.Matrix[j][c], m.Matrix[i][c]
	}
}

// multRow ...
// locally multiply two rows in a matrix
func (m *GF8Matrix) multRow(i int, coef byte) {
	for c := 0; c < m.numCol; c++ {
		m.Matrix[i][c] = mulTbl[m.Matrix[i][c]][coef]
	}
}

// mulAddToRow
// multiply one row and add to another
func (m *GF8Matrix) mulAddToRow(from int, coef byte, to int) {
	for c := 0; c < m.numCol; c++ {
		m.Matrix[to][c] ^= mulTbl[m.Matrix[from][c]][coef]
	}

}

var ErrNotSquareMatrix = errors.New("Matrix have no inverse for shape issue")
var ErrSingularMatrix = errors.New("Square Matrix is Singular")

// Inv() ...
// return the inverse of the matrix by Gaussian-Jordan Elimination
// not use the parrent pointer, make new instead, so that the modifications of GJE do not affect the original one
func (mPointer *GF8Matrix) Inv() (*GF8Matrix, error) {
	m := NewMatrix(mPointer.numRow, mPointer.numCol)
	copy(m.Matrix, mPointer.Matrix)
	if m.numCol != m.numRow {
		log.Fatal(ErrNotSquareMatrix)
	}
	// extend identity matrix to right half
	for r := 0; r < m.numRow; r++ {
		prepandArray := make([]byte, m.numCol)
		prepandArray[r] = 1
		m.Matrix[r] = append(m.Matrix[r], prepandArray...)
	}
	originalNumCol := m.numCol
	m.numCol += originalNumCol
	// make left half identity
	for r := 0; r < m.numRow; r++ {
		// find head
		singularFlag := true
		for iFindHead := r; iFindHead < m.numRow; iFindHead++ {
			if m.Matrix[iFindHead][r] != 0 {
				if iFindHead != r {
					m.swapRow(iFindHead, r)
				}
				singularFlag = false
				break
			}
		}
		if singularFlag {
			return nil, ErrSingularMatrix
		}
		// make this diagnal unit
		head := m.Matrix[r][r]
		m.multRow(r, inverseTbl[head])
		// erase other rows
		for rFollow := 0; rFollow < m.numRow; rFollow++ {
			if rFollow == r {
				continue
			}
			m.mulAddToRow(r, m.Matrix[rFollow][r], rFollow)
		}
	}
	//remove left half
	for r := 0; r < m.numRow; r++ {
		m.Matrix[r] = m.Matrix[r][originalNumCol:]
	}
	m.numCol = originalNumCol

	if SanityCheck {
		shouldBeIdentity := m.Mul(mPointer)
		shouldBeIdentity.ShowMatrix(SanityCheck)
		for r := 0; r < m.numRow; r++ {
			for c := 0; c < m.numCol; c++ {
				if r == c && shouldBeIdentity.Matrix[r][c] != 1 {
					log.Fatalf("Fail the sanity check, %v != 1", shouldBeIdentity.Matrix[r][c])
				} else if r != c && shouldBeIdentity.Matrix[r][c] != 0 {
					log.Fatalf("Fail the sanity check, %v != 0", shouldBeIdentity.Matrix[r][c])
				}
			}
		}
	}

	return m, nil
}

func (m *GF8Matrix) ShowMatrix(debug bool) {
	if !debug {
		return
	}
	for r := 0; r < m.numRow; r++ {
		for c := 0; c < m.numCol; c++ {
			fmt.Printf("%d\t", m.Matrix[r][c])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (m *GF8Matrix) GetRow(rId int) *GF8Matrix {
	outMat := NewMatrix(1, m.numCol)
	outMat.Matrix[0] = m.Matrix[rId]
	return outMat
}

func (m *GF8Matrix) GetShape() (int, int) {
	return m.numRow, m.numCol
}
