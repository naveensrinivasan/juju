// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package payloads

import (
	"github.com/juju/errors"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/payload"
	"github.com/juju/juju/payload/api"
	"github.com/juju/juju/rpc/params"
)

// NewFacade provides the signature required for facade registration.
func NewFacade(ctx facade.Context) (*API, error) {
	authorizer := ctx.Auth()
	if !authorizer.AuthClient() {
		return nil, apiservererrors.ErrPerm
	}
	backend, err := ctx.State().ModelPayloads()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return NewAPI(backend), nil
}

// payloadBackend exposes the State functionality for payloads in a model.
type payloadBackend interface {
	// ListAll returns information on the payload with the id on the unit.
	ListAll() ([]payload.FullPayloadInfo, error)
}

// API serves payload-specific API methods.
type API struct {
	// State exposes the payload aspect of Juju's state.
	backend payloadBackend
}

// NewAPI builds a new facade for the given backend.
func NewAPI(backend payloadBackend) *API {
	return &API{backend: backend}
}

// List builds the list of payloads being tracked for
// the given unit and IDs. If no IDs are provided then all tracked
// payloads for the unit are returned.
func (a API) List(args params.PayloadListArgs) (params.PayloadListResults, error) {
	var r params.PayloadListResults

	payloads, err := a.backend.ListAll()
	if err != nil {
		return r, errors.Trace(err)
	}

	filters, err := payload.BuildPredicatesFor(args.Patterns)
	if err != nil {
		return r, errors.Trace(err)
	}
	payloads = payload.Filter(payloads, filters...)

	for _, payload := range payloads {
		apiInfo := api.Payload2api(payload)
		r.Results = append(r.Results, apiInfo)
	}
	return r, nil
}
