package controller

import (
	"URL_Shortener/internal/models"
	"fmt"
	"regexp"
)

var matchHTTP = regexp.MustCompile(`^http://`)
var matchHTTPS = regexp.MustCompile(`^https://`)
var matchExpire = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func validReq(req *models.NewShortUrlRequest) error {

	if !matchExpire.MatchString(req.ExpireAt) {
		return fmt.Errorf("validReq: expireAt formant error")
	}

	if matchHTTPS.MatchString(req.Url) {
		return nil
	} else if matchHTTP.MatchString(req.Url) {
		return nil
	}

	return fmt.Errorf("validReq: url formant error")
}
