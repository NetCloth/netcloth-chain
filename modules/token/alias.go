package token

import (
	"github.com/NetCloth/netcloth-chain/modules/token/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	DefaultCodespace   = types.DefaultCodespace
	CodeInvalidMoniker = types.CodeInvalidMoniker

	ErrInvalidMoniker = types.ErrInvalidMoniker
	NewMsgIssue       = types.NewMsgIssue
)

type (
	MsgIssue = types.MsgIssue
)