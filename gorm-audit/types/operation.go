package types

// Operation 操作类型（共享定义）
type Operation string

const (
	OperationCreate Operation = "create"
	OperationUpdate Operation = "update"
	OperationDelete Operation = "delete"
	OperationQuery  Operation = "query"
)

// String 实现 Stringer 接口
func (o Operation) String() string {
	return string(o)
}

// IsValid 验证操作类型是否有效
func (o Operation) IsValid() bool {
	switch o {
	case OperationCreate, OperationUpdate, OperationDelete, OperationQuery:
		return true
	}
	return false
}
