package model

type Operation string

const (
	OperationFileAccess   Operation = "FILE_ACCESS"
	OperationProcessSpawn Operation = "PROCESS_SPAWN"
	OperationNetworkCall  Operation = "NETWORK_CALL"
)

type OperationDetail interface {
	IsOperationDetail()
}
