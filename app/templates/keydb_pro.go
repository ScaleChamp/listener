package templates

import (
	"text/template"
)

type KeyDBProConf struct {
	RequirePass     string
	Secret          string
	MaxMemory       uint64
	MaxMemoryPolicy string
	ServerThreads   int
	EnablePro       string
}

var KeyDBProConfTemplate = template.Must(template.New("keydb-pro.conf").Parse(`
protected-mode yes
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300
daemonize yes
supervised no
pidfile "/var/run/keydb/keydb-server.pid"
loglevel notice
logfile "/var/log/keydb/keydb-server.log"
databases 100
always-show-logo yes
save 900 1
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename "dump.rdb"
dir "/var/lib/keydb"
replica-serve-stale-data yes
replica-read-only yes
repl-diskless-sync yes
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
replica-priority 100
lazyfree-lazy-eviction no
lazyfree-lazy-expire no
lazyfree-lazy-server-del no
replica-lazy-flush no
appendonly no
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
aof-use-rdb-preamble yes
lua-time-limit 5000
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
list-compress-depth 0
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
stream-node-max-bytes 4096
stream-node-max-entries 100
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
dynamic-hz yes
aof-rewrite-incremental-fsync yes
rdb-save-incremental-fsync yes

requirepass {{ .RequirePass }}
rename-command ACL {{ .Secret }}-ACL
rename-command SHUTDOWN {{ .Secret }}-SHUTDOWN
rename-command CONFIG {{ .Secret }}-CONFIG
rename-command SAVE {{ .Secret }}-SAVE
rename-command BGSAVE {{ .Secret }}-BGSAVE
rename-command BGREWRITEAOF {{ .Secret }}-BGREWRITEAOF
rename-command LATENCY {{ .Secret }}-LATENCY
rename-command CLIENT {{ .Secret }}-CLIENT
rename-command PSYNC {{ .Secret }}-PSYNC
rename-command MONITOR {{ .Secret }}-MONITOR
rename-command SLOWLOG {{ .Secret }}-SLOWLOG
rename-command SLAVEOF {{ .Secret }}-SLAVEOF
rename-command REPLICAOF {{ .Secret }}-REPLICAOF
maxmemory {{ .MaxMemory }}
maxmemory-policy {{ .MaxMemoryPolicy }}
server-threads {{ .ServerThreads }}
server-thread-affinity true
enable-pro {{ .EnablePro }}
storage-provider flash /flash
`))
