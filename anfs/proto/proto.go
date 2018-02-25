package proto

//go:generate gencode go -package proto -schema proto.schema

const (
	TypeFile uint16 = iota
	TypeDir
	TypeSymlink
)