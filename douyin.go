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

type DouyinVideo struct {
	Cli *http.Client
}

func NewDouyinVideo(opts ...Option) *DouyinVideo {
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

	return &DouyinVideo{
		Cli: cli,
	}
}

func (c *DouyinVideo) GetVideoId(ctx context.Context, url string) (string, error) {
	var vid string
	_, err := strconv.ParseInt(url, 10, 64)
	if err == nil {
		vid = url
	} else {
		if strings.Contains(url, "www.douyin.com/video") {
			vid = DigitReg.FindString(url)
			if vid == "" {
				return "", errors.New("无效的地址")
			}
		} else {
			surl := UrlReg.FindString(url)
			if surl == "" {
				return "", errors.New("无效的地址")
			}

			videoId, er := c.getDouyinVideoId(ctx, surl)
			if er != nil {
				return "", er
			}
			vid = videoId
		}
	}

	return vid, nil
}

func (c *DouyinVideo) GetVideo(ctx context.Context, url string, sessionidss string) (string, error) {
	var vid string
	_, err := strconv.ParseInt(url, 10, 64)
	if err == nil {
		vid = url
	} else {
		if strings.Contains(url, "www.douyin.com/video") {
			vid = DigitReg.FindString(url)
			if vid == "" {
				return "", errors.New("无效的地址")
			}
		} else {
			surl := UrlReg.FindString(url)
			if surl == "" {
				return "", errors.New("无效的地址")
			}

			videoId, er := c.getDouyinVideoId(ctx, surl)
			if er != nil {
				return "", er
			}
			vid = videoId
		}
	}

	query := fmt.Sprintf("device_platform=webapp&aid=6383&channel=channel_pc_web&aweme_id=%v&pc_client_type=1&version_code=190500&version_name=19.5.0&cookie_enabled=true&screen_width=1344&screen_height=756&browser_language=zh-CN&browser_platform=Win32&browser_name=Firefox&browser_version=118.0&browser_online=true&engine_name=Gecko&engine_version=109.0&os_name=Windows&os_version=10&cpu_core_num=16&device_memory=&platform=PC&webid=7284189800734082615&msToken=B1N9FM825TkvFbayDsDvZxM8r5suLrsfQbC93TciS0O9Iii8iJpAPd__FM2rpLUJi5xtMencSXLeNn8xmOS9q7bP0CUsrt9oVTL08YXLPRzZm0dHKLc9PGRlyEk=",
		vid)
	apiUrl := fmt.Sprintf("https://www.douyin.com/aweme/v1/web/aweme/detail/?%v", query)
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	req, er := http.NewRequest("GET", apiUrl, nil)
	if er != nil {
		return "", er
	}

	req.Header.Set("User-Agent", ua)
	req.Header.Set("Referer", fmt.Sprintf("https://www.douyin.com/video/%v", vid))
	req.Header.Set("Cookie", fmt.Sprintf("sessionid_ss=%v", sessionidss))

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

func (c *DouyinVideo) getDouyinVideoId(ctx context.Context, surl string) (string, error) {
	res, er := c.Cli.Get(surl)
	if er != nil {
		return "", er
	}
	if res.StatusCode != 302 {
		return "", errors.New("获取视频ID失败")
	}
	locUrl, _ := res.Location()
	result := DigitReg.FindString(locUrl.String())
	if result == "" {
		return "", errors.New("解析参数失败 ->" + locUrl.String())
	}
	return result, nil
}
