package counter

import (
	"github.com/cilium/ebpf/link"
	"net"
)

type Counter struct {
	objs counterObjects
	link link.Link
}

func NewCounter() (*Counter, error) {
	var objs counterObjects
	if err := loadCounterObjects(&objs, nil); err != nil {
		return nil, err
	}
	return &Counter{objs: objs}, nil
}

func (counter *Counter) Attach(ifName string) error {
	iface, err := net.InterfaceByName(ifName)
	if err != nil {
		return err
	}

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   counter.objs.CountPackets,
		Interface: iface.Index,
	})
	if err != nil {
		return err
	}
	counter.link = l
	return nil
}

func (counter *Counter) ReadCount() (uint64, error) {
	var count uint64
	err := counter.objs.PktCount.Lookup(uint32(0), &count)
	return count, err
}

func (counter *Counter) Close() error {
	if counter.link != nil {
		counter.link.Close()
	}
	return counter.objs.Close()
}
