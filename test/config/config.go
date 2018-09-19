package config

import "os"

func init() {
	os.Setenv("AliLogEndpoint", "<endpoint>")
	os.Setenv("AliLogAccessKeyID", "<access_key>")
	os.Setenv("AliLogAccessKeySecret", "<access_secret>")
	os.Setenv("AliLogName", "<log_project>")
	os.Setenv("AliLogStoreName", "<log_store>")
	os.Setenv("AliLogFile", "./json.log")
}