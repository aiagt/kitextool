package main

import (
	"context"

	ktconf "github.com/ahaostudy/kitextool/conf"
	echo "github.com/ahaostudy/kitextool/example/kitex_gen/echo/echoservice"
	ktresolver "github.com/ahaostudy/kitextool/option/client/resolver"
	ktclient "github.com/ahaostudy/kitextool/suite/client"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

var conf = new(ktconf.ClientConf)

func init() {
	ktconf.LoadFiles(conf, "conf.yaml")
}

func main() {
	cli := echo.MustNewClient("echo",
		// Use the KitexTool suite in the client
		client.WithSuite(ktclient.NewKitexToolSuite(
			// Global configuration of KitexTool suite
			conf,
			// Use nacos for service discovery
			ktresolver.WithResolver(ktresolver.NewNacosResolver),
		)),
	)

	ctx := context.Background()
	req := "hello"
	resp, err := cli.Echo(ctx, req)
	if err != nil {
		klog.CtxErrorf(ctx, "%v", err.Error())
		return
	}
	klog.Infof("resp: %v", resp)
}
