package types

const (
	// ModuleName is the name of the vm module
	ModuleName = "vm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	CodeKey = StoreKey + "_code"

	// QuerierRoute is the querier route for the vm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the vm module
	RouterKey = ModuleName
)
