package types

import (
	"github.com/netcloth/netcloth-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgContractCreate{}, "nch/MsgContractCreate", nil)
	cdc.RegisterConcrete(MsgContractCall{}, "nch/MsgContractCall", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
