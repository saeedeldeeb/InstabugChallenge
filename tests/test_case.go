package tests

import (
	"github.com/goravel/framework/testing"

	"chat/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
