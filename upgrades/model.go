// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrades

import (
	"github.com/juju/errors"
	"github.com/juju/version/v2"
)

// MinMajorUpgradeVersion defines the minimum version all models
// must be running before a major version upgrade.
var MinMajorUpgradeVersion = map[int]version.Number{
	3: version.MustParse("2.8.9"),
}

// UpgradeAllowed returns true if a major version upgrade is allowed
// for the specified from and to versions.
func UpgradeAllowed(from, to version.Number) (bool, version.Number, error) {
	if from.Major == to.Major {
		return true, version.Number{}, nil
	}
	// Downgrades not allowed.
	if from.Major > to.Major {
		return false, version.Number{}, nil
	}
	minVer, ok := MinMajorUpgradeVersion[to.Major]
	if !ok {
		return false, version.Number{}, errors.Errorf("unknown version %q", to)
	}
	// Allow upgrades from rc etc.
	from.Tag = ""
	return from.Compare(minVer) >= 0, minVer, nil
}
