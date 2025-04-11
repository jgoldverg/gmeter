//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -target bpfel,bpfeb counter ./counter.c

package counter
