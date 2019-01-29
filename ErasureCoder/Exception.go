package ErasureCoder

import "errors"

var ErrShortData = errors.New("not enough data to fill the number of requested shards")
var ErrNotMatchNumDataBlocks = errors.New("# datablocks does not match")
var ErrNoSufficientBlocks = errors.New("No sufficient blocks for recover")
