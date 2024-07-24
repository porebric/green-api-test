package models

type Settings struct {
	InstanceId int64
	Settings   []Setting
}

type Setting struct {
	Id    int64
	Name  string
	Value string
}

func (m Settings) GetInstanceId() int64 {
	return m.InstanceId
}
