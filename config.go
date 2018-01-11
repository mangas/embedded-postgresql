package embeded_postgresql

type StartupConfig struct {
	CleanDir bool
	Version  string
}

type RuntimeConfig struct {
	ExecDir string
	DataDir string
}

type DBConfig struct {
	Port int
	Username string
}