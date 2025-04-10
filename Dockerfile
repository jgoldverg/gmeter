FROM ubuntu:24.10

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    clang llvm gcc make iproute2 iputils-ping \
    libelf-dev libbpf-dev \
    git curl ca-certificates \
    bpftool \
    linux-headers-$(uname -r) \
    golang

# Create workdir
WORKDIR /app

# Copy source
COPY . .

# Optional: run make or go build
# RUN make

CMD ["/bin/bash"]
