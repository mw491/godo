package model

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type Todo struct {
	Id   uint64
	Task string
	Done bool
}
