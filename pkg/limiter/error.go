package limiter

import "github.com/pkg/errors"

var ErrExceedMaxCount = errors.New("exceed max count of limiter")
