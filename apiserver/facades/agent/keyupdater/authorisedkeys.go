// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keyupdater

import (
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"github.com/juju/utils/v3/ssh"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	"github.com/juju/juju/state/watcher"
)

// KeyUpdater defines the methods on the keyupdater API end point.
type KeyUpdater interface {
	AuthorisedKeys(args params.Entities) (params.StringsResults, error)
	WatchAuthorisedKeys(args params.Entities) (params.NotifyWatchResults, error)
}

// KeyUpdaterAPI implements the KeyUpdater interface and is the concrete
// implementation of the api end point.
type KeyUpdaterAPI struct {
	state      *state.State
	model      *state.Model
	resources  facade.Resources
	authorizer facade.Authorizer
	getCanRead common.GetAuthFunc
}

var _ KeyUpdater = (*KeyUpdaterAPI)(nil)

// NewKeyUpdaterAPI creates a new server-side keyupdater API end point.
func NewKeyUpdaterAPI(ctx facade.Context) (*KeyUpdaterAPI, error) {
	authorizer := ctx.Auth()
	// Only machine agents have access to the keyupdater service.
	if !authorizer.AuthMachineAgent() {
		return nil, apiservererrors.ErrPerm
	}
	// No-one else except the machine itself can only read a machine's own credentials.
	getCanRead := func() (common.AuthFunc, error) {
		return authorizer.AuthOwner, nil
	}
	st := ctx.State()
	m, err := st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &KeyUpdaterAPI{
		state:      st,
		model:      m,
		resources:  ctx.Resources(),
		authorizer: authorizer,
		getCanRead: getCanRead,
	}, nil
}

// WatchAuthorisedKeys starts a watcher to track changes to the authorised ssh keys
// for the specified machines.
// The current implementation relies on global authorised keys being stored in the model config.
// This will change as new user management and authorisation functionality is added.
func (api *KeyUpdaterAPI) WatchAuthorisedKeys(arg params.Entities) (params.NotifyWatchResults, error) {
	results := make([]params.NotifyWatchResult, len(arg.Entities))

	canRead, err := api.getCanRead()
	if err != nil {
		return params.NotifyWatchResults{}, err
	}
	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = apiservererrors.ServerError(err)
			continue
		}
		// 1. Check permissions
		if !canRead(tag) {
			results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			continue
		}
		// 2. Check entity exists
		if _, err := api.state.FindEntity(tag); err != nil {
			if errors.IsNotFound(err) {
				results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			} else {
				results[i].Error = apiservererrors.ServerError(err)
			}
			continue
		}
		// 3. Watch for changes
		watch := api.model.WatchForModelConfigChanges()
		// Consume the initial event.
		if _, ok := <-watch.Changes(); ok {
			results[i].NotifyWatcherId = api.resources.Register(watch)
		} else {
			err = watcher.EnsureErr(watch)
		}
		results[i].Error = apiservererrors.ServerError(err)
	}
	return params.NotifyWatchResults{Results: results}, nil
}

// AuthorisedKeys reports the authorised ssh keys for the specified machines.
// The current implementation relies on global authorised keys being stored in the model config.
// This will change as new user management and authorisation functionality is added.
func (api *KeyUpdaterAPI) AuthorisedKeys(arg params.Entities) (params.StringsResults, error) {
	if len(arg.Entities) == 0 {
		return params.StringsResults{}, nil
	}
	results := make([]params.StringsResult, len(arg.Entities))

	// For now, authorised keys are global, common to all machines.
	var keys []string
	config, configErr := api.model.ModelConfig()
	if configErr == nil {
		keys = ssh.SplitAuthorisedKeys(config.AuthorizedKeys())
	}

	canRead, err := api.getCanRead()
	if err != nil {
		return params.StringsResults{}, err
	}
	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = apiservererrors.ServerError(err)
			continue
		}
		// 1. Check permissions
		if !canRead(tag) {
			results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			continue
		}
		// 2. Check entity exists
		if _, err := api.state.FindEntity(tag); err != nil {
			if errors.IsNotFound(err) {
				results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			} else {
				results[i].Error = apiservererrors.ServerError(err)
			}
			continue
		}
		// 3. Get keys
		if configErr == nil {
			results[i].Result = keys
		} else {
			err = configErr
		}
		results[i].Error = apiservererrors.ServerError(err)
	}
	return params.StringsResults{Results: results}, nil
}
