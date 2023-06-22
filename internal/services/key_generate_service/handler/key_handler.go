package handler

import (
	"URL_Shortener/internal/data/key_data"
	"URL_Shortener/internal/models"
	"fmt"
	"sync"
)

type defaultKeyHandler struct {
	mu          sync.Mutex
	keyBuffer   chan string
	latestKeyId int64
	keyData     key_data.KeyData
}

func NewDefaultKeyHandler() (KeyHandler, error) {
	return &defaultKeyHandler{}, nil
}

func (d *defaultKeyHandler) GetKeys(num int) (result []string, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var count int
	var keys []models.KeyRow

Loop:
	for count = 0; count < num; count++ {
		select {
		case key := <-d.keyBuffer:
			result = append(result, key)
		default:
			keys, err = d.keyData.GetKey(num*2-count, d.latestKeyId)
			if err != nil {
				d.insertKeyToBuf(result)
				return []string{}, fmt.Errorf("GetKeys: %w", err)
			}
			break Loop
		}
	}

	if len(keys) == 0 {
		return result, nil
	}

	for index, key := range keys {
		keys[index].Used = 1
		if len(result) < num {
			result = append(result, key.Key)
			continue
		}
		d.keyBuffer <- key.Key
	}

	//update keys to used
	if _, err = d.keyData.UpdateKey(keys); err != nil {
		return []string{}, fmt.Errorf("GetKeys: %w", err)
	}

	d.latestKeyId = keys[len(keys)-1].Id

	return result, nil
}

func (d *defaultKeyHandler) GenerateKey() {

}

func (d *defaultKeyHandler) insertKeyToBuf(keys []string) {
	for _, key := range keys {
		d.keyBuffer <- key
	}
}
