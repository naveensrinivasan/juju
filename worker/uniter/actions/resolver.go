// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package actions

import (
	"github.com/juju/juju/worker/uniter/operation"
	"github.com/juju/juju/worker/uniter/remotestate"
	"github.com/juju/juju/worker/uniter/resolver"
)

// Logger is here to stop the desire of creating a package level Logger.
// Don't do this, instead use the one passed into the NewResolver as needed.
type logger interface{}

var _ logger = struct{}{}

// Logger represents the logging methods used by the actions resolver.
type Logger interface {
	Infof(string, ...interface{})
}

type actionsResolver struct {
	logger Logger
}

// NewResolver returns a new resolver with determines which action related operation
// should be run based on local and remote uniter states.
//
// TODO(axw) 2015-10-27 #1510333
// Use the same method as in the runcommands resolver
// for updating the remote state snapshot when an
// action is completed.
func NewResolver(logger Logger) resolver.Resolver {
	return &actionsResolver{logger: logger}
}

func nextAction(pendingActions []string, completedActions map[string]struct{}) (string, error) {
	for _, action := range pendingActions {
		if _, ok := completedActions[action]; !ok {
			return action, nil
		}
	}
	return "", resolver.ErrNoOperation
}

// NextOp implements the resolver.Resolver interface.
func (r *actionsResolver) NextOp(
	localState resolver.LocalState,
	remoteState remotestate.Snapshot,
	opFactory operation.Factory,
) (operation.Operation, error) {
	// During CAAS unit initialization action operations are
	// deferred until the unit is running. If the remote charm needs
	// updating, hold off on action running.
	if remoteState.ActionsBlocked || localState.OutdatedRemoteCharm {
		// If we were somehow running an action during remote container changes/restart
		// we need to fail it and move on.
		if localState.Kind == operation.RunAction {
			if localState.Hook != nil {
				r.logger.Infof("found incomplete action %v; ignoring", localState.ActionId)
				r.logger.Infof("recommitting prior %q hook", localState.Hook.Kind)
				return opFactory.NewSkipHook(*localState.Hook)
			}
			return opFactory.NewFailAction(*localState.ActionId)
		}
		return nil, resolver.ErrNoOperation
	}
	// If there are no operation left to be run, then we cannot return the
	// error signaling such here, we must first check to see if an action is
	// already running (that has been interrupted) before we declare that
	// there is nothing to do.
	nextAction, err := nextAction(remoteState.ActionsPending, localState.CompletedActions)
	if err != nil && err != resolver.ErrNoOperation {
		return nil, err
	}
	switch localState.Kind {
	case operation.RunHook:
		// We can still run actions if the unit is in a hook error state.
		if localState.Step == operation.Pending && err == nil {
			return opFactory.NewAction(nextAction)
		}
	case operation.RunAction:
		if localState.Hook != nil {
			r.logger.Infof("found incomplete action %v; ignoring", localState.ActionId)
			r.logger.Infof("recommitting prior %q hook", localState.Hook.Kind)
			return opFactory.NewSkipHook(*localState.Hook)
		}

		r.logger.Infof("%q hook is nil", operation.RunAction)
		// If the next action is the same as what the uniter is
		// currently running then this means that the uniter was
		// some how interrupted (killed) when running the action
		// and before updating the remote state to indicate that
		// the action was completed. The only safe thing to do
		// is fail the action, since rerunning an arbitrary
		// command can potentially be hazardous.
		if nextAction == *localState.ActionId {
			return opFactory.NewFailAction(*localState.ActionId)
		}

		// If the next action is different then what the uniter
		// is currently running, then the uniter may have been
		// interrupted while running the action but the remote
		// state was updated. Thus, the semantics of
		// (re)preparing the running operation should move the
		// uniter's state along safely. Thus, we return the
		// running action.
		return opFactory.NewAction(*localState.ActionId)
	case operation.Continue:
		if err != resolver.ErrNoOperation {
			return opFactory.NewAction(nextAction)
		}
	}
	return nil, resolver.ErrNoOperation
}
