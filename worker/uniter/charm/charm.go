// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	"errors"

	"github.com/juju/charm/v8"
	"github.com/juju/collections/set"
	"github.com/juju/utils/v3"
)

// Logger is here to stop the desire of creating a package level Logger.
// Don't do this, instead pass one in as needed.
type logger interface{}

var _ logger = struct{}{}

// CharmURLPath is the path within a charm directory to which Deployers
// commonly write the charm URL of the latest deployed charm.
const CharmURLPath = ".juju-charm"

// Bundle allows access to a charm's files.
type Bundle interface {

	// ArchiveMembers returns a set of slash-separated strings representing files,
	// directories, and symlinks stored in the bundle.
	ArchiveMembers() (set.Strings, error)

	// ExpandTo unpacks the entities referenced in the manifest into the
	// supplied directory. If it returns without error, every file referenced
	// in the charm must be present in the directory; implementations may vary
	// in the details of what they do with other files present.
	ExpandTo(dir string) error
}

// BundleInfo describes a Bundle.
type BundleInfo interface {

	// URL returns the charm URL identifying the bundle.
	URL() *charm.URL

	// ArchiveSha256 returns the hex-encoded SHA-256 digest of the bundle data.
	ArchiveSha256() (string, error)
}

// BundleReader provides a mechanism for getting a Bundle from a BundleInfo.
type BundleReader interface {

	// Read returns the bundle identified by the supplied info. The abort chan
	// can be used to notify an implementation that it need not complete the
	// operation, and can immediately error out if it is convenient to do so.
	Read(bi BundleInfo, abort <-chan struct{}) (Bundle, error)
}

// Deployer is responsible for installing and upgrading charms.
type Deployer interface {

	// Stage must be called to prime the Deployer to install or upgrade the
	// bundle identified by the supplied info. The abort chan can be used to
	// notify an implementation that it need not complete the operation, and
	// can immediately error out if it convenient to do so. It must always
	// be safe to restage the same bundle, or to stage a new bundle.
	Stage(info BundleInfo, abort <-chan struct{}) error

	// Deploy will install or upgrade the most recently staged bundle.
	// Behaviour is undefined if Stage has not been called. Failures that
	// can be resolved by user intervention will be signalled by returning
	// ErrConflict.
	Deploy() error
}

// Logger represents the logging functions used by the charm package.
type Logger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

// ErrConflict indicates that an upgrade failed and cannot be resolved
// without human intervention.
var ErrConflict = errors.New("charm upgrade has conflicts")

// ReadCharmURL reads a charm identity file from the supplied path.
func ReadCharmURL(path string) (*charm.URL, error) {
	surl := ""
	if err := utils.ReadYaml(path, &surl); err != nil {
		return nil, err
	}
	return charm.ParseURL(surl)
}

// WriteCharmURL writes a charm identity file into the supplied path.
func WriteCharmURL(path string, url *charm.URL) error {
	return utils.WriteYaml(path, url.String())
}
