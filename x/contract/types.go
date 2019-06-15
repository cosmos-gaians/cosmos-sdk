package contract

// Defines contract module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Contract stores data about a contract
type Contract struct {
	ByteCode []byte
}
