package goget

import (
	"context"
	"errors"
	"fmt"
	"github.com/eininst/flog"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type KsVideo struct {
	Cli *http.Client
}

func NewKsVideo(opts ...Option) *KsVideo {
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

	return &KsVideo{
		Cli: cli,
	}
}

func (c *KsVideo) GetVideoId(ctx context.Context, url string) (string, error) {
	var vid string
	_, err := strconv.ParseInt(url, 10, 64)
	if err == nil {
		vid = url
	} else {
		if strings.Contains(url, "v.kuaishou.com") ||
			strings.Contains(url, "https://www.kuaishou.com/f") {
			_videoUrl, er := c.getVideoId(ctx, url)
			if er != nil {
				return "", er
			}
			vid = _videoUrl
		}
	}
	if vid == "" {
		return "", errors.New("获取视频ID失败")
	}
	return vid, nil
}

func (c *KsVideo) GetVideo(ctx context.Context, url string) (map[string]string, error) {
	var videoUrl = url
	_, err := strconv.ParseInt(url, 10, 64)
	if err == nil {
		_videoUrl := fmt.Sprintf("https://www.kuaishou.com/short-video/%v", url)
		videoUrl = _videoUrl
	} else {
		if strings.Contains(url, "v.kuaishou.com") ||
			strings.Contains(url, "https://www.kuaishou.com/f") {
			_videoUrl, er := c.getVideoId(ctx, url)
			if er != nil {
				return nil, er
			}
			videoUrl = _videoUrl
		}
	}
	req, er := http.NewRequest("GET", videoUrl, nil)
	if er != nil {
		return nil, er
	}

	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Cookie", "kpf=PC_WEB; clientid=3; did=web_94ab9f62dd5b20fe38839602c22df142; userId=2753518500; kuaishou.server.web_st=ChZrdWFpc2hvdS5zZXJ2ZXIud2ViLnN0EqABiGCSATDTG77uOLvYG0_SJzaUi8lELcdKNZg1ASY51CMadxEW8zBV8smPNYdhPgaHHgnwbn560iJXDhDSERtyDpucTsPk_p5deUNsvw3Wx87G4Cu03hVvXb_uqM_kQyjE4OSAaxiMrDcC4YUDeWfWsW6BIM8nuII9ABQcNcQqtJdxDw00QEd3I3KDN2XvyWWKF4RfQoMDDINy-Nb-DfQIOhoS6uws2LN-siMyPVYdMaXTUH7FIiCHjwFNfyY7Y_HiRFuN9Jq_hB02oahp2dThKwUbIvccgigFMAE; kuaishou.server.web_ph=b42b72f8f03da4ad0284b7f2b8c6d5f5b5b6; kpn=KUAISHOU_VISION; _dd_s=rum=0&expire=1701776357360")

	res, er := c.Cli.Do(req)
	if er != nil {
		return nil, er
	}
	defer res.Body.Close()
	body, er := io.ReadAll(res.Body)
	if er != nil {
		return nil, er
	}

	bodyStr := string(body)

	caption := regexp.MustCompile(`"caption":"(.*?),`).FindString(bodyStr)
	duration := regexp.MustCompile(`"duration":(.*?),`).FindString(bodyStr)
	coverUrl := regexp.MustCompile(`"coverUrl":"(.*?)"`).FindString(bodyStr)
	photoUrl := regexp.MustCompile(`"photoUrl":"(.*?)"`).FindString(bodyStr)

	caption = strings.Replace(caption, `"caption":"`, "", -1)
	caption = strings.Replace(caption, `",`, "", -1)

	duration = strings.Replace(duration, `"duration":`, "", -1)
	duration = strings.Replace(duration, `,`, "", -1)

	m := map[string]string{
		"caption":  caption,
		"duration": duration,
		"coverUrl": covertUrl(coverUrl, `"coverUrl":`),
		"photoUrl": covertUrl(photoUrl, `"photoUrl":`),
	}

	return m, nil
}

func covertUrl(url string, key string) string {
	rr := strings.Replace(url, key, "", -1)
	rr = strings.Replace(rr, `"`, "", -1)
	results, _ := strconv.Unquote("\"" + rr + "\"")

	flog.Info(results)
	return results
}

func (c *KsVideo) getVideoId(ctx context.Context, surl string) (string, error) {
	res, er := c.Cli.Get(surl)
	if er != nil {
		return "", er
	}
	if res.StatusCode != 302 {
		return "", errors.New("获取视频ID失败")
	}
	locUrl, _ := res.Location()

	return locUrl.String(), nil
}
