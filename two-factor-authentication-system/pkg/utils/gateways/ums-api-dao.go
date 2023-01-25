package gateways

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetUserEmailInfo(userId int64) string {
	req, err := http.NewRequest("GET", "http://3.36.21.109:10210/user-management-service-1.0/api/users/email", nil)
	if err != nil {
		log.Print(err)
	}
	q := req.URL.Query()
	q.Add("userId", strconv.Itoa(int(userId)))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		log.Panic("user email could not get from ums")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	sb := string(body)
	log.Printf("%v", sb)
	return sb
}
