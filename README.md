# Go Meter

A tool for monitoring networking traffic using Cillium eBPF

# Build the image
docker build -t ebpf-dev .

# Run the container with direct mounts and privileges
docker run --rm -it \
--privileged \
--network host \
-v /sys:/sys \
-v /lib/modules:/lib/modules \
-v /proc:/proc \
-v /path/to/your/project:/home/devuser/app \
-v /sys/kernel/btf:/sys/kernel/btf \
gmeter

