package remote

type Remote interface {
	Connect(pathKey string) error
	Close() error
	Load(address string, user string, port int, dir string)
	SetDebug(debug bool)
	Run(cmd string) error
}
