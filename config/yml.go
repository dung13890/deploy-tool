package config

var SourceYaml = `server:
  address: "{{ .Ip }}"
  user: "ec2-user"
  port: 22
  project: "{{ .Project }}"
`
