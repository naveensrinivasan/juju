// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//go:build !windows
// +build !windows

package unit_test

import (
	"testing"

	"github.com/juju/collections/set"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	coretesting "github.com/juju/juju/testing"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type importSuite struct{}

var _ = gc.Suite(&importSuite{})

func (*importSuite) TestImports(c *gc.C) {
	// TODO(sidecar) - improve test performance
	c.Skip("test times out on Jenkins")
	found := set.NewStrings(
		coretesting.FindJujuCoreImports(c, "github.com/juju/juju/cmd/containeragent/unit")...)

	expected := set.NewStrings(
		"agent",
		"agent/tools",
		"api",
		"api/agent/agent",
		"api/authentication",
		"api/base",
		"api/client/block",
		"api/agent/caasoperator",
		"api/common",
		"api/common/cloudspec",
		"api/controller/controller",
		"api/controller/instancepoller",
		"api/agent/keyupdater",
		"api/agent/leadership",
		"api/agent/logger",
		"api/logsender",
		"api/agent/machiner",
		"api/agent/migrationflag",
		"api/agent/migrationminion",
		"api/client/modelmanager",
		"api/agent/proxyupdater",
		"api/agent/reboot",
		"api/agent/retrystrategy",
		"api/agent/unitassigner",
		"api/agent/uniter",
		"api/agent/upgrader",
		"api/client/usermanager",
		"api/watcher",
		"cmd/containeragent/utils",
		"apiserver/errors",
		"rpc/params",
		"apiserver/apiserverhttp",
		"caas",
		"caas/kubernetes/clientconfig",
		"caas/kubernetes/provider",
		"caas/kubernetes/provider/application",
		"caas/kubernetes/provider/constants",
		"caas/kubernetes/provider/proxy",
		"caas/kubernetes/provider/resources",
		"caas/kubernetes/provider/specs",
		"caas/kubernetes/provider/storage",
		"caas/kubernetes/provider/utils",
		"caas/kubernetes/provider/watcher",
		"caas/specs",
		"charmhub",
		"charmhub/path",
		"charmhub/transport",
		"charmstore",
		"cloud",
		"cloudconfig",
		"cloudconfig/cloudinit",
		"cloudconfig/instancecfg",
		"cloudconfig/podcfg",
		"cmd",
		"cmd/juju/common",
		"cmd/juju/interact",
		"cmd/jujud/agent/addons",
		"cmd/jujud/agent/agentconf",
		"cmd/jujud/agent/config",
		"cmd/jujud/agent/engine",
		"cmd/jujud/agent/errors",
		"cmd/modelcmd",
		"cmd/output",
		"controller",
		"core/actions",
		"core/annotations",
		"core/application",
		"core/constraints",
		"core/devices",
		"core/globalclock",
		"core/instance",
		"core/leadership",
		"core/lease",
		"core/life",
		"core/lxdprofile",
		"core/machinelock",
		"core/migration",
		"core/model",
		"core/network",
		"core/network/firewall",
		"core/paths",
		"core/paths/transientfile",
		"core/permission",
		"core/presence",
		"core/quota",
		"core/raftlease",
		"core/relation",
		"core/resources",
		"core/secrets",
		"core/series",
		"core/snap",
		"core/status",
		"core/watcher",
		"downloader",
		"environs",
		"environs/bootstrap",
		"environs/cloudspec",
		"environs/config",
		"environs/context",
		"environs/filestorage",
		"environs/gui",
		"environs/imagemetadata",
		"environs/instances",
		"environs/simplestreams",
		"environs/storage",
		"environs/sync",
		"environs/tags",
		"environs/tools",
		"environs/utils",
		"feature",
		"juju",
		"juju/keys",
		"juju/names",
		"juju/osenv",
		"juju/sockets",
		"jujuclient",
		"logfwd",
		"logfwd/syslog",
		"mongo", // TODO: move mongo dependency from JUJU CLI if we decide to split the `agent.Config` for controller and machineagent/unitagent/containeragent.
		"network",
		"network/debinterfaces",
		"network/netplan",
		"packaging",
		"packaging/dependency",
		"pki",
		"pki/tls",
		"proxy",
		"provider/lxd/lxdnames",
		"pubsub/agent",
		"resource",
		"rpc",
		"rpc/jsoncodec",
		"service",
		"service/common",
		"service/snap",
		"service/systemd",
		"service/upstart",
		"service/windows",
		"state/errors",
		"storage",
		"storage/provider",
		"tools",
		"utils/proxy",
		"utils/scriptrunner",
		"version",
		"worker",
		"worker/agent",
		"worker/apiaddressupdater",
		"worker/apicaller",
		"worker/apiconfigwatcher",
		"worker/caasprober",
		"worker/common/charmrunner",
		"worker/common/reboot",
		"worker/fortress",
		"worker/introspection",
		"worker/introspection/pprof",
		"worker/leadership",
		"worker/logger",
		"worker/logsender",
		"worker/migrationflag",
		"worker/migrationminion",
		"worker/muxhttpserver",
		"worker/proxyupdater",
		"worker/retrystrategy",
		"worker/uniter",
		"worker/uniter/actions",
		"worker/uniter/charm",
		"worker/uniter/container",
		"worker/uniter/hook",
		"worker/uniter/leadership",
		"worker/uniter/operation",
		"worker/uniter/reboot",
		"worker/uniter/relation",
		"worker/uniter/remotestate",
		"worker/uniter/resolver",
		"worker/uniter/runcommands",
		"worker/uniter/runner",
		"worker/uniter/runner/context",
		"worker/uniter/runner/debug",
		"worker/uniter/runner/jujuc",
		"worker/uniter/storage",
		"worker/uniter/upgradeseries",
		"worker/uniter/verifycharmprofile",
		"wrench",
	)

	unexpected := found.Difference(expected)
	// TODO: review if there are any un-expected imports!
	// Show the values rather than just checking the length so a failing
	// test shows them.
	c.Check(unexpected.SortedValues(), jc.DeepEquals, []string{})
	// If unneeded show any values this is good as we've reduced
	// dependencies, and they should be removed from expected above.
	unneeded := expected.Difference(found)
	c.Check(unneeded.SortedValues(), jc.DeepEquals, []string{})
}
