//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -tags linux counter ../ebpf/counter.c -- -I../headers

package gen
