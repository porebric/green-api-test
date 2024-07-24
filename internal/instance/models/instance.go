package models

type Instance struct {
	Id    int64
	Name  string
	Token string
}

func (m Instance) GetInstanceId() int64 {
	return m.Id
}
