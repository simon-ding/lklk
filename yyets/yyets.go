package yyets

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	FeedURL     = "http://rss.rrys.tv/rss/feed/"
	SearchURL   = "http://pc.zmzapi.com/index.php?g=api/pv2&m=index&a=search&accesskey=519f9cab85c8059d17544947k361a827&limit=200&k="
	DetailURL   = "http://pc.zmzapi.com/index.php?g=api/pv2&m=index&a=resource&accesskey=519f9cab85c8059d17544947k361a827&id="
	LoginURL    = "http://www.zmz2019.com/User/Login/ajaxLogin"
	FavURL      = "http://www.zmz2019.com/user/fav"
	ResourceURL = "http://www.zmz2019.com/resource/"
	UserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"
)


type Client struct {
	username string
	password string
	cookies  []*http.Cookie
}

func (c *Client) SetLogin(username, password string) {
	c.username = username
	c.password = password
}

func (c *Client) UserFavs() ([]string, error) {
	if err := c.login(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", FavURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)

	for _, cc := range c.cookies {
		req.AddCookie(cc)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %s", resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var res []string
	doc.Find(".user-favlist .fl-img ").Each(func(i int, selection *goquery.Selection) {
		href := selection.Find("a").AttrOr("href", "")
		splits := strings.Split(href, "/")
		resourceID := splits[len(splits)-1]
		res = append(res, resourceID)
	})

	return res, nil
}

func (c *Client) login() error {
	if c.username == "" || c.password == "" {
		return fmt.Errorf("username and password is needed")
	}
	v := url.Values{}
	v.Add("account", c.username)
	v.Add("password", c.password)
	v.Add("remember", "1")
	v.Add("url_back", FavURL)
	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var r struct {
		Status int    `json:"status"`
		Info   string `json:"info"`
	}
	json.Unmarshal(data, &r)

	if resp.StatusCode != 200 || r.Status != 1 {
		return fmt.Errorf("http status code: %s, response %v", resp.Status, r)
	}
	c.cookies = resp.Cookies()[2:]
	return nil
}


type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Description string  `xml:"description"`
	Link        string  `xml:"link"`
	Title       string  `xml:"title"`
	Item        []*Item `xml:"item"`
}

type Item struct {
	Link          string    `xml:"link"`
	Title         string    `xml:"title"`
	Guid          string    `xml:"guid"`
	PubDate       string    `xml:"pubDate"`
	DateFormatted time.Time `xml:"date"`
	Magnet        string    `xml:"magnet"`
	Ed2k          string    `xml:"ed2k"`
}

func (c *Client) ParseRssURL(url string) (*Feed, error) {
	MaxRetries := 5
	var data []byte
	for i := 0; i < MaxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			data, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			break
		}
		resp.Body.Close()
	}
	var feed Feed
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	for _, item := range feed.Channel.Item {
		item.DateFormatted, _ = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
	}
	return &feed, err
}

type Detail struct {
	ID            string `json:"id"`
	Cnname        string `json:"cnname"`
	Enname        string `json:"enname"`
	Channel       string `json:"channel"`
	ChannelCN     string `json:"channel_cn"`
	Category      string `json:"category"`
	CloseResource string `json:"close_resource"`
	PlayStatus    string `json:"play_status"`
	Poster        string `json:"poster"`
	URL           string `json:"url"`
}

func (c *Client) GetDetail(resourceID string) (*Detail, error) {
	var res struct {
		Data struct {
			Detail Detail `json:"detail"`
		} `json:"data"`
	}
	MaxRetries := 5
	var data []byte
	url := DetailURL + resourceID
	for i := 0; i < MaxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			data, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			break
		} else if err == nil {
			resp.Body.Close()
		}
	}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	res.Data.Detail.URL = ResourceURL + resourceID
	return &res.Data.Detail, nil
}
