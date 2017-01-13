package cmd

import "mime/multipart"

const indexfile = "index.yaml"
const logfile = "log"

type ChartFile struct {
	Data *multipart.FileHeader `form:"Data"`
}

type fileData struct {
	name string
	data string
}
