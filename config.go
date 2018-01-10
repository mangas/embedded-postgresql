package embeded_postgresql

type StartupConfig struct {
	cleanDir bool
	version  string
}

type RuntimeConfig struct {
	execDir string
	dataDir string
}

type DBConfig struct {
	port int
	username string
}