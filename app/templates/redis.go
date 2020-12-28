package templates

import (
	"gitlab.com/scalablespace/listener/app/models"
	"text/template"
)

type RedisConf struct {
	RequirePass     string
	Secret          string
	MaxMemory       uint64
	MaxMemoryPolicy string
	IOThreads       int
}

var RedisConfTemplate = template.Must(template.New("redis.conf").Parse(`
protected-mode yes
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300
daemonize yes
supervised no
pidfile "/var/run/redis/redis-server.pid"
loglevel notice
logfile "/var/log/redis/redis-server.log"
databases 100
always-show-logo yes
save 900 1
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename "dump.rdb"
dir "/var/lib/redis"
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
io-threads {{ .IOThreads }}
`))

type Firewall struct {
	Instance *models.Instance
	Node     *models.Node
}

var UfwRedisTemplate = template.Must(template.New("redis").Parse(`*filter
:ufw-user-input - [0:0]
:ufw-user-output - [0:0]
:ufw-user-forward - [0:0]
:ufw-before-logging-input - [0:0]
:ufw-before-logging-output - [0:0]
:ufw-before-logging-forward - [0:0]
:ufw-user-logging-input - [0:0]
:ufw-user-logging-output - [0:0]
:ufw-user-logging-forward - [0:0]
:ufw-after-logging-input - [0:0]
:ufw-after-logging-output - [0:0]
:ufw-after-logging-forward - [0:0]
:ufw-logging-deny - [0:0]
:ufw-logging-allow - [0:0]
:ufw-user-limit - [0:0]
:ufw-user-limit-accept - [0:0]
### RULES ###

### tuple ### allow any 22 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 22 -j ACCEPT
-A ufw-user-input -p udp --dport 22 -j ACCEPT

### tuple ### allow any 6666 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 6666 -j ACCEPT
-A ufw-user-input -p udp --dport 6666 -j ACCEPT

{{if .Instance.Whitelist}}
{{range .Node.Whitelist}}
### tuple ### allow any 6379 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 6379 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 6379 -s {{.}} -j ACCEPT
### END RULES ###
{{end}}
{{range .Instance.Whitelist}}
### tuple ### allow any 6379 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 6379 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 6379 -s {{.}} -j ACCEPT
{{end}}
{{else}}
### tuple ### allow any 6379 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 6379 -j ACCEPT
-A ufw-user-input -p udp --dport 6379 -j ACCEPT
{{end}}

### LOGGING ###
-A ufw-after-logging-input -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-A ufw-after-logging-forward -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-I ufw-logging-deny -m conntrack --ctstate INVALID -j RETURN -m limit --limit 3/min --limit-burst 10
-A ufw-logging-deny -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-A ufw-logging-allow -j LOG --log-prefix "[UFW ALLOW] " -m limit --limit 3/min --limit-burst 10
### END LOGGING ###

### RATE LIMITING ###
-A ufw-user-limit -m limit --limit 3/minute -j LOG --log-prefix "[UFW LIMIT BLOCK] "
-A ufw-user-limit -j REJECT
-A ufw-user-limit-accept -j ACCEPT
### END RATE LIMITING ###
COMMIT
`))

const RedisAgentConf = `DB=redis
GOGC=40
REDIS_URL='redis://:%s@localhost:6379'
SECRET=%s
`
