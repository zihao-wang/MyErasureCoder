package ErasureCoder

import "github.com/pkg/errors"

var ErrShortData = errors.New("not enough data to fill the number of requested shards")