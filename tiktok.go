package goget

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Tiktok struct {
	Cli              *http.Client
	ProxyUrl         string
	ProxyRedirectUrl string
}

func NewTiktok(opts ...Option) *Tiktok {
	options := &Options{
		Timeout:   time.Second * 60,
		Transport: http.DefaultTransport,
	}
	options.Apply(opts)

	cli := &http.Client{
		Timeout:   options.Timeout,
		Transport: options.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 1 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return &Tiktok{
		Cli:              cli,
		ProxyUrl:         options.ProxyUrl,
		ProxyRedirectUrl: options.ProxyRedirectUrl,
	}
}

func (c *Tiktok) GetTiktok(ctx context.Context, url string) (string, error) {
	var vid string
	_, err := strconv.ParseInt(url, 10, 64)
	if err == nil {
		vid = url
	} else {
		if strings.Contains(url, "www.tiktok.com/") {
			vid = DigitReg.FindString(url)
			if vid == "" {
				return "", errors.New("无效的地址")
			}
		} else {
			surl := UrlReg.FindString(url)
			if surl == "" {
				return "", errors.New("无效的地址")
			}

			videoId, er := c.GetTiktokVideoId(ctx, surl)
			if er != nil {
				return "", er
			}
			vid = videoId
		}
	}
	var apiUrl string
	if c.ProxyUrl != "" {
		apiUrl = fmt.Sprintf("%v/aweme/v1/feed/?aweme_id=%v", c.ProxyUrl, vid)
	} else {
		apiUrl = fmt.Sprintf("https://api16-normal-c-useast1a.tiktokv.com/aweme/v1/feed/?aweme_id=%v", vid)
	}
	ua := "com.ss.android.ugc.trill/494+Mozilla/5.0+(Linux;+Android+12;+2112123G+Build/SKQ1.211006.001;+wv)+AppleWebKit/537.36+(KHTML,+like+Gecko)+Version/4.0+Chrome/107.0.5304.105+Mobile+Safari/537.36"

	req, er := http.NewRequest("GET", apiUrl, nil)
	if er != nil {
		return "", er
	}

	req.Header.Set("User-Agent", ua)
	res, er := c.Cli.Do(req)
	if er != nil {
		return "", er
	}
	defer res.Body.Close()
	body, er := io.ReadAll(res.Body)
	if er != nil {
		return "", er
	}

	return string(body), nil
}

func (c *Tiktok) GetTiktokVideoId(ctx context.Context, surl string) (string, error) {
	if c.ProxyRedirectUrl != "" {
		surl = strings.Replace(surl, "https://vt.tiktok.com", c.ProxyRedirectUrl, -1)
	}
	res, er := c.Cli.Get(surl)
	if er != nil {
		return "", er
	}

	if res.StatusCode != 302 && res.StatusCode != 301 {
		return "", errors.New("获取视频ID失败")
	}

	locUrl, _ := res.Location()

	locStr := strings.Split(locUrl.String(), "?")[0]
	locStr = strings.Split(locStr, "/video/")[1]

	return locStr, nil
}
