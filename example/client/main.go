package main

import (
	"context"
	ktconf "github.com/aiagt/kitextool/conf"
	echo "github.com/aiagt/kitextool/example/kitex_gen/echo/echoservice"
	ktresolver "github.com/aiagt/kitextool/option/client/resolver"
	ktclient "github.com/aiagt/kitextool/suite/client"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

var conf = new(ktconf.MultiClientConf)

func init() {
	ktconf.LoadFiles(conf, "conf.yaml", "client/conf.yaml")
}

func main() {
	cli := echo.MustNewClient("echo",
		// Use the KitexTool suite in the client
		client.WithSuite(ktclient.NewKitexToolSuite(
			// Global configuration of KitexTool suite
			conf.GetClientConf("echo"),
			// Use consul for service discovery
			ktresolver.WithResolver(ktresolver.NewConsulResolver),
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
