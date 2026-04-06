package store

import (
	"maphoto/internal/util"
	"path/filepath"
)

var (
	PICPATH   string = filepath.Join(util.ExcutePath(), "maphoto_data", "photo")
	THUMBPATH string = filepath.Join(util.ExcutePath(), "maphoto_data", "thumbnail")
)
