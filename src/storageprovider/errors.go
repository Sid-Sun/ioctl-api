package storageprovider

import "errors"

var ErrNotFound = errors.New("object not found")
var ErrAlreadyExists = errors.New("object already exists")
