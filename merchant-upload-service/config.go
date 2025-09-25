package main

type Config struct {
	BaseUploadDir string `json:"baseUploadDir"`
	MaxFileSize   int64  `json:"maxFileSize"`
}
