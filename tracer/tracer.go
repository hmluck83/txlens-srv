package tracer

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/hmluck83/txlens-srv/internal/stack"
	"github.com/holiman/uint256"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/debug"
	"github.com/lmittmann/w3/module/eth"
)

var EthAddress common.Address
var transferLog *uint256.Int
var depositLog *uint256.Int
var withDrawalLog *uint256.Int

const (
	erc20Transer = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	// de Facto
	deposit    = "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c"
	withdrawal = "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"
)

func init() {
	EthAddress = common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	transferLog = uint256.MustFromHex(erc20Transer)
	depositLog = uint256.MustFromHex(deposit)
	withDrawalLog = uint256.MustFromHex(withdrawal)
}

func transferValue(debugStructlog *debug.StructLog) *big.Int {
	// 상남자는 예외 처리 따위는 하지 않는다네
	// 언젠가 메모리가 터지겠지?
	offset := debugStructlog.Stack[len(debugStructlog.Stack)-1].Uint64()
	size := debugStructlog.Stack[len(debugStructlog.Stack)-2].Uint64()

	memorySlice := debugStructlog.Memory[offset : offset+size]
	result := new(big.Int).SetBytes(memorySlice)
	return result
}

func stackOffsetValue(debugStructLog *debug.StructLog, offset int) *uint256.Int {
	return &debugStructLog.Stack[len(debugStructLog.Stack)-offset]
}

func FundFlowFromTx(txHash common.Hash) (FundFlows, AddrLabels, error) {
	fundflows := FundFlows{}
	addrLabels := AddrLabels{}

	// panic if the PRC connection fail to establish
	var client = w3.MustDial(os.Getenv("RPCNODE"))

	var tx *types.Transaction
	err := client.Call(eth.Tx(txHash).Returns(&tx))

	if err != nil {
		panic(err)
	}

	// EOA로 부터 첫번째 트랜잭션의 0이 아닐 경우 기록
	if tx.Value().Cmp(big.NewInt(0)) != 0 {
		sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			panic(err)
		}
		ff := fundFlow{
			From:  sender,
			To:    *tx.To(),
			Value: *tx.Value(),
			Token: EthAddress,
		}
		fundflows = append(fundflows, ff)
		addressLabeler(&ff, addrLabels)
	}

	addressStack := stack.NewStack()
	addressStack.Push(tx.To())

	// 현재 이용하는 RPC-NODE가 개복치 같아, 자꾸 죽어버린다. 어쩔 수 없이
	// 마감을 앞두고 있는 상황에서 땜빵한다
	// 다른 node로 변경하고 살아 있음
	var debugTrace *debug.Trace

	err = client.Call(
		debug.TraceTx(txHash, &debug.TraceConfig{
			EnableStack:   true,
			EnableMemory:  true,
			EnableStorage: false,
		},
		).Returns(&debugTrace))

	if err != nil {
		panic(err)
	}

	for _, val := range debugTrace.StructLogs {
		switch val.Op {
		case vm.CALL:
			receiveAddress := val.Stack[len(val.Stack)-2]
			value := val.Stack[len(val.Stack)-3]

			if value.Cmp(uint256.NewInt(0)) > 0 {
				ff := fundFlow{
					From:  *addressStack.Peek(),
					To:    common.BytesToAddress(receiveAddress.Bytes()),
					Value: *value.ToBig(),
					Token: EthAddress,
				}
				fundflows = append(fundflows, ff)
				addressLabeler(&ff, addrLabels)
			}
			fallthrough

		case vm.CALLCODE, vm.DELEGATECALL, vm.STATICCALL:
			pushedAddress := common.BytesToAddress(val.Stack[len(val.Stack)-2].Bytes())
			addressStack.Push(&pushedAddress)

		case vm.RETURN, vm.REVERT, vm.STOP:
			addressStack.Pop()

		case vm.LOG3: // ERC20 Transfer
			topic1 := stackOffsetValue(val, 3)
			if topic1.Cmp(transferLog) == 0 {
				from := stackOffsetValue(val, 4)
				to := stackOffsetValue(val, 5)
				value := transferValue(val)

				ff := fundFlow{
					From:  from.Bytes20(),
					To:    to.Bytes20(),
					Value: *value,
					Token: *addressStack.Peek(),
				}
				fundflows = append(fundflows, ff)
				addressLabeler(&ff, addrLabels)
			}

		case vm.LOG2: // Non-Standard Transfer Log
			topic1 := stackOffsetValue(val, 3)

			if topic1.Cmp(depositLog) == 0 {
				from := *addressStack.Peek()
				to := stackOffsetValue(val, 4)
				value := transferValue(val)

				ff := fundFlow{
					From:  from,
					To:    to.Bytes20(),
					Value: *value,
					Token: from,
				}
				fundflows = append(fundflows, ff)
			}

			if topic1.Cmp(withDrawalLog) == 0 {
				from := stackOffsetValue(val, 4)
				to := *addressStack.Peek()
				value := transferValue(val)

				ff := fundFlow{
					From:  from.Bytes20(),
					To:    to,
					Value: *value,
					Token: to,
				}
				fundflows = append(fundflows, ff)
			}

		}
	}

	return fundflows, addrLabels, nil
}
