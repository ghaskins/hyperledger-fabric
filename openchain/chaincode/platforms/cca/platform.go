package cca

import (
	pb "github.com/openblockchain/obc-peer/protos"
)

type Platform struct {
}

func (self *Platform) ValidateSpec(spec *pb.ChaincodeSpec) error {
	return nil
}
