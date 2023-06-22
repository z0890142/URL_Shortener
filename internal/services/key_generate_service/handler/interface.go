package handler

type KeyHandler interface {
	GenerateKey()
	GetKeys(num int) ([]string, error)
}
