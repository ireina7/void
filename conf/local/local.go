package local

// A configuration of local server
type LocalConf struct {
	version string
	host    string
	port    string
}

func (self *LocalConf) Version() string {
	return self.version
}

func (self *LocalConf) Host() string {
	return self.host
}

func (self *LocalConf) Port() string {
	return self.port
}

// Local configuration
func Instance() LocalConf {
	return LocalConf{
		version: "0.1.0",
		host:    "localhost",
		port:    "3210",
	}
}
