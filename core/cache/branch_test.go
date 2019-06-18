// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cache_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/cache"
	"github.com/juju/juju/core/settings"
	"github.com/juju/juju/testing"
)

type BranchSuite struct {
	cache.EntitySuite
}

var _ = gc.Suite(&BranchSuite{})

func (s *BranchSuite) TestBranchSetDetailsPublishesCopy(c *gc.C) {
	hub := s.EnsureHub(nil)

	rcv := make(chan interface{}, 1)
	_ = hub.Subscribe("branch-change", func(_ string, msg interface{}) { rcv <- msg })
	_ = s.NewBranch(branchChange, hub)

	select {
	case msg := <-rcv:
		b, ok := msg.(cache.Branch)
		if !ok {
			c.Fatal("wrong type published; expected Branch.")
		}
		c.Check(b.Name(), gc.Equals, branchChange.Name)

	case <-time.After(testing.LongWait):
		c.Fatal("branch change message not Received")
	}
}

var branchChange = cache.BranchChange{
	ModelUUID:     "model-uuid",
	Id:            "0",
	Name:          "testing-branch",
	AssignedUnits: map[string][]string{"redis": {"redis/0", "redis/1"}},
	Config:        map[string]settings.ItemChanges{"redis": {settings.MakeAddition("password", "pass666")}},
	Completed:     0,
	GenerationId:  0,
}
