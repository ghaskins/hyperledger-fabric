/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package car

import (
	"archive/tar"
	"io/ioutil"
	"strings"

	"bytes"
	"fmt"
	"io"

	"github.com/fsouza/go-dockerclient"
	cutil "github.com/hyperledger/fabric/core/container/util"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Platform for the CAR type
type Platform struct {
}

// ValidateSpec validates the chaincode specification for CAR types to satisfy
// the platform interface.  This chaincode type currently doesn't
// require anything specific so we just implicitly approve any spec
func (carPlatform *Platform) ValidateSpec(spec *pb.ChaincodeSpec) error {
	return nil
}

func (carPlatform *Platform) ValidateDeploymentSpec(cds *pb.ChaincodeDeploymentSpec) error {
	// CAR platform will validate the code package within chaintool
	return nil
}

func (carPlatform *Platform) GetDeploymentPayload(spec *pb.ChaincodeSpec) ([]byte, error) {

	return ioutil.ReadFile(spec.ChaincodeId.Path)
}

func (carPlatform *Platform) GenerateDockerfile(cds *pb.ChaincodeDeploymentSpec) (string, error) {

	var buf []string

	//let the executable's name be chaincode ID's name
	buf = append(buf, "FROM "+cutil.GetDockerfileFromConfig("chaincode.car.runtime"))
	buf = append(buf, "ADD chaincode.tar /usr/local/bin")

	dockerFileContents := strings.Join(buf, "\n")

	return dockerFileContents, nil
}

func (carPlatform *Platform) GenerateDockerBuild(cds *pb.ChaincodeDeploymentSpec, tw *tar.Writer) error {
	client, err := cutil.NewDockerClient()
	if err != nil {
		return fmt.Errorf("Error creating docker client: %s", err)
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        cutil.GetDockerfileFromConfig("chaincode.car.builder"),
			Cmd:          []string{"java", "-jar", "/usr/local/bin/chaintool", "buildcar", "/tmp/codepackage.car", "-o", "/tmp/output/chaincode"},
			AttachStdout: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Error creating container: %s", err)
	}
	defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})

	codepackage, output := io.Pipe()

	go func() {
		tw := tar.NewWriter(output)

		err := cutil.WriteBytesToPackage("codepackage.car", cds.CodePackage, tw)

		tw.Close()
		output.CloseWithError(err)
	}()

	err = client.UploadToContainer(container.ID, docker.UploadToContainerOptions{
		Path:        "/tmp",
		InputStream: codepackage,
	})
	if err != nil {
		return fmt.Errorf("Error uploading codepackage to container: %s", err)
	}

	stdout := bytes.NewBuffer(nil)
	_, err = client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: stdout,
		Logs:         true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		return fmt.Errorf("Error attaching to container: %s", err)
	}

	err = client.StartContainer(container.ID, nil)
	if err != nil {
		return fmt.Errorf("Error building CAR: %s \"%s\"", err, stdout.String())
	}

	_, err = client.WaitContainer(container.ID)
	if err != nil {
		return fmt.Errorf("Error waiting for container to complete: %s", err)
	}

	payload := bytes.NewBuffer(nil)
	err = client.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
		Path:         "/tmp/output/.",
		OutputStream: payload,
	})
	if err != nil {
		return fmt.Errorf("Error downloading payload: %s", err)
	}

	return cutil.WriteBytesToPackage("chaincode.tar", payload.Bytes(), tw)
}
