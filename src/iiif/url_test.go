package iiif

import (
	"fmt"
	"strings"
	"testing"

	"github.com/uoregon-libraries/gopkg/assert"
)

var weirdID = "identifier-foo-bar%2Fbaz,,,,,chameleon"
var simplePath = weirdID + "/full/full/30/default.jpg"

func TestInvalid(t *testing.T) {
	badURL := strings.Replace(simplePath, "/full/full", "/bad/full", 1)
	badURL = strings.Replace(badURL, "default.jpg", "default.foo", 1)
	i, err := NewURL(badURL)
	assert.Equal("invalid region, invalid format", err.Error(), "NewURL error message", t)
	assert.False(i.Valid(), "IIIF URL is invalid", t)

	// All other data should still be extracted despite this being a bad IIIF URL
	assert.Equal(weirdID, i.ID.String(), "identifier should be extracted", t)
	assert.Equal("identifier-foo-bar/baz,,,,,chameleon", i.ID.Path(), "ID path", t)
	assert.Equal(RTNone, i.Region.Type, "bad Region is RTNone", t)
	assert.Equal(STFull, i.Size.Type, "Size is STFull", t)
	assert.Equal(30.0, i.Rotation.Degrees, "i.Rotation.Degrees", t)
	assert.True(!i.Rotation.Mirror, "!i.Rotation.Mirror", t)
	assert.Equal(QDefault, i.Quality, "i.Quality == QDefault", t)
	assert.Equal(FmtUnknown, i.Format, "i.Format == FmtJPG", t)
	assert.Equal(false, i.Info, "not an info request", t)
}

func TestValid(t *testing.T) {
	i, err := NewURL(simplePath)
	assert.NilError(err, "NewURL has no error", t)

	assert.True(i.Valid(), fmt.Sprintf("Expected %s to be valid", simplePath), t)
	assert.Equal(weirdID, i.ID.String(), "identifier should be extracted", t)
	assert.Equal("identifier-foo-bar/baz,,,,,chameleon", i.ID.Path(), "ID path", t)
	assert.Equal(RTFull, i.Region.Type, "Region is RTFull", t)
	assert.Equal(STFull, i.Size.Type, "Size is STFull", t)
	assert.Equal(30.0, i.Rotation.Degrees, "i.Rotation.Degrees", t)
	assert.True(!i.Rotation.Mirror, "!i.Rotation.Mirror", t)
	assert.Equal(QDefault, i.Quality, "i.Quality == QDefault", t)
	assert.Equal(FmtJPG, i.Format, "i.Format == FmtJPG", t)
	assert.Equal(false, i.Info, "not an info request", t)
}

func TestInfo(t *testing.T) {
	i, err := NewURL("some%2Fvalid%2Fpath.jp2/info.json")
	assert.NilError(err, "info request isn't an error", t)
	assert.Equal("some%2Fvalid%2Fpath.jp2", i.ID.String(), "identifier", t)
	assert.Equal(true, i.Info, "is an info request", t)
}

func TestInfoBaseRedirect(t *testing.T) {
	i, err := NewURL("some%2Fvalid%2Fpath.jp2")
	assert.NilError(err, "info request isn't an error", t)
	assert.Equal("some%2Fvalid%2Fpath.jp2", i.ID.String(), "identifier", t)
	assert.Equal(true, i.BaseURIRedirect, "is an info request", t)
}
