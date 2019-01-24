package synse

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet_NoVarsSet(t *testing.T) {
	v := Get()

	assert.NotNil(t, v)
	assert.Empty(t, v.BuildDate)
	assert.Empty(t, v.GitCommit)
	assert.Empty(t, v.GitTag)
	assert.Empty(t, v.GoVersion)
	assert.Empty(t, v.BinVersion)
	assert.Equal(t, runtime.GOARCH, v.Arch)
	assert.Equal(t, runtime.GOOS, v.OS)
}

func TestGet_VarsSet(t *testing.T) {
	BuildDate = "foo"
	GitCommit = "bar"
	GitTag = "foo"
	GoVersion = "bar"
	BinVersion = "foo"

	v := Get()

	assert.NotNil(t, v)
	assert.Equal(t, "foo", v.BuildDate)
	assert.Equal(t, "bar", v.GitCommit)
	assert.Equal(t, "foo", v.GitTag)
	assert.Equal(t, "bar", v.GoVersion)
	assert.Equal(t, "foo", v.BinVersion)
	assert.Equal(t, runtime.GOARCH, v.Arch)
	assert.Equal(t, runtime.GOOS, v.OS)
}
