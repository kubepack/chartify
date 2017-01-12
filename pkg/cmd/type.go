package cmd

import "mime/multipart"

type ChartFile struct {
	Data *multipart.FileHeader `form:"Data"`
}
