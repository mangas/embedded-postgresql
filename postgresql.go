package embeded_postgresql

import (
	"runtime"
	"fmt"
	"os"
	"path"
	"os/exec"
	"time"
)

func StartPostgres(startConfig StartupConfig, dbConfig DBConfig) RuntimeConfig {
	downloadDir := downloadPostgres(runtime.GOOS, startConfig.Version)
	dataDir := path.Join(downloadDir, "data")
	dirExists := pathExists(dataDir)

	if dirExists && !startConfig.CleanDir {
		panic(fmt.Sprintf("Directory %s must not exist", downloadDir))
	}

	os.RemoveAll(dataDir)

	out, err := exec.Command(
		path.Join(downloadDir, "bin", "initdb"),
		"-A trust", "-U", dbConfig.Username, "-D", dataDir, "-E UTF-8",
	).CombinedOutput()

	if err == nil {
		fmt.Println(string(out))
	} else {
		fmt.Println(string(out))
		panic(err)
	}

	cmd := exec.Command(
		path.Join(downloadDir, "bin", "pg_ctl"),
		"-w", "-D", dataDir, fmt.Sprintf("-o -F -p %v", dbConfig.Port),
		"-l", path.Join(downloadDir, "log"), "start",
	)

	out, err = cmd.CombinedOutput()

	if err == nil {
		fmt.Println(string(out))
	} else {
		fmt.Println(string(out))
		panic(err)
	}

	rc := RuntimeConfig{downloadDir, dataDir}

	checkPostgresStarted(rc, dbConfig)

	return rc
}

func StopPostGres(runConfig RuntimeConfig) error {
	out, err := exec.Command(
		path.Join(
			runConfig.ExecDir, "bin", "pg_ctl"),
			"-D", runConfig.DataDir, "stop", "--mode=fast", "-t 5", "-w",
			).CombinedOutput()

	if err == nil {
		fmt.Println(string(out))
	} else {
		fmt.Println(string(out))
		panic(fmt.Sprintf("Can't stop database: %v", err))
	}

	return err
}

func checkPostgresStarted(runtimeConfig RuntimeConfig, dbConfig DBConfig) error {

	var err error

	for i := 0; i<3; {
		i++
		time.Sleep(1*time.Second)

		out, err := exec.Command(
			path.Join(runtimeConfig.ExecDir, "bin", "pg_isready"),
			fmt.Sprintf("--port=%v ", dbConfig.Port), "--host=127.0.0.1",
		).CombinedOutput()

		if err == nil {
			fmt.Println("Database is accepting connections...")
			return nil
		} else {
			fmt.Println(string(out))
			fmt.Println(fmt.Sprintf("Database failed to start with: %v", err))
		}
	}

	return err
}