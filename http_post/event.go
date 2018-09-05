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

// PostEvent triggers an http POST message to be consumed by RestAPI.
// This is a testclass to test the REST services.
func main() {
	cmdline.AppName = "Produce Http Event for Image Intake Service"
	cmdline.AppCmdName = "produce_http_event"
	cmdline.AppVersion = "1.0"
	cmdline.CopyrightYears = "2018"
	cmdline.CopyrightHolder = "Fanatics, Inc."
	cl := cmdline.New(true)
	cl.Description = "Generate Test Post event"

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

	requestURL := "http://localhost:8080/upload"

	logger := argos.NewLogger(argos.DefaultConfigByEnv(
		cmdline.AppCmdName, service.VersionForLog(), cfg.DeploymentEnv, cfg.InstanceName))
	span := logger.NewRootSpan("startup", nil)
	defer span.End()
	span.Debugf("Using the '%s' environment", env)

	productInfo := make(map[string]string)
	productInfo["ProductId"] = "null"
	productInfo["AltIndex"] = "37970"
	productInfo["IsAlt"] = "false"

	productData := make(map[string]interface{})
	productData["AnselImageType"] = "product"
	productData["IsNewFile"] = "false"
	productData["FileAction"] = "update"
	productData["SourceFilePath"] = "productimages/_37000/G:\\Services\\Fanatics.ImageUploadAPI\\jax-img-cls010.ff.p10\\images\\ProductImages\\_37000\\ff_37970_full.jpg"
	productData["SourceFileHost"] = "images.footballfanatics.com" // "192.168.125.99"
	productData["DirectUploadS3KeyName"] = "wwwroot/images/productimages/_2991000/ff_2991461_full.jpg"
	productData["StagedImageFileID"] = "https://"
	productData["OriginalOwner"] = "null"
	productData["ProductImageInfo"] = productInfo

	maindata := make(map[string]interface{})
	maindata["AnselImage"] = productData

	b := new(bytes.Buffer)
	ue := json.NewEncoder(b).Encode(maindata)
	if ue != nil {
		logger.Error("Error Encoding....", ue)
	}
	span.Debugf("Sending post data '%v'", maindata)
	res, err := http.Post(requestURL, "application/json; charset=utf-8", b)
	if err != nil {
		span.Debugf("Errored sending post request", err)
	}
	span.Debugf("Response sent --> '%v'", res)
}
