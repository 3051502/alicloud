//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package images

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/zclconf/go-cty/cty"

	packerecs "github.com/hashicorp/packer-plugin-alicloud/builder/ecs"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type Datasource struct {
	config Config
}

type Config struct {
	common.PackerConfig            `mapstructure:",squash"`
	packerecs.AlicloudAccessConfig `mapstructure:",squash"`
	ImagesFilterOptions            `mapstructure:",squash"`

	ctx interpolate.Context
}

type DatasourceOutput struct {
	ImageId string `mapstructure:"imageId"`

	ImageName string `mapstructure:"image_name"`

	CreationTime string `mapstructure:"creation_time"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &d.config.ctx,
	}, raws...)

	if err != nil {
		return err
	}

	var errs *packersdk.MultiError
	errs = packersdk.MultiErrorAppend(errs, d.config.AlicloudAccessConfig.Prepare(&d.config.ctx)...)

	if d.config.Empty() {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("Both Region and Image Name must be provided"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	client, err := d.config.Client()

	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	image, err := d.config.ImagesFilterOptions.GetFilteredImage(&ecs.DescribeImagesRequest{}, client.Client)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	output := DatasourceOutput{
		ImageId:      image.ImageId,
		ImageName:    image.ImageName,
		CreationTime: image.CreationTime,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
