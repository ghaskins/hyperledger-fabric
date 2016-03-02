package util

import (
	"strings"
	"io/ioutil"
	"io"
	"os"
	"fmt"
	"net/http"
)

func Download(path string) (string, error) {
	if strings.HasPrefix(path, "http://") {
		// The file is remote, so we need to download it to a temporary location first

		var tmp *os.File
		var err error
		tmp, err = ioutil.TempFile("", "obc")
		if err != nil {
			return "", fmt.Errorf("Error creating temporary file: %s", err)
		}
		defer os.Remove(tmp.Name())
		defer tmp.Close()

		resp, err := http.Get(path)
		if err != nil {
			return "", fmt.Errorf("Error with HTTP GET: %s", err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(tmp, resp.Body)
		if err != nil {
			return "", fmt.Errorf("Error downloading bytes: %s", err)
		}

		return tmp.Name(), nil
	} else {
		return path, nil
	}
}
