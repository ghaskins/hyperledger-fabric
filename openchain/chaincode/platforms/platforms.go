package platforms

import (
	"archive/tar"
	"fmt"
	"github.com/openblockchain/obc-peer/openchain/chaincode/platforms/golang"
	"github.com/openblockchain/obc-peer/openchain/chaincode/platforms/cca"
	pb "github.com/openblockchain/obc-peer/protos"
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
