package bublyk

//nolint: goimports // unknown problem
import (
	"errors"
	"github.com/kaatinga/const-errs"
)

const ErrUnrecognizedFormat const_errs.Error = "unknown date format"

var errUndefined = errors.New("cannot encode status undefined")
