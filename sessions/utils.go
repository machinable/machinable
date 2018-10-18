package sessions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mssola/user_agent"
)

func getGeoIP(ip string) (string, error) {
	// ... this should be changed to get the access key from a config or environment variable
	accessKey := "85a38b87f3b696c7dcbf8f6f58c3c6a9"
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, accessKey)

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

func CreateSession(userID, ip, userAgent string, collection *mongo.Collection) (*Session, error) {
	location, _ := getGeoIP(ip)

	ua := user_agent.New(userAgent)

	bname, bversion := ua.Browser()
	session := &Session{
		ID:           objectid.New(),
		UserID:       userID,
		Location:     location,
		Mobile:       ua.Mobile(),
		IP:           ip,
		LastAccessed: time.Now(),
		Browser:      bname + " " + bversion,
		OS:           ua.OS(),
	}

	// save the user
	_, err := collection.InsertOne(
		context.Background(),
		session,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}
