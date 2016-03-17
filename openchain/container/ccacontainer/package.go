package ccacontainer

import (
	"archive/tar"
	"fmt"
	cutil "github.com/openblockchain/obc-peer/openchain/container/util"
	pb "github.com/openblockchain/obc-peer/protos"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func WritePackage(spec *pb.ChaincodeSpec, tw *tar.Writer) error {

	path, err := cutil.Download(spec.ChaincodeID.Path)
	if err != nil {
		return err
	}

	spec.ChaincodeID.Name, err = generateHashcode(spec, path)
	if err != nil {
		return fmt.Errorf("Error generating hashcode: %s", err)
	}

	copyobcc := viper.GetBool("chaincode.obcc.copyhost")

	buf := make([]string, 0)

	//let the executable's name be chaincode ID's name
	buf = append(buf, viper.GetString("chaincode.Dockerfile"))
	buf = append(buf, "COPY protoc-gen-go /usr/local/bin")
	if copyobcc {
		buf = append(buf, "COPY obcc /usr/local/bin")
	} else {
		buf = append(buf, "RUN apt-get install --yes obcc")
	}
	buf = append(buf, "COPY package.cca /tmp/package.cca")
	buf = append(buf, fmt.Sprintf("RUN obcc buildcca /tmp/package.cca -o $GOPATH/bin/%s && rm /tmp/package.cca", spec.ChaincodeID.Name))

	dockerFileContents := strings.Join(buf, "\n")
	dockerFileSize := int64(len([]byte(dockerFileContents)))

	//Make headers identical by using zero time
	var zeroTime time.Time
	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: zeroTime, AccessTime: zeroTime, ChangeTime: zeroTime})
	tw.Write([]byte(dockerFileContents))

	err = cutil.WriteFileToPackage(path, "package.cca", tw)
	if err != nil {
		return err
	}

	err = cutil.WriteExecutableToPackage("protoc-gen-go", tw)
	if err != nil {
		return err
	}

	if copyobcc {
		err := cutil.WriteExecutableToPackage("obcc", tw)
		if err != nil {
			return err
		}
	}

	err = cutil.WriteGopathSrc(tw, "")
	if err != nil {
		return fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}

	return nil
}
