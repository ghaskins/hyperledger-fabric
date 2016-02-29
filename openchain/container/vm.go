/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package container

import (
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
	pb "github.com/openblockchain/obc-peer/protos"
)

func newDockerClient() (client *docker.Client, err error) {
	//QQ: is this ok using config properties here so deep ? ie, should we read these in main and stow them away ?
	endpoint := viper.GetString("vm.endpoint")
	vmLogger.Info("Creating VM with endpoint: %s", endpoint)
	if viper.GetBool("vm.docker.tls.enabled") {
		cert := viper.GetString("vm.docker.tls.cert.file")
		key := viper.GetString("vm.docker.tls.key.file")
		ca := viper.GetString("vm.docker.tls.ca.file")
		client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	} else {
		client, err = docker.NewClient(endpoint)
	}
	return
}

// VM implemenation of VM management functionality.
type VM struct {
	Client *docker.Client
}

// NewVM creates a new VM instance.
func NewVM() (*VM, error) {
	client, err := newDockerClient()
	if err != nil {
		return nil, err
	}
	VM := &VM{Client: client}
	return VM, nil
}

var vmLogger = logging.MustGetLogger("container")

// ListImages list the images available
func (vm *VM) ListImages(context context.Context) error {
	imgs, err := vm.Client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return err
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}

	return nil
}

// BuildChaincodeContainer builds the container for the supplied chaincode specification
func (vm *VM) BuildChaincodeContainer(spec *pb.ChaincodeSpec) ([]byte, error) {
	chaincodePkgBytes, err := GetChaincodePackageBytes(spec)
	if err != nil {
		return nil, fmt.Errorf("Error getting chaincode package bytes: %s", err)
	}
	err = vm.buildChaincodeContainerUsingDockerfilePackageBytes(spec, chaincodePkgBytes)
	if err != nil {
		return nil, fmt.Errorf("Error building Chaincode container: %s", err)
	}
	return chaincodePkgBytes, nil
}

// GetChaincodePackageBytes creates bytes for docker container generation using the supplied chaincode specification
func GetChaincodePackageBytes(spec *pb.ChaincodeSpec) ([]byte, error) {
	if spec == nil || spec.ChaincodeID == nil {
		return nil, fmt.Errorf("invalid chaincode spec")
	}
	if spec.ChaincodeID.Path == "" {
		return nil, fmt.Errorf("Cannot generate hashcode from empty chaincode path")
	}
	if spec.ChaincodeID.Name != "" {
		return nil, fmt.Errorf("chaincode name exists")
	}

	path := spec.ChaincodeID.Path
	var localpath string
	var err error

	if strings.HasPrefix(path, "http://") {
		// The file is remote, so we need to download it to a temporary location first

		var tmp *os.File
		tmp, err = ioutil.TempFile("", "cca")
		if err != nil {
			return nil, fmt.Errorf("Error creating temporary file: %s", err)
		}
		defer os.Remove(tmp.Name())
		defer tmp.Close()

		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("Error with HTTP GET: %s", err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(tmp, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error downloading bytes: %s", err)
		}

		localpath = tmp.Name()
	} else {
		localpath = path
	}

	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	spec.ChaincodeID.Name, err = generateHashcode(spec, localpath)
	if err != nil {
		return nil, fmt.Errorf("Error generating hashcode: %s", err)
	}
	err = writeChaincodePackage(spec, localpath, tw)
	if err != nil {
		return nil, fmt.Errorf("Error writing chaincode package: %s", err)
	}

	tw.Close()
	gw.Close()

	chaincodePkgBytes := inputbuf.Bytes()

	return chaincodePkgBytes, nil
}

// Builds the Chaincode image using the supplied Dockerfile package contents
func (vm *VM) buildChaincodeContainerUsingDockerfilePackageBytes(spec *pb.ChaincodeSpec, code []byte) error {
	outputbuf := bytes.NewBuffer(nil)
	vmName := spec.ChaincodeID.Name
	inputbuf := bytes.NewReader(code)
	opts := docker.BuildImageOptions{
		Name:         vmName,
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := vm.Client.BuildImage(opts); err != nil {
		vmLogger.Debug(fmt.Sprintf("Failed Chaincode docker build:\n%s\n", outputbuf.String()))
		return fmt.Errorf("Error building Chaincode container: %s", err)
	}
	return nil
}

// BuildPeerContainer builds the image for the Peer to be used in development network
func (vm *VM) BuildPeerContainer() error {
	//inputbuf, err := vm.GetPeerPackageBytes()
	inputbuf, err := vm.getPackageBytes(vm.writePeerPackage)

	if err != nil {
		return fmt.Errorf("Error building Peer container: %s", err)
	}
	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         "openchain-peer",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := vm.Client.BuildImage(opts); err != nil {
		vmLogger.Debug(fmt.Sprintf("Failed Peer docker build:\n%s\n", outputbuf.String()))
		return fmt.Errorf("Error building Peer container: %s\n", err)
	}
	return nil
}

// BuildObccaContainer builds the image for the obcca to be used in development network
func (vm *VM) BuildObccaContainer() error {
	inputbuf, err := vm.getPackageBytes(vm.writeObccaPackage)

	if err != nil {
		return fmt.Errorf("Error building obcca container: %s", err)
	}
	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         "obcca",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := vm.Client.BuildImage(opts); err != nil {
		vmLogger.Debug(fmt.Sprintf("Failed obcca docker build:\n%s\n", outputbuf.String()))
		return fmt.Errorf("Error building obcca container: %s\n", err)
	}
	return nil
}

// BuildChaincodeBaseContainer builds the base-image for the chaincode
func (vm *VM) BuildChaincodeBaseContainer() error {
	inputbuf, err := vm.getPackageBytes(vm.writeChaincodeBasePackage)

	if err != nil {
		return fmt.Errorf("Error building chaincode-base container: %s", err)
	}
	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         "chaincode-base",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := vm.Client.BuildImage(opts); err != nil {
		vmLogger.Debug(fmt.Sprintf("Failed chaincode-base docker build:\n%s\n", outputbuf.String()))
		return fmt.Errorf("Error building chaincode-base container: %s\n", err)
	}
	return nil
}

// GetPeerPackageBytes returns the gzipped tar image used for docker build of Peer
func (vm *VM) GetPeerPackageBytes() (io.Reader, error) {
	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tr := tar.NewWriter(gw)
	// Get the Tar contents for the image
	err := vm.writePeerPackage(tr)
	tr.Close()
	gw.Close()
	if err != nil {
		return nil, fmt.Errorf("Error getting Peer package: %s", err)
	}
	return inputbuf, nil
}

//type tarWriter func()

func (vm *VM) getPackageBytes(writerFunc func(*tar.Writer) error) (io.Reader, error) {
	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tr := tar.NewWriter(gw)
	// Get the Tar contents for the image
	err := writerFunc(tr)
	tr.Close()
	gw.Close()
	if err != nil {
		return nil, fmt.Errorf("Error getting package bytes: %s", err)
	}
	return inputbuf, nil
}

func writeFileToPackage(fqpath string, filename string, tw *tar.Writer) error {
	info, err := os.Lstat(fqpath)
	if err != nil {
		return fmt.Errorf("Error lstat on %s: %s", fqpath, err)
	}

	fd, err := os.Open(fqpath)
	if err != nil {
		return fmt.Errorf("Error opening %s: %s", fqpath, err)
	}
	defer fd.Close()

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

// Find the instance of obcc installed on the host's $PATH and inject it into the package
func writeObccToPackage(tw *tar.Writer) error {
	cmd := exec.Command("which", "obcc")
	path, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Error determining obcc path dynamically")
	}

	return writeFileToPackage(string(path), "obcc", tw)
}

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
	err = writeFileToPackage(path, "package.cca", tw)
	if err != nil {
		return err
	}

	return nil
}

func (vm *VM) writePeerPackage(tw *tar.Writer) error {
	startTime := time.Now()

	dockerFileContents := viper.GetString("peer.Dockerfile")
	dockerFileSize := int64(len([]byte(dockerFileContents)))

	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: startTime, AccessTime: startTime, ChangeTime: startTime})
	tw.Write([]byte(dockerFileContents))
	err := writeGopathSrc(tw, "")
	if err != nil {
		return fmt.Errorf("Error writing Peer package contents: %s", err)
	}
	return nil
}

func (vm *VM) writeObccaPackage(tw *tar.Writer) error {
	startTime := time.Now()

	dockerFileContents := viper.GetString("peer.Dockerfile")
	dockerFileContents = dockerFileContents + "WORKDIR obc-ca\nRUN go install && cp obcca.yaml $GOPATH/bin\n"
	dockerFileSize := int64(len([]byte(dockerFileContents)))

	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: startTime, AccessTime: startTime, ChangeTime: startTime})
	tw.Write([]byte(dockerFileContents))
	err := writeGopathSrc(tw, "")
	if err != nil {
		return fmt.Errorf("Error writing obcca package contents: %s", err)
	}
	return nil
}

func (vm *VM) writeChaincodeBasePackage(tw *tar.Writer) error {

	copyobcc := viper.GetBool("chaincode.obcc.copyhost")

	buf := make([]string, 0)

	buf = append(buf, viper.GetString("chaincode.Dockerfile"))
	if copyobcc {
		buf = append(buf, "COPY obcc /usr/local/bin")
	} else {
		buf = append(buf, "RUN apt-get install --yes obcc")
	}

	dockerFileContents := strings.Join(buf, "\n")
	dockerFileSize := int64(len([]byte(dockerFileContents)))

	//Make headers identical by using zero time
	var zeroTime time.Time
	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: zeroTime, AccessTime: zeroTime, ChangeTime: zeroTime})
	tw.Write([]byte(dockerFileContents))

	if copyobcc {
		err := writeObccToPackage(tw)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeGopathSrc(tw *tar.Writer, excludeDir string) error {
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
