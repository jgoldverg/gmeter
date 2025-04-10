// SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define IPPROTO_IPV4 0x0800

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, __u32);     // source IP
    __type(value, __u64);   // packet count
    __uint(max_entries, 1024);
} packet_counter SEC(".maps");

SEC("xdp")
int xdp_counter_prog(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    struct ethhdr *eth = data;
    if ((void*)(eth + 1) > data_end)
        return XDP_PASS;

    // check if it's an IPv4 packet
    if (bpf_ntohs(eth->h_proto) != IPPROTO_IPV4)
        return XDP_PASS;

    struct iphdr *iph = data + sizeof(*eth);
    if ((void*)(iph + 1) > data_end)
        return XDP_PASS;

    __u32 src_ip = iph->saddr;

    // update the packet count map
    __u64 *value = bpf_map_lookup_elem(&packet_counter, &src_ip);
    if (value) {
        __sync_fetch_and_add(value, 1);
    } else {
        __u64 init_val = 1;
        bpf_map_update_elem(&packet_counter, &src_ip, &init_val, BPF_ANY);
    }

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
