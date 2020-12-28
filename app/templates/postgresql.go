package templates

import "text/template"

const PgAgentConf = `DB=pg
GOGC=40
POSTGRES_URL=postgres://postgres:%s@localhost:5432
SECRET=%s
`

var UfwPgTemplate = template.Must(template.New("pg").Parse(`*filter
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
### tuple ### allow any 5432 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 5432 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 5432 -s {{.}} -j ACCEPT
### END RULES ###
{{end}}
{{range .Instance.Whitelist}}
### tuple ### allow any 5432 0.0.0.0/0 any {{.}} in
-A ufw-user-input -p tcp --dport 5432 -s {{.}} -j ACCEPT
-A ufw-user-input -p udp --dport 5432 -s {{.}} -j ACCEPT
{{end}}
{{else}}
### tuple ### allow any 5432 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 5432 -j ACCEPT
-A ufw-user-input -p udp --dport 5432 -j ACCEPT
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
