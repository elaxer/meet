package app

type ctxKey int

const (
	CtxKeyUser ctxKey = iota
	CtxKeyQuestionnaire
)

type UploadType string

const (
	UploadTypeImage UploadType = "image"
)

const LogFilename = "logs/log.log"
