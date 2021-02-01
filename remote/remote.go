package remote

type Remote interface {
	Connect(pathKey string) error
	Close() error
	Load(address string, user string, port int, dir string, debug bool)
	SetDebug(debug bool)
	Run(cmd string) error
}
