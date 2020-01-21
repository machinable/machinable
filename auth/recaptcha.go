package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anothrnick/machinable/config"
)

// RecaptchaSiteVerify verifies the client response with recaptcha
func RecaptchaSiteVerify(clientResponse string) error {

	resp, err := http.Post(fmt.Sprintf("https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s", config.ReCaptchaSecret, clientResponse), "application/json", nil)
	if err != nil {
		log.Fatalln(err)
	}

	var result struct {
		Success     bool        `json:"success"`
		ChallengeTs time.Time   `json:"challenge_ts"`
		Hostname    string      `json:"hostname"`
		ErrorCods   interface{} `json:"error-codes"`
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("unexpected status code from recaptcha verification: %d", resp.StatusCode))
		return errors.New("invalid reCaptcha")
	}

	json.NewDecoder(resp.Body).Decode(&result)

	if !result.Success {
		return errors.New("invalid reCaptcha")
	}

	return nil
}
