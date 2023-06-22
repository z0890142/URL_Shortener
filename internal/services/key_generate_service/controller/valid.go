package controller

import (
	"URL_Shortener/internal/models"
	"fmt"
)

func validReq(req *models.GetKeysRequest) error {

	if req.Nums <= 0 {
		return fmt.Errorf("validReq: nums must be greater than 0")
	}
	return nil
}
