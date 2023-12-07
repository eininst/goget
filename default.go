package goget

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

var xg = NewXiguaVideo()
var ks = NewKsVideo()
var dy = NewDouyinVideo()
var bilibili = NewBiliBiliVideo()

type Video struct {
	Kind string
	Data string
}

func GetDouyinVideo(ctx context.Context, url string, sessionidss string) (string, error) {
	return dy.GetVideo(ctx, url, sessionidss)
}

func GetKsVideo(ctx context.Context, url string) (map[string]string, error) {
	return ks.GetVideo(ctx, url)
}

func GetXgVideo(ctx context.Context, url string) (string, error) {
	return xg.GetVideo(ctx, url)
}

func GetBilibiliVideo(ctx context.Context, url string) (map[string]any, error) {
	return bilibili.GetVideo(ctx, url)
}

func GetVideo(ctx context.Context, url string, sessionidss string) (*Video, error) {
	if strings.Contains(url, "douyin") {
		data, er := dy.GetVideo(ctx, url, sessionidss)
		if er != nil {
			return nil, er
		}
		return &Video{
			Kind: "dy",
			Data: data,
		}, nil
	} else if strings.Contains(url, "kuaishou") {
		data, er := ks.GetVideo(ctx, url)
		if er != nil {
			return nil, er
		}
		dataStr, _ := json.Marshal(data)
		return &Video{
			Kind: "ks",
			Data: string(dataStr),
		}, nil
	} else if strings.Contains(url, "ixigua") {
		data, er := xg.GetVideo(ctx, url)
		if er != nil {
			return nil, er
		}
		return &Video{
			Kind: "xg",
			Data: data,
		}, nil
	} else if strings.Contains(url, "bilibili") {
		data, er := bilibili.GetVideo(ctx, url)
		if er != nil {
			return nil, er
		}
		dataStr, _ := json.Marshal(data)
		return &Video{
			Kind: "bilibili",
			Data: string(dataStr),
		}, nil
	}

	return nil, errors.New("无效的地址")
}
