package handler

import "fmt"

type KeyHandler interface {
	GetKeys(num int) ([]string, error)
	Shutdown()
}

func NewKeyHandler(conf interface{}) (KeyHandler, error) {
	switch c := conf.(type) {
	case *DefaultKeyHandlerConf:
		return newDefaultKeyHandler(c)
	}
	return nil, fmt.Errorf("config type not found")
}
