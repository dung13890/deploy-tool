package config

var SourceYaml = `server:
	address: "{{ .Ip }}"
	user: "ec2-user"
	port: 22
	dir: "/data/sites"
	project: "{{ .Project }}"
repository:
	url: "{{ .Repo }}"
	branch: "master"
	# tag: "1.0.2"
shared:
	folders:
		- "vendor"
		- "node_modules"
	files:
		- ".env"
tasks:
	- "echo 'Hello!!'"
cluster:
	hosts:
		# - "127.0.0.1"
	rsync:
		excludes:
			- ".env"
			- "*.log"
	cmds:
		- "uname -a"
`
