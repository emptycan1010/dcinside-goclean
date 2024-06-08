package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func main() {
	fileread, err := os.ReadFile("ENV.json")
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(string(fileread))
	var useracc USERACC
	err = json.Unmarshal(fileread, &useracc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(useracc.Id)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	login(client, useracc)
	Getglist(client, useracc)
}

type USERACC struct {
	Id string `json:"id"`
	Pw string `json:"pw"`
}

func login(client *http.Client, useracc USERACC) {
	// https://msign.dcinside.com/login
	//body := fmt.Sprintf("code=%s&password=%s&loginCash=on&conKey=31a18171b48560f33de68eed&r_url=https%3A%2F%2Fm.dcinside.com&_token=", useracc.Id, useracc.Pw)
	value := url.Values{
		"code":       {useracc.Id},
		"password":   {useracc.Pw},
		"loginCash":  {"on"},
		"conKey":     {"31a18171b48560f33de68eed"},
		"r_url":      {"https://m.dcinside.com"},
		"_token":     {""},
		"return_url": {""},
	}
	req, err := http.NewRequest(
		"POST",
		"https://msign.dcinside.com/login",
		strings.NewReader(value.Encode()),
	)
	if err != nil {
		log.Fatalln(err)

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0")
	req.Header.Set("Referer", "https://msign.dcinside.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Host", "msign.dcinside.com")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	// print body as string
	//bt, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	return
}

func Getglist(client *http.Client, useracc USERACC) {
	// g_id=emptycan1010&menu=G_all&page=2&list_more=1
	bdy := url.Values{
		"g_id":      {useracc.Id},
		"menu":      {"G_all"},
		"page":      {"1"},
		"list_more": {"1"},
	}
	req, err := http.NewRequest(
		"POST",
		"https://m.dcinside.com/ajax/response-galloglist",
		strings.NewReader(bdy.Encode()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "m.dcinside.com")
	req.Header.Set("Origin", "https://m.dcinside.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://m.dcinside.com/gallog/emptycan1010?menu=G_all&s_menu=N")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	// print body as string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var glist GALLOGRESP
	err = json.Unmarshal(body, &glist)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

type HEADTEXT struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"` // G == 정갤, E == 마갤, N == 미갤
}

type GList struct {
	Subject       string `json:"subject"`
	No            string `json:"no"`
	Name          string `json:"name"`
	Pno           string `json:"pno"`
	GallCode      string `json:"gall_code"`
	Check_Comment string `json:"check_comment"`
	TotalComment  string `json:"total_comment"`
}

type GALLOGRESP struct {
	gallog_info struct {
		cate_set struct {
			code  string `json:"code"`
			level string `json:"level"`
			name  string `json:"name"`
			no    string `json:"no"`
			Type  string `json:"type"`
			num   string `json:"num"`
		} `json:"cate_set"`
		DcbestCnt     int        `json:"dcbest_cnt"`
		Lastupdate    string     `json:"lastupdate"`
		ManagerSkill  bool       `json:"manager_skill"`
		Profile_image string     `json:"profile_image"`
		Total_cnt     int        `json:"total_cnt"`
		HeadText      []HEADTEXT `json:"head_text"`
	} `json:"gallog_info"`
	gallog_list struct {
		Last_page int     `json:"last_page"`
		Data      []GList `json:"data"`
	} `json:"gallog_list"`
}
