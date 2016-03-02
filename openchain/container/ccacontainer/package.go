package ccacontainer

import (
	"archive/tar"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
	"bufio"
	"os/exec"

	"golang.org/x/net/context"

	"github.com/spf13/viper"

	"github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"


	cutil "github.com/openblockchain/obc-peer/openchain/container/util"
	pb "github.com/openblockchain/obc-peer/protos"
)

func writeChaincodePackage(spec *pb.ChaincodeSpec, path string, tw *tar.Writer) error {

	buf := make([]string, 0)

	//let the executable's name be chaincode ID's name
	buf = append(buf, "FROM chaincode-base")
	buf = append(buf, "COPY package.cca /tmp/package.cca")
	buf = append(buf, fmt.Sprintf("RUN obcc buildcca /tmp/package.cca -o $GOPATH/bin/%s && rm /tmp/package.cca", spec.ChaincodeID.Name))

	dockerFileContents := strings.Join(buf, "\n")
	dockerFileSize := int64(len([]byte(dockerFileContents)))

	//Make headers identical by using zero time
	var zeroTime time.Time
	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: zeroTime, AccessTime: zeroTime, ChangeTime: zeroTime})
	tw.Write([]byte(dockerFileContents))

	var err error
	err = cutil.WriteFileToPackage(path, "package.cca", tw)
	if err != nil {
		return err
	}

	return nil
}

func WritePackage(spec *pb.ChaincodeSpec, tw *tar.Writer) error {

	path, err := cutil.Download(spec.ChaincodeID.Path)
	if err != nil {
		return nil
	}

	spec.ChaincodeID.Name, err = generateHashcode(spec, path)
	if err != nil {
		return nil, fmt.Errorf("Error generating hashcode: %s", err)
	}

	err = writeChaincodePackage(spec, path, tw)
	if err != nil {
		return nil, fmt.Errorf("Error writing chaincode package: %s", err)
	}

	return nil
}
