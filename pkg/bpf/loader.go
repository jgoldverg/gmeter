package bpf

import (
	"fmt"
	"github.com/cilium/ebpf"
)

func LoadBPFProgram(path string) (*ebpf.Collection, error) {
	coll, err := ebpf.LoadCollection(path)
	if err != nil {
		return nil, fmt.Errorf("error loading bpf program: %s", err)
	}
	return coll, nil
}
