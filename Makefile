.PHONY: all build clean

BPF_SRC := pkg/bpf/programs/monitor.c
BPF_OUT := build/bpf/monitor.o
BINARY := bin/gmeter

all: build

build: $(BPF_OUT)
	go build -o $(BINARY) .

$(BPF_OUT): $(BPF_SRC)
	@mkdir -p $(dir $(BPF_OUT))
	clang -O2 -g -Wall -target bpf -c $(BPF_SRC) -o $(BPF_OUT)

clean:
	rm -rf build bin
