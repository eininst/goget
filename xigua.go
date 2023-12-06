package goget

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type XiguaVideo struct {
	Cli1 *http.Client
	Cli2 *http.Client
}

func NewXiguaVideo(opts ...Option) *XiguaVideo {
	options := &Options{
		Timeout:   time.Second * 60,
		Transport: http.DefaultTransport,
	}
	options.Apply(opts)

	cli1 := &http.Client{
		Timeout:   options.Timeout,
		Transport: options.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 1 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	cli2 := &http.Client{
		Timeout:   options.Timeout,
		Transport: options.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 2 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return &XiguaVideo{
		Cli1: cli1,
		Cli2: cli2,
	}
}

func (c *XiguaVideo) GetVideo(ctx context.Context, url string) (string, error) {
	videoId := url
	_, err := strconv.ParseInt(url, 10, 64)
	if err != nil {
		if strings.Contains(url, "v.ixigua.com") {
			res, er := c.Cli1.Get(url)
			if er != nil {
				return "", er
			}
			locationUrl, er := res.Location()
			if er != nil {
				return "", er
			}
			vidUrl := strings.Split(locationUrl.String(), "?")[0]
			_vid := DigitReg.FindString(vidUrl)
			if _vid == "" {
				return "", errors.New("无法解析该地址")
			}
			videoId = _vid
		} else {
			_vid := UrlReg.FindString(url)
			if _vid == "" {
				return "", errors.New("无法解析该地址")
			}
			videoId = _vid
		}
	}
	videoUrl := fmt.Sprintf("https://m.ixigua.com/video/%v?wid_try=1", videoId)
	reg := regexp.MustCompile(`\"vid\":\"([^\"]+)+`)

	vres, ver := c.Cli2.Get(videoUrl)
	if ver != nil {
		return "", ver
	}
	defer vres.Body.Close()
	vbody, er := io.ReadAll(vres.Body)
	dd := reg.FindString(string(vbody))
	vid := strings.Replace(dd, `"vid":"`, "", -1)

	r := time.Now().UnixNano() / 1000
	urlpart := fmt.Sprintf("/video/urls/v/1/toutiao/mp4/%v?r=%v", vid, r)

	originalData := []byte(urlpart)
	crc32.Checksum(originalData, crc32.IEEETable)
	i3 := crc32.ChecksumIEEE(originalData)

	jsonurl := fmt.Sprintf("https://ib.365yg.com%v&s=%v&nobase64=true", urlpart, i3)
	res, er := c.Cli2.Get(jsonurl)
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
