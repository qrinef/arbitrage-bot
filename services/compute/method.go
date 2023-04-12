package compute

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

var (
	emptyBytes = string(make([]byte, 31))

	type1Indexes = Indexes{
		path:      [2]int{164, 196},
		countPath: [2]int{132, 163},
		deadline:  [2]int{100, 132},
	}
	type2Indexes = Indexes{
		path:      [2]int{196, 228},
		countPath: [2]int{164, 195},
		deadline:  [2]int{132, 164},
	}

	type1Methods = [][]byte{
		common.Hex2Bytes("7ff36ab5"), // swapExactETHForTokens
	}
	type2Methods = [][]byte{
		common.Hex2Bytes("38ed1739"), // swapExactTokensForTokens
		common.Hex2Bytes("5c11d795"), // swapExactTokensForTokensSupportingFeeOnTransferTokens
		common.Hex2Bytes("18cbafe5"), // swapExactTokensForETH
		common.Hex2Bytes("791ac947"), // swapExactTokensForETHSupportingFeeOnTransferTokens
	}

	type1JoinedMethods = append(bytes.Join(type1Methods, []byte{0}), byte(0))
	type2JoinedMethods = append(bytes.Join(type2Methods, []byte{0}), byte(0))
)

type Method struct {
	AmountIn *uint256.Int
	Path     []common.Address
	Deadline *uint256.Int
}

type Indexes struct {
	path      [2]int
	countPath [2]int
	deadline  [2]int
}

func (s *Service) decodeMethod(data []byte, value *uint256.Int) (result *Method, err error) {
	method, indexes, err := s.createMethod(data, value)
	if err != nil {
		return result, err
	}

	if string(data[indexes.countPath[0]:indexes.countPath[1]]) != emptyBytes {
		return result, errors.New("tx path incorrect")
	}

	countPath := int(data[indexes.countPath[1]])
	if countPath < 2 || countPath > 10 {
		return result, errors.Errorf("tx path len %d, max. %d paths", countPath, 10)
	}

	for i := 0; i < countPath*32; i += 32 {
		to := indexes.path[1] + i

		if len(data) < to {
			return result, errors.Errorf("tx data bytes %d < %d request", len(data), to)
		}

		method.Path = append(method.Path, common.BytesToAddress(data[indexes.path[0]+i:to]))
	}

	method.Deadline = new(uint256.Int).SetBytes(data[indexes.deadline[0]:indexes.deadline[1]])

	return &method, err
}

func (s *Service) createMethod(data []byte, value *uint256.Int) (method Method, indexes Indexes, err error) {
	if len(data) < 196 {
		return method, indexes, errors.New("tx data size incorrect")
	}

	methodSignature := append(data[:4], byte(0))

	if bytes.Contains(type1JoinedMethods, methodSignature) {
		indexes = type1Indexes
		method.AmountIn = value
	} else if bytes.Contains(type2JoinedMethods, methodSignature) {
		indexes = type2Indexes
		method.AmountIn = new(uint256.Int).SetBytes(data[4:36])
	} else {
		return method, indexes, errors.New("tx method not found")
	}

	return method, indexes, err
}
