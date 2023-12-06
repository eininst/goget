package goget

import "context"

var xg = NewXiguaVideo()
var ks = NewKsVideo()
var dy = NewDouyinVideo()
var bilibili = NewBiliBiliVideo()

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
