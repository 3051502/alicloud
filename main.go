package main

import (
	"fmt"
	"os"

	imageFilter "github.com/3051502/packer-plugin-alicloud/datasource/images"

	ecsbuilder "github.com/hashicorp/packer-plugin-alicloud/builder/ecs"
	importpp "github.com/hashicorp/packer-plugin-alicloud/post-processor/alicloud-import"
	version "github.com/hashicorp/packer-plugin-alicloud/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("ecs", new(ecsbuilder.Builder))
	pps.RegisterPostProcessor("import", new(importpp.PostProcessor))
	pps.RegisterDatasource("images", new(imageFilter.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
