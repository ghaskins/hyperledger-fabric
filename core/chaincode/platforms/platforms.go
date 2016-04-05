package platforms

import (
	"archive/tar"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/platforms/cca"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	pb "github.com/hyperledger/fabric/protos"
)

type Platform interface {
	ValidateSpec(spec *pb.ChaincodeSpec) error
	WritePackage(spec *pb.ChaincodeSpec, tw *tar.Writer) error
}

func Find(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {

	switch chaincodeType {
	case pb.ChaincodeSpec_GOLANG:
		return &golang.Platform{}, nil
	case pb.ChaincodeSpec_CCA:
		return &cca.Platform{}, nil
	default:
		return nil, fmt.Errorf("Unknown chaincodeType: %s", chaincodeType)
	}

}
