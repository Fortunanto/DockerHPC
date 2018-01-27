package Interfaces

type BackingStore interface {
	InsertValue(value interface{}, keyHashParameters ...string) (string, error)
	GetValue(key string) (interface{}, error)
	GenerateKeyHash(string) string
}
