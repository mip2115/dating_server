package zoom

import (
	"encoding/json"
	//"github.com/gorilla/mux"
	"../../types"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// gets the zoom access token
func ZoomCallback(w http.ResponseWriter, r *http.Request) {

	res := make(map[string]interface{})
	res["res"] = "Hello world"

	/*
		u, _ := url.Parse(r.URL.RawQuery)
		fmt.Println(u)

		values, _ := url.ParseQuery(r.URL.RawQuery)

		code := values.Get("code")
	*/
	code := getCode(r.URL.RawQuery)
	zat, err := getAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(*zat)
	res["code"] = code

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func SetZoomMeeting() {

	// https://api.zoom.us/v2/users/me/meetings
	// Authoriation bearer accesstoken

	/*
		{
	  "topic": "Sex education",
	  "type": 2,
	  "start_time": "2020-06-08T12:15:00Z",
	  "duration": 30,
	  "timezone": "America/Los_Angeles",
	  "agenda": "Talk about sex ed"
	}

	*/

}

func getCode(urlstring string) string {
	values, _ := url.ParseQuery(urlstring)

	code := values.Get("code")
	return code
}

func getAccessToken(code string) (*types.ZoomAccessToken, error) {

	str := "https://zoom.us/oauth/token?grant_type=authorization_code&code=" + code + "&redirect_uri=http://www.localhost:5000/api/auth/zoom/callback"
	req, err := http.NewRequest("POST", str, nil)
	if err != nil {
		return nil, err
	}

	head := "Basic cFdES0txY2dTcmE5eWhsRkNmNDNBdzpSd0FYaTc3UmVEZ2xkVUIwTG55RDJtcVFLb1BINTVNeg=="
	req.Header.Set("Authorization", head)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//result := make(map[string]interface{})
	zoomAccessToken := types.ZoomAccessToken{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &zoomAccessToken)
	return &zoomAccessToken, nil

}
