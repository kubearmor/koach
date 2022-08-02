package model

type Operation string

const (
	OperationFileAccess   Operation = "File"
	OperationProcessSpawn Operation = "Process"
	OperationNetworkCall  Operation = "Network"
)

type OperationDetail interface {
	IsOperationDetail()
}
