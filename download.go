package embeded_postgresql

import (
	"path"
	"os"
	"fmt"
	"strings"
	"os/exec"
)

const (
	base_download_url = "http://get.enterprisedb.com/postgresql/postgresql-"
	zip = "zip"
	tar = "tar"
	linux = "linux"
	osx = "darwin"
)

var binaries = map[string]string {
	linux: "-linux-x64-binaries.tar.gz",
	osx: "-osx-binaries.zip",
}

type Dir string

func downloadPostgres(opSystem string, version string) string {

	var workdir = path.Join(os.Getenv("HOME"), ".postgres-embedded", version)
	var tmp = path.Join(workdir, "postgres.tmp")
	var dist = path.Join(workdir, "pgsql")

	if pathExists(dist) {
		return dist
	}

	if ok :=ensureDirectory(workdir); ok != nil {
		fmt.Println(ok.Error())
		panic("Can't create directory")
	}

	out, err := downloadAndExtract(
		workdir,
		tmp,
		strings.Join([]string{
			base_download_url,
			version,
			binaries[opSystem],
		}, ""))

	if err != nil {
		fmt.Println(out)
		panic(err.Error())
	}

	return dist
}

func downloadAndExtract(workDir string,tmp string, url string) (out string, e error) {
	os.Chdir(workDir)

	if out, err := exec.Command("wget", "-O", tmp, url).CombinedOutput(); err!=nil {
		return string(out), err
	} else {
		fmt.Println(out)
	}

	if strings.HasSuffix(url, zip) {
		if out, err := exec.Command("unzip", "-q", tmp).CombinedOutput(); err!=nil {
			return string(out), err
		} else {
			fmt.Println(out)
		}
	} else {
		if out, err := exec.Command("tar", "-xzf", tmp).CombinedOutput(); err!=nil {
			return string(out), err
		} else {
			fmt.Println(out)
		}
	}

	return "", nil
}

func ensureDirectory(dir string) error {
	if !pathExists(dir) {
		return os.MkdirAll(dir, 0700)
	}

	return nil
}

func pathExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	return true
}