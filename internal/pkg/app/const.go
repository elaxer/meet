package app

import (
	"path/filepath"
	"runtime"

	"github.com/huandu/go-sqlbuilder"
)

const (
	ListLimitDefault = 10
	ListLimitMax     = 100
)

const SQLBuilderFlavor = sqlbuilder.PostgreSQL

var (
	_, b, _, _ = runtime.Caller(0)
	RootDir, _ = filepath.Abs(filepath.Dir(b) + "/../..")
)
