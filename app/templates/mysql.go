package templates

import (
	"text/template"
)

var MySQLAgentConf = `DB=mysql
GOGC=40
MYSQL_URL=mysql://ssnuser:%s@localhost:3306/ssndb
SECRET=%s
`

var UfwMySQLTemplate = template.Must(template.New("mysql").Parse(`*filter
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
### tuple ### allow any 3306 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 3306 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 3306 -s {{.}} -j ACCEPT
### END RULES ###
{{end}}
{{range .Instance.Whitelist}}
### tuple ### allow any 3306 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 3306 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 3306 -s {{.}} -j ACCEPT
{{end}}
{{else}}
### tuple ### allow any 3306 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 3306 -j ACCEPT
-A ufw-user-input -p udp --dport 3306 -j ACCEPT
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

var BackupPushMySQLConfTemplate = template.Must(template.New("backup-push.sh").Parse(`#!/bin/bash
export WALG_MYSQL_DATASOURCE_NAME='root@unix(/var/run/mysqld/mysqld.sock)/mysql'
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export WALG_PGP_KEY_PATH=/etc/public.pgp
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"

export AWS_S3_FORCE_PATH_STYLE="true"
export AWS_ACCESS_KEY_ID={{ .AwsKey }}
export AWS_SECRET_ACCESS_KEY={{ .AwsSecretKey }}
export AWS_REGION={{ .Region }}

export WALG_STREAM_CREATE_COMMAND="xtrabackup --backup --stream=xbstream --datadir=/var/lib/mysql"
export WALG_STREAM_RESTORE_COMMAND="xbstream -x -C /var/lib/mysql"
export WALG_MYSQL_BACKUP_PREPARE_COMMAND="xtrabackup --prepare --target-dir=/var/lib/mysql"
export WALG_MYSQL_BINLOG_REPLAY_COMMAND='mysqlbinlog --stop-datetime="$WALG_MYSQL_BINLOG_END_TS" "$WALG_MYSQL_CURRENT_BINLOG" | mysql'
export WALG_MYSQL_BINLOG_DST=/tmp/binlog

mysql -uroot -e "FLUSH BINARY LOGS;"

/usr/local/bin/wal-g backup-push
`))

var BackupMySQLFetchConfTemplate = template.Must(template.New("backup-fetch.sh").Parse(`#!/bin/bash
export WALG_MYSQL_DATASOURCE_NAME='root@unix(/var/run/mysqld/mysqld.sock)/mysql'
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export WALG_PGP_KEY_PATH=/etc/private.pgp
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"

export AWS_S3_FORCE_PATH_STYLE="true"
export AWS_ACCESS_KEY_ID="{{ .AwsKey }}"
export AWS_SECRET_ACCESS_KEY="{{ .AwsSecretKey }}"
export AWS_REGION="{{ .Region }}"

export WALG_STREAM_CREATE_COMMAND="xtrabackup --backup --stream=xbstream --datadir=/var/lib/mysql"
export WALG_STREAM_RESTORE_COMMAND="xbstream -x -C /var/lib/mysql"
export WALG_MYSQL_BACKUP_PREPARE_COMMAND="xtrabackup --prepare --target-dir=/var/lib/mysql"
export WALG_MYSQL_BINLOG_REPLAY_COMMAND='mysqlbinlog --stop-datetime="$WALG_MYSQL_BINLOG_END_TS" "$WALG_MYSQL_CURRENT_BINLOG" | mysql'
export WALG_MYSQL_BINLOG_DST=/tmp/binlog

/usr/local/bin/wal-g backup-fetch LATEST
`))

var BinLogPushConfTemplate = template.Must(template.New("binlog-push.sh").Parse(`#!/bin/bash
export WALG_MYSQL_DATASOURCE_NAME='root@unix(/var/run/mysqld/mysqld.sock)/mysql'
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export WALG_PGP_KEY_PATH=/etc/public.pgp
export WALG_S3_PREFIX='s3://{{ .Bucket }}/{{ .InstanceId }}'

export AWS_S3_FORCE_PATH_STYLE="true"
export AWS_ACCESS_KEY_ID='{{ .AwsKey }}'
export AWS_SECRET_ACCESS_KEY='{{ .AwsSecretKey }}'
export AWS_REGION='{{ .Region }}'

export WALG_STREAM_CREATE_COMMAND="xtrabackup --backup --stream=xbstream --datadir=/var/lib/mysql"
export WALG_STREAM_RESTORE_COMMAND="xbstream -x -C /var/lib/mysql"
export WALG_MYSQL_BACKUP_PREPARE_COMMAND="xtrabackup --prepare --target-dir=/var/lib/mysql"
export WALG_MYSQL_BINLOG_REPLAY_COMMAND='mysqlbinlog --stop-datetime="$WALG_MYSQL_BINLOG_END_TS" "$WALG_MYSQL_CURRENT_BINLOG" | mysql'
export WALG_MYSQL_BINLOG_DST=/tmp/binlog

mysql -uroot -e "FLUSH BINARY LOGS;"

if test -f /tmp/maintenance; then
    exit
fi

/usr/local/bin/wal-g binlog-push
`))

var BinLogReplayConfTemplate = template.Must(template.New("binlog-replay.sh").Parse(`#!/bin/bash
export WALG_MYSQL_DATASOURCE_NAME='root@unix(/var/run/mysqld/mysqld.sock)/mysql'
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export WALG_PGP_KEY_PATH=/etc/private.pgp
export WALG_S3_PREFIX='s3://{{ .Bucket }}/{{ .InstanceId }}'

export AWS_S3_FORCE_PATH_STYLE="true"
export AWS_ACCESS_KEY_ID='{{ .AwsKey }}'
export AWS_SECRET_ACCESS_KEY='{{ .AwsSecretKey }}'
export AWS_REGION='{{ .Region }}'

export WALG_STREAM_CREATE_COMMAND="xtrabackup --backup --stream=xbstream --datadir=/var/lib/mysql"
export WALG_STREAM_RESTORE_COMMAND="xbstream -x -C /var/lib/mysql"
export WALG_MYSQL_BACKUP_PREPARE_COMMAND="xtrabackup --prepare --target-dir=/var/lib/mysql"
export WALG_MYSQL_BINLOG_REPLAY_COMMAND='mysqlbinlog --stop-datetime="$WALG_MYSQL_BINLOG_END_TS" "$WALG_MYSQL_CURRENT_BINLOG" | mysql'
export WALG_MYSQL_BINLOG_DST=/tmp/binlog

/usr/local/bin/wal-g binlog-replay --since LATEST
`))
