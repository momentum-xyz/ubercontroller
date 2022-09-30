package auth

import (
	// Std
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func VerifyToken(token string, introspectUrl string) error {
	type introPayload struct {
		Token string `json:"token"`
	}
	client := &http.Client{}
	bearer := "Bearer " + token

	payload := introPayload{Token: token}
	payloadBody, err := json.Marshal(payload)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal payload")
	}

	req, err := http.NewRequest("POST", introspectUrl, bytes.NewBuffer(payloadBody))
	if err != nil {
		return errors.WithMessage(err, "failed to create request")
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, "failed to do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("invalid status code: %s", resp.Status)
	}

	ret := make([]byte, resp.ContentLength)
	resp.Body.Read(ret)
	if string(ret) != "1" {
		return errors.Errorf("invalid ret: %s", ret)
	}

	return nil
}
