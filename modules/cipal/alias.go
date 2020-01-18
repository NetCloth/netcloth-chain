package cipal

import (
	"github.com/netcloth/netcloth-chain/modules/cipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = keeper.DefaultParamspace
)

var (
	RegisterCodec                     = types.RegisterCodec
	NewIPALObject                     = types.NewCIPALObject
	NewQuerier                        = keeper.NewQuerier
	NewADParam                        = types.NewADParam
	NewIPALUserRequest                = types.NewCIPALUserRequest
	NewMsgIPALClaim                   = types.NewMsgCIPALClaim
	NewKeeper                         = keeper.NewKeeper
	ErrEmptyInputs                    = types.ErrEmptyInputs
	ErrStringTooLong                  = types.ErrStringTooLong
	ErrInvalidSignature               = types.ErrInvalidSignature
	ErrIPALClaimUserRequestExpired    = types.ErrIPALClaimUserRequestExpired
	ErrCIPALClaimUserRequestSigVerify = types.ErrCIPALClaimUserRequestSigVerify
	ModuleCdc                         = types.ModuleCdc
	AttributeValueCategory            = types.AttributeValueCategory
)

type (
	Keeper          = keeper.Keeper
	MsgIPALClaim    = types.MsgCIPALClaim
	IPALUserRequest = types.CIPALUserRequest
	ADParam         = types.ADParam
	CIPALObject     = types.CIPALObject
)
