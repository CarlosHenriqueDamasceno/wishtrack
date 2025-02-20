package test

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/stretchr/testify/suite"
)

func PrepareBody(body any, suite *suite.Suite) io.Reader {
	b, err := json.Marshal(body)
	suite.Assert().Nil(err, "Fail to marshal payload")
	return bytes.NewReader(b)
}
