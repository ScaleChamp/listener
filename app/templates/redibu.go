package templates

import (
	uuid "github.com/satori/go.uuid"
	"text/template"
)

type RedibuConf struct {
	InstanceId  uuid.UUID
	AccessKeyId string
	SecretKey   string
	AwsRegion   string
	Bucket      string
}

var RedibuConfTemplate = template.Must(template.New("redibu.conf").Parse(`S3_BUCKET={{ .Bucket }}
S3_PREFIX={{ .InstanceId }}
AWS_ACCESS_KEY_ID={{ .AccessKeyId }}
AWS_SECRET_KEY={{ .SecretKey }}
AWS_REGION={{ .AwsRegion }}
DUMP_PATH=/var/lib/redis/dump.rdb
PGP_PRIVATE_KEY_PATH=/etc/private.pgp
PGP_PUBLIC_KEY_PATH=/etc/public.pgp
GOGC=40
`))
