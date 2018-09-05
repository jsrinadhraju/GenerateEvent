package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/richardwilkes/atexit"
	"github.com/richardwilkes/cmdline"

	"jaxf-github.fanatics.corp/forge/furnace"
	"jaxf-github.fanatics.corp/forge/furnace/fanflow"
	"jaxf-github.fanatics.corp/forge/furnace/log/argos"
	"jaxf-github.fanatics.corp/forge/furnace/service"
)

// FanflowEvent triggers a kafka fanflow event to be consumed by Fanflow Consumer.
func main() {
	cmdline.AppName = "Image Upload Service"
	cmdline.AppCmdName = "image_upload_service" // override whatever might have actually been used
	cmdline.AppVersion = "1.0"
	cmdline.CopyrightYears = "2018"
	cmdline.CopyrightHolder = "Fanatics, Inc."
	cl := cmdline.New(true)
	cl.Description = "Generate Test Fanflow event"

	env := furnace.EnforceEnv(os.Getenv("ENV"), false)
	cl.NewStringOption(&env).SetSingle('e').SetName("env").SetUsage("The environment to run in. May be " + furnace.AvailableOptions(true))
	var configPath string
	cl.NewStringOption(&configPath).SetSingle('c').SetName("config").SetArg("path").SetUsage("The path to a configuration file")
	cl.Parse(os.Args[1:])

	cfg := service.NewConfigByEnv(furnace.DeploymentEnv{
		Realm:  furnace.Ecom,
		Env:    env,
		Region: furnace.AWSUSEast1,
	})

	logger := argos.NewLogger(argos.DefaultConfigByEnv(
		cmdline.AppCmdName, service.VersionForLog(), cfg.DeploymentEnv, cfg.InstanceName))
	span := logger.NewRootSpan("startup", nil)
	defer span.End()
	span.Debugf("Using the '%s' environment", env)

	ff := fanflow.New(fanflow.Env("master", cfg.Env))
	productInfo := make(map[string]string)
	productInfo["ProductId"] = "null"
	productInfo["AltIndex"] = "37970"
	productInfo["IsAlt"] = "false"

	productData := make(map[string]interface{})
	productData["AnselImageType"] = "product"
	productData["IsNewFile"] = "false"
	productData["FileAction"] = "update"
	productData["SourceFilePath"] = "productimages/_37000/G:\\Services\\Fanatics.ImageUploadAPI\\jax-img-cls010.ff.p10\\images\\ProductImages\\_37000\\ff_37970_full.jpg"
	productData["SourceFileHost"] = "192.168.125.99" // "images.footballfanatics.com"
	productData["DirectUploadS3KeyName"] = "wwwroot/images/productimages/_2991000/ff_2991461_full.jpg"
	productData["StagedImageFileID"] = "https://"
	productData["OriginalOwner"] = "null"
	productData["ProductImageInfo"] = productInfo

	maindata := make(map[string]interface{})
	maindata["AnselImage"] = productData

	ff.Send(fanflow.Event{
		Type:    "UPLOAD_IMAGE.ANSEL",
		Payload: maindata,
	})
	// span.Debug("Sent the event successfully -->", maindata)
	ff.Shutdown()
	atexit.Exit(0)
}
