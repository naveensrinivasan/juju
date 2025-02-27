// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/base"
	apiwatcher "github.com/juju/juju/api/watcher"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/rpc/params"
)

// Client provides access to the action facade.
type Client struct {
	base.ClientFacade
	facade base.FacadeCaller
}

// NewClient returns a new actions client.
func NewClient(st base.APICallCloser) *Client {
	frontend, backend := base.NewClientFacade(st, "Action")
	return &Client{ClientFacade: frontend, facade: backend}
}

// Actions takes a list of ActionTags, and returns the full
// Action for each ID.
func (c *Client) Actions(actionIDs []string) ([]ActionResult, error) {
	arg := params.Entities{Entities: make([]params.Entity, len(actionIDs))}
	for i, ID := range actionIDs {
		arg.Entities[i].Tag = names.NewActionTag(ID).String()
	}
	results := params.ActionResults{}
	err := c.facade.FacadeCall("Actions", arg, &results)
	return unmarshallActionResults(results.Results), err
}

// ListOperations fetches the operation summaries for specified apps/units.
func (c *Client) ListOperations(arg OperationQueryArgs) (Operations, error) {
	if v := c.BestAPIVersion(); v < 6 {
		return Operations{}, errors.Errorf("ListOperations not supported by this version (%d) of Juju", v)
	}
	args := params.OperationQueryArgs{
		Applications: arg.Applications,
		Units:        arg.Units,
		Machines:     arg.Machines,
		ActionNames:  arg.ActionNames,
		Status:       arg.Status,
		Offset:       arg.Offset,
		Limit:        arg.Limit,
	}
	results := params.OperationResults{}
	err := c.facade.FacadeCall("ListOperations", args, &results)
	if params.ErrCode(err) == params.CodeNotFound {
		err = nil
	}
	return unmarshallOperations(results), err
}

// Operation fetches the operation with the specified id.
func (c *Client) Operation(id string) (Operation, error) {
	if v := c.BestAPIVersion(); v < 6 {
		return Operation{}, errors.Errorf("Operations not supported by this version (%d) of Juju", v)
	}
	arg := params.Entities{
		Entities: []params.Entity{{names.NewOperationTag(id).String()}},
	}
	var results params.OperationResults
	err := c.facade.FacadeCall("Operations", arg, &results)
	if err != nil {
		return Operation{}, err
	}
	if len(results.Results) != 1 {
		return Operation{}, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return Operation{}, maybeNotFound(result.Error)
	}
	return unmarshallOperation(result), nil
}

// maybeNotFound returns an error satisfying errors.IsNotFound
// if the supplied error has a CodeNotFound error.
func maybeNotFound(err *params.Error) error {
	if err == nil || !params.IsCodeNotFound(err) {
		return err
	}
	return errors.NewNotFound(err, "")
}

// FindActionTagsByPrefix takes a list of string prefixes and finds
// corresponding ActionTags that match that prefix.
func (c *Client) FindActionTagsByPrefix(arg params.FindTags) (params.FindTagsResults, error) {
	results := params.FindTagsResults{}
	err := c.facade.FacadeCall("FindActionTagsByPrefix", arg, &results)
	return results, err
}

// Enqueue takes a list of Actions and queues them up to be executed by
// the designated ActionReceiver, returning the params.Action for each
// queued Action, or an error if there was a problem queueing up the
// Action.
func (c *Client) Enqueue(actions []Action) ([]ActionResult, error) {
	arg := params.Actions{Actions: make([]params.Action, len(actions))}
	for i, a := range actions {
		arg.Actions[i] = params.Action{
			Receiver:   a.Receiver,
			Name:       a.Name,
			Parameters: a.Parameters,
		}
	}
	results := params.ActionResults{}
	err := c.facade.FacadeCall("Enqueue", arg, &results)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return unmarshallActionResults(results.Results), nil
}

// EnqueueOperation takes a list of Actions and queues them up to be executed as
// an operation, each action running as a task on the the designated ActionReceiver.
// We return the ID of the overall operation and each individual task.
func (c *Client) EnqueueOperation(actions []Action) (EnqueuedActions, error) {
	if v := c.BestAPIVersion(); v < 6 {
		return EnqueuedActions{}, errors.Errorf("EnqueueOperation not supported by this version (%d) of Juju", v)
	}
	arg := params.Actions{Actions: make([]params.Action, len(actions))}
	for i, a := range actions {
		arg.Actions[i] = params.Action{
			Receiver:   a.Receiver,
			Name:       a.Name,
			Parameters: a.Parameters,
		}
	}
	results := params.EnqueuedActions{}
	err := c.facade.FacadeCall("EnqueueOperation", arg, &results)
	if err != nil {
		return EnqueuedActions{}, errors.Trace(err)
	}
	return unmarshallEnqueuedActions(results)
}

// FindActionsByNames takes a list of action names and returns actions for
// every name.
func (c *Client) FindActionsByNames(arg params.FindActionsByNames) (map[string][]ActionResult, error) {
	results := params.ActionsByNames{}
	err := c.facade.FacadeCall("FindActionsByNames", arg, &results)
	result := make(map[string][]ActionResult)
	for _, a := range results.Actions {
		result[a.Name] = unmarshallActionResults(a.Actions)
	}
	return result, err
}

// Cancel attempts to cancel a queued up Action from running.
func (c *Client) Cancel(actionIDs []string) ([]ActionResult, error) {
	arg := params.Entities{Entities: make([]params.Entity, len(actionIDs))}
	for i, ID := range actionIDs {
		arg.Entities[i].Tag = names.NewActionTag(ID).String()
	}
	results := params.ActionResults{}
	err := c.facade.FacadeCall("Cancel", arg, &results)
	return unmarshallActionResults(results.Results), err
}

// applicationsCharmActions is a batched query for the charm.Actions for a slice
// of applications by Entity.
func (c *Client) applicationsCharmActions(arg params.Entities) (params.ApplicationsCharmActionsResults, error) {
	results := params.ApplicationsCharmActionsResults{}
	err := c.facade.FacadeCall("ApplicationsCharmsActions", arg, &results)
	return results, err
}

// ApplicationCharmActions is a single query which uses ApplicationsCharmsActions to
// get the charm.Actions for a single Application by tag.
func (c *Client) ApplicationCharmActions(appName string) (map[string]ActionSpec, error) {
	tag := names.NewApplicationTag(appName)
	tags := params.Entities{Entities: []params.Entity{{Tag: tag.String()}}}
	results, err := c.applicationsCharmActions(tags)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("%d results, expected 1", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	if result.ApplicationTag != tag.String() {
		return nil, errors.Errorf("action results received for wrong application %q", result.ApplicationTag)
	}
	return unmarshallActionSpecs(result.Actions), nil
}

// WatchActionProgress returns a watcher that reports on action log messages.
// The result strings are json formatted core.actions.ActionMessage objects.
func (c *Client) WatchActionProgress(actionId string) (watcher.StringsWatcher, error) {
	if v := c.BestAPIVersion(); v < 5 {
		return nil, errors.Errorf("WatchActionProgress not supported by this version (%d) of Juju", v)
	}
	var results params.StringsWatchResults
	args := params.Entities{
		Entities: []params.Entity{
			{Tag: names.NewActionTag(actionId).String()},
		},
	}
	err := c.facade.FacadeCall("WatchActionsProgress", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := apiwatcher.NewStringsWatcher(c.facade.RawAPICaller(), result)
	return w, nil
}
