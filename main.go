//go:generate go fmt ./...

package main

import "gitlab.com/scalablespace/listener/config"

// iptables-save | tee iptables_backup.conf | grep -v '\-A' | iptables-restore
/*
sysctl -w net.ipv4.ip_forward=1

TARGET=1.1.1.1 PORT=6379 /usr/local/bin/iptables-atomic.sh /usr/local/bin/failover-proxy.sh
*/

/*
/usr/local/bin/failover-proxy.sh

#!/bin/bash
iptables -t nat -A PREROUTING -p tcp --dport $PORT -j DNAT --to-destination $TARGET:$PORT
iptables -t nat -A POSTROUTING -j MASQUERADE
*/

/*
/usr/local/bin/iptables-atomic.sh

#!/bin/sh

ip netns add staging
ip netns exec staging $@
ip netns exec staging iptables-save >/tmp/iptables_dump
iptables-restore /tmp/iptables_dump
ip netns exec staging ip6tables-save >/tmp/ip6tables_dump
ip6tables-restore /tmp/ip6tables_dump
ip netns delete staging
*/

func main() {
	config.NewApp().Run()
}
