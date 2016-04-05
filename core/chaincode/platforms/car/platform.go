package car

import (
	pb "github.com/hyperledger/fabric/protos"
)

type Platform struct {
}

func (self *Platform) ValidateSpec(spec *pb.ChaincodeSpec) error {
	return nil
}
