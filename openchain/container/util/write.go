package util

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"bufio"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"github.com/op/go-logging"
)

var vmLogger = logging.MustGetLogger("container")

func WriteGopathSrc(tw *tar.Writer, excludeDir string) error {
	gopath := os.Getenv("GOPATH")
	if strings.LastIndex(gopath, "/") == len(gopath)-1 {
		gopath = gopath[:len(gopath)]
	}
	rootDirectory := fmt.Sprintf("%s%s%s", os.Getenv("GOPATH"), string(os.PathSeparator), "src")
	vmLogger.Info("rootDirectory = %s", rootDirectory)

	//append "/" if necessary
	if excludeDir != "" && strings.LastIndex(excludeDir, "/") < len(excludeDir)-1 {
		excludeDir = excludeDir + "/"
	}

	rootDirLen := len(rootDirectory)
	walkFn := func(path string, info os.FileInfo, err error) error {

		// If path includes .git, ignore
		if strings.Contains(path, ".git") {
			return nil
		}

		if info.Mode().IsDir() {
			return nil
		}

		//exclude any files with excludeDir prefix. They should already be in the tar
		if excludeDir != "" && strings.Index(path, excludeDir) == rootDirLen+1 { //1 for "/"
			return nil
		}
		// Because of scoping we can reference the external rootDirectory variable
		newPath := fmt.Sprintf("src%s", path[rootDirLen:])
		//newPath := path[len(rootDirectory):]
		if len(newPath) == 0 {
			return nil
		}

		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		h, err := tar.FileInfoHeader(info, newPath)
		if err != nil {
			vmLogger.Error(fmt.Sprintf("Error getting FileInfoHeader: %s", err))
			return err
		}
		//Let's take the variance out of the tar, make headers identical everywhere by using zero time
		var zeroTime time.Time
		h.AccessTime = zeroTime
		h.ModTime = zeroTime
		h.ChangeTime = zeroTime
		h.Name = newPath
		if err = tw.WriteHeader(h); err != nil {
			vmLogger.Error(fmt.Sprintf("Error writing header: %s", err))
			return err
		}
		if _, err := io.Copy(tw, fr); err != nil {
			return err
		}
		return nil
	}

	if err := filepath.Walk(rootDirectory, walkFn); err != nil {
		vmLogger.Info("Error walking rootDirectory: %s", err)
		return err
	}
	// Write the tar file out
	if err := tw.Close(); err != nil {
		return err
	}
	//ioutil.WriteFile("/tmp/chaincode_deployment.tar", inputbuf.Bytes(), 0644)
	return nil
}

func WriteFileToPackage(fqpath string, filename string, tw *tar.Writer) error {
	fd, err := os.Open(fqpath)
	if err != nil {
		return fmt.Errorf("%s: %s", fqpath, err)
	}
	defer fd.Close()

	info, err := os.Stat(fqpath)
	if err != nil {
		return fmt.Errorf("%s: %s", fqpath, err)
	}

	is := bufio.NewReader(fd)

	header, err := tar.FileInfoHeader(info, fqpath)
	if err != nil {
		return fmt.Errorf("Error getting FileInfoHeader: %s", err)
	}

	//Let's take the variance out of the tar, make headers identical by using zero time
	var zeroTime time.Time
	header.AccessTime = zeroTime
	header.ModTime = zeroTime
	header.ChangeTime = zeroTime
	header.Name = filename

	tw.WriteHeader(header)
	if _, err := io.Copy(tw, is); err != nil {
		return fmt.Errorf("Error copying package into docker payload: %s", err)
	}

	return nil
}

// Find the instance of "bin" installed on the host's $PATH and inject it into the package
func WriteExecutableToPackage(bin string, tw *tar.Writer) error {
	cmd := exec.Command("which", bin)
	path, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Error determining %s path dynamically", bin)
	}

	return WriteFileToPackage(strings.Trim(string(path), "\n"), bin, tw)
}
