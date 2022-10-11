package model

type Operation string

const (
	OperationFileAccess   Operation = "File"
	OperationProcessSpawn Operation = "Process"
	OperationNetworkCall  Operation = "Network"
	OperationSystemCall   Operation = "Syscall"
)
