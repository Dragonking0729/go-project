package core

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func VmStdErrFormat(logs []vm.StructLog) {
	fmt.Fprintf(os.Stderr, "VM Stats %d ops\n", len(logs))
	for _, log := range logs {
		fmt.Fprintf(os.Stderr, "PC %-3d - %-14s\n", log.Pc, log.Op)
		fmt.Fprintln(os.Stderr, "STACK =", len(log.Stack))
		for i, item := range log.Stack {
			fmt.Fprintf(os.Stderr, "%04d: %x\n", i, common.LeftPadBytes(item.Bytes(), 32))
		}

		const maxMem = 10
		addr := 0
		fmt.Fprintln(os.Stderr, "MEM =", len(log.Memory))
		for i := 0; i+16 <= len(log.Memory) && addr < maxMem; i += 16 {
			data := log.Memory[i : i+16]
			str := fmt.Sprintf("%04d: % x  ", addr*16, data)
			for _, r := range data {
				if r == 0 {
					str += "."
				} else if utf8.ValidRune(rune(r)) {
					str += fmt.Sprintf("%s", string(r))
				} else {
					str += "?"
				}
			}
			addr++
			fmt.Fprintln(os.Stderr, str)
		}
	}
}
