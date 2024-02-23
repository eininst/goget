package main

import (
	"context"
	"github.com/eininst/flog"
	"github.com/eininst/goget"
)

func main() {
	tiktok := goget.NewTiktok(
		goget.WithProxyUrl("http://hk.youpinsay.com/tk"),
		goget.WithProxyRedirectUrl("http://hk.youpinsay.com/vtk"))
	//body, er := tiktok.GetTiktok(context.TODO(), "7338615050246114603", "")
	//flog.Info(er)
	//flog.Info(body)

	//https://www.tiktok.com/@openai/video/7338615050246114603
	//https://vt.tiktok.com/ZSFrcG2Q5/

	//vt.tiktok.com
	s, x := tiktok.GetTiktok(context.TODO(), "https://vt.tiktok.com/ZSFrcG2Q5/")
	flog.Info(x)
	flog.Info(s)

	//goget.NewDouyinVideo().GetDouyinVideoId(context.TODO(), "https://v.douyin.com/iNVBE1WV/")
}
