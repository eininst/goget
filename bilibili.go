package goget

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
	"time"
)

type BiliBiliVideo struct {
	Cli *http.Client
}

func NewBiliBiliVideo(opts ...Option) *BiliBiliVideo {
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

	return &BiliBiliVideo{
		Cli: cli,
	}
}

func (c *BiliBiliVideo) GetVideoId(ctx context.Context, url string) (string, error) {
	videoId := url
	if strings.Contains(url, "https://www.bilibili.com") {
		xx := UrlReg.FindString(url)
		xx = strings.Replace(xx, "https://www.bilibili.com/video/", "", -1)
		lindex := strings.Index(xx, "/")
		if lindex == -1 {
			lindex = len(xx)
		}
		videoId = xx[0:lindex]
	}

	videoId = strings.Replace(videoId, "video/BV", "", -1)
	if strings.Contains(videoId, "video/av") {
		videoId = strings.Replace(videoId, "video/av", "", -1)
	}

	return videoId, nil
}

func (c *BiliBiliVideo) GetVideo(ctx context.Context, url string) (map[string]any, error) {
	videoId := url
	if strings.Contains(url, "https://www.bilibili.com") {
		xx := UrlReg.FindString(url)
		xx = strings.Replace(xx, "https://www.bilibili.com/video/", "", -1)
		lindex := strings.Index(xx, "/")
		if lindex == -1 {
			lindex = len(xx)
		}
		videoId = xx[0:lindex]
	}

	videoId = strings.Replace(videoId, "video/BV", "", -1)
	apiUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%v", videoId)

	if strings.Contains(videoId, "video/av") {
		videoId = strings.Replace(videoId, "video/av", "", -1)
		apiUrl = fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?aid=%v", videoId)
	}

	res, er := c.Cli.Get(apiUrl)
	if er != nil {
		return nil, er
	}
	defer res.Body.Close()
	body, er := io.ReadAll(res.Body)
	if er != nil {
		return nil, er
	}

	asd := gjson.Parse(string(body))
	data := asd.Get("data")

	aid := data.Get("aid").String()
	cid := data.Get("cid").String()

	playUrlApi := fmt.Sprintf("https://api.bilibili.com/x/player/playurl?avid=%v&cid=%v&platform=html5", aid, cid)

	pres, per := c.Cli.Get(playUrlApi)
	if per != nil {
		return nil, per
	}

	defer pres.Body.Close()
	pbody, er := io.ReadAll(pres.Body)
	if er != nil {
		return nil, er
	}
	playData := gjson.Parse(string(pbody)).Get("data")

	durl := playData.Get("durl").Array()[0]
	playUrl := durl.Get("url").String()
	playSize := durl.Get("size").Int()
	playTime := durl.Get("length").Int() / 1000

	result := map[string]any{
		"playUrl":  playUrl,
		"playSize": playSize,
		"playTime": playTime,
		"tname":    data.Get("tname").String(),
		"pic":      data.Get("pic").String(),
		"title":    data.Get("title").String(),
		"desc":     data.Get("desc"),
		"owner":    data.Get("owner").Raw,
	}

	return result, nil
}
