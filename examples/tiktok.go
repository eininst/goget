package main

import (
	"context"
	"github.com/eininst/flog"
	"github.com/eininst/goget"
)

func main() {
	tiktok := goget.NewTiktok(
		goget.WithProxyUrl("http://xxxx"),
		goget.WithProxyRedirectUrl("http://xxxxx/vtk"))

	s, x := tiktok.GetTiktok(context.TODO(), "https://vt.tiktok.com/ZSFrcG2Q5/")
	flog.Info(x)
	flog.Info(s)

}
