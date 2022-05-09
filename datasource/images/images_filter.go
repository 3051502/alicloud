//go:generate packer-sdc struct-markdown
package images

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type ImagesFilterOptions struct {

	// [DescribeImages](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/describeimages)

	//the region that the image is in
	Region string `mapstructure:"region"`

	//the name of the image to search for
	ImageName string `mapstructure:"image_name"`
}

func (d *ImagesFilterOptions) Empty() bool {
	return len(d.Region) == 0 || len(d.ImageName) == 0
}

func (d *ImagesFilterOptions) GetFilteredImage(params *ecs.DescribeImagesRequest, ecsconn *ecs.Client) (*ecs.Image, error) {

	request := ecs.CreateDescribeImagesRequest()
	request.Scheme = "https"
	request.ImageName = params.ImageName
	request.RegionId = params.RegionId

	response, err := ecsconn.DescribeImages(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	if response.TotalCount == 0 {
		err := fmt.Errorf("No Image was found matching filters: %v", params)
		return nil, err
	}

	fmt.Printf("response is %#v\n", response)

	var image *ecs.Image

	image.Description = "testing"

	return image, nil
}
