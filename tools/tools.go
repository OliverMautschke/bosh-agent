// +build tools

package tools

import (
	_ "github.com/onsi/ginkgo/ginkgo" // comment to make golint happy
	_ "golang.org/x/lint/golint" // comment to make golint happy
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
