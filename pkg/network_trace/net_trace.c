#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

typedef struct {
    __u32 saddr;
    __u32 daddr;
    __u16 sport;
    __u16 dport;
    __u8 proto;
} conn_key_t;

typedef struct {
    __u64 packets;
    __u64 bytes;
} conn_value_t;

struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH),
    __uint(max_entries, MAX_ENTRIES),
    __type(key, conn_key_t),
    __type(value, conn_value_t)
}

#define MAX_ENTRIES 1024
#define TASK_COMM_LEN 16

// Configurable target (set from userspace)
volatile const char TARGET_COMM[TASK_COMM_LEN] = "java"; // e.g., Minecraft
volatile const pid_t TARGET_PID = 0;                     // 0 = match by name


SEC("kprobe/tcp_sendmsg")
int BPT_KPROBE(tcp_sendmsg, struct sock *sk, struct msg_hdr *msg, size_t size) {
    return track_traffic(sk, size);
}

SEC("kprobe/tcp_cleanup_rbuf")
int BPF_KPROBE(tcp_cleanup_rbuf, struct sock *sk, int copied) {
    if(copied < 0) return 0;
    return track_traffic(sk, copied);
}


static int track_traffic(struct sock *sk, size_t bytes) {
        pid_t pid = bpf_get_current_pid_tgid() >> 32;
        char comm[TASK_COMM_LEN];
        bpf_get_current_comm(&comm, sizeof(comm));
        if (TARGET_PID != 0 && pid != TARGET_PID) return 0;
        if (TARGET_PID == 0 && __builtin_memcmp(comm, TARGET_COMM, sizeof(TARGET_COMM)) != 0) {
            return 0;
        }

        conn_key_t key = {0};

}