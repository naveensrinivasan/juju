// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package diskmanager

import (
	"github.com/juju/names/v4"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	"github.com/juju/juju/storage"
)

// DiskManagerAPI provides access to the DiskManager API facade.
type DiskManagerAPI struct {
	st          stateInterface
	authorizer  facade.Authorizer
	getAuthFunc common.GetAuthFunc
}

var getState = func(st *state.State) stateInterface {
	return stateShim{st}
}

// NewDiskManagerAPI creates a new server-side DiskManager API facade.
func NewDiskManagerAPI(ctx facade.Context) (*DiskManagerAPI, error) {
	authorizer := ctx.Auth()
	if !authorizer.AuthMachineAgent() {
		return nil, apiservererrors.ErrPerm
	}

	authEntityTag := authorizer.GetAuthTag()
	getAuthFunc := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			// A machine agent can always access its own machine.
			return tag == authEntityTag
		}, nil
	}

	st := ctx.State()
	return &DiskManagerAPI{
		st:          getState(st),
		authorizer:  authorizer,
		getAuthFunc: getAuthFunc,
	}, nil
}

func (d *DiskManagerAPI) SetMachineBlockDevices(args params.SetMachineBlockDevices) (params.ErrorResults, error) {
	result := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.MachineBlockDevices)),
	}
	canAccess, err := d.getAuthFunc()
	if err != nil {
		return result, err
	}
	for i, arg := range args.MachineBlockDevices {
		tag, err := names.ParseMachineTag(arg.Machine)
		if err != nil {
			result.Results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			continue
		}
		if !canAccess(tag) {
			err = apiservererrors.ErrPerm
		} else {
			// TODO(axw) create volumes for block devices without matching
			// volumes, if and only if the block device has a serial. Under
			// the assumption of unique (to a machine) serial IDs, this
			// gives us a guaranteed *persistently* unique way of identifying
			// the volume.
			//
			// NOTE: we must predicate the above on there being no unprovisioned
			// volume attachments for the machine, otherwise we would have
			// a race between the volume attachment info being recorded and
			// the diskmanager publishing block devices and erroneously creating
			// volumes.
			err = d.st.SetMachineBlockDevices(tag.Id(), stateBlockDeviceInfo(arg.BlockDevices))
			// TODO(axw) set volume/filesystem attachment info.
		}
		result.Results[i].Error = apiservererrors.ServerError(err)
	}
	return result, nil
}

func stateBlockDeviceInfo(devices []storage.BlockDevice) []state.BlockDeviceInfo {
	result := make([]state.BlockDeviceInfo, len(devices))
	for i, dev := range devices {
		result[i] = state.BlockDeviceInfo{
			DeviceName:     dev.DeviceName,
			DeviceLinks:    dev.DeviceLinks,
			Label:          dev.Label,
			UUID:           dev.UUID,
			HardwareId:     dev.HardwareId,
			WWN:            dev.WWN,
			BusAddress:     dev.BusAddress,
			Size:           dev.Size,
			FilesystemType: dev.FilesystemType,
			InUse:          dev.InUse,
			MountPoint:     dev.MountPoint,
			SerialId:       dev.SerialId,
		}
	}
	return result
}
