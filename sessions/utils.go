package sessions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mssola/user_agent"
)

// GetGeoIP retrieves location information of the client ip from IP Stack
func GetGeoIP(ip, apiKey string) (string, error) {
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey)

	ipStackData := struct {
		City       string `json:"city"`
		RegionCode string `json:"region_code"`
	}{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("error creating request")
	}

	// set client with 10 second timeout
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("error making request")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ipStackData); err != nil {
		return "", errors.New("error decoding response")
	}

	location := ""

	if ipStackData.City != "" && ipStackData.RegionCode != "" {
		location = ipStackData.City + ", " + ipStackData.RegionCode
	}

	return location, nil
}

// CreateSession creates a new session object after looking up geoip
func CreateSession(userID, ip, userAgent string, config *config.AppConfig) *models.Session {
	location, _ := GetGeoIP(ip, config.IPStackKey)

	ua := user_agent.New(userAgent)

	bname, bversion := ua.Browser()
	session := &models.Session{
		UserID:       userID,
		Location:     location,
		Mobile:       ua.Mobile(),
		IP:           ip,
		LastAccessed: time.Now(),
		Browser:      bname + " " + bversion,
		OS:           ua.OS(),
	}

	return session
}
