package templates

import (
	uuid "github.com/satori/go.uuid"
	"text/template"
)

type WalGConf struct {
	Version      string
	Password     string
	AwsKey       string
	AwsSecretKey string
	Bucket       string
	InstanceId   uuid.UUID
	Region       string
}

var BackupFetchConfTemplate = template.Must(template.New("backup-fetch.conf").Parse(`#!/bin/bash
export PGHOST=localhost
export PGPORT=5432
export PGUSER=postgres
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export AWS_S3_FORCE_PATH_STYLE="true"
export WALG_PGP_KEY_PATH=/etc/private.pgp
export PGPASSWORD={{ .Password }}
export AWS_ACCESS_KEY_ID={{ .AwsKey }}
export AWS_SECRET_ACCESS_KEY={{ .AwsSecretKey }}
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"
export AWS_REGION={{ .Region }}

/usr/local/bin/wal-g backup-fetch /var/lib/postgresql/{{ .Version }}/main LATEST
`))

var BackupPushConfTemplate = template.Must(template.New("backup-push.conf").Parse(`#!/bin/bash
export PGHOST=localhost
export PGPORT=5432
export PGUSER=postgres
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export AWS_S3_FORCE_PATH_STYLE="true"
export WALG_PGP_KEY_PATH=/etc/public.pgp

export PGPASSWORD={{ .Password }}
export AWS_ACCESS_KEY_ID={{ .AwsKey }}
export AWS_SECRET_ACCESS_KEY={{ .AwsSecretKey }}
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"
export AWS_REGION={{ .Region }}

/usr/local/bin/wal-g backup-push /var/lib/postgresql/{{ .Version }}/main
`))

var WalPushConfTemplate = template.Must(template.New("wal-push.conf").Parse(`#!/bin/bash
export PGHOST=localhost
export PGPORT=5432
export PGUSER=postgres
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export WALG_PGP_KEY_PATH=/etc/public.pgp
export AWS_S3_FORCE_PATH_STYLE="true"

export PGPASSWORD={{ .Password }}
export AWS_ACCESS_KEY_ID={{ .AwsKey }}
export AWS_SECRET_ACCESS_KEY={{ .AwsSecretKey }}
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"
export AWS_REGION={{ .Region }}

if test -f /tmp/maintenance; then
    exit
fi

/usr/local/bin/wal-g wal-push $1
`))

var WalFetchConfTemplate = template.Must(template.New("wal-fetch.conf").Parse(`#!/bin/bash
export PGHOST=localhost
export PGPORT=5432
export PGUSER=postgres
export WALG_DOWNLOAD_CONCURRENCY=2
export WALG_UPLOAD_CONCURRENCY=2
export WALG_DELTA_MAX_STEPS=7
export WALG_COMPRESSION_METHOD=brotli
export AWS_S3_FORCE_PATH_STYLE="true"
export WALG_PGP_KEY_PATH=/etc/private.pgp

export PGPASSWORD={{ .Password }}
export AWS_ACCESS_KEY_ID={{ .AwsKey }}
export AWS_SECRET_ACCESS_KEY={{ .AwsSecretKey }}
export WALG_S3_PREFIX="s3://{{ .Bucket }}/{{ .InstanceId }}"
export AWS_REGION={{ .Region }}

/usr/local/bin/wal-g wal-fetch $1 $2
`))
