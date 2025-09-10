// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"fmt"

	"github.com/juju/collections/transform"

	"github.com/juju/juju/core/changestream"
	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/uuid"
)

// WatchTaskAbortingForReceiver watches for any task for the receiverUUID to
// have a status of aborting. TaskIDs are sent on via the returned strings
// watcher.
func (w *WatchableService) WatchTaskAbortingForReceiver(
	ctx context.Context,
	receiverUUID uuid.UUID,
) (watcher.StringsWatcher, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	initialQuery := func(ctx context.Context, txn database.TxnRunner) ([]string, error) {
		ctx, span := trace.Start(ctx, "WatchTaskAbortingForReceiver.initialQuery")
		defer span.End()

		ids, err := w.st.GetIDsForAbortingTaskOfReceiver(ctx, receiverUUID)
		if err != nil {
			return nil, errors.Errorf("%q: %w", receiverUUID, err)
		}

		return ids, err
	}

	mapper := func(ctx context.Context, changes []changestream.ChangeEvent) ([]string, error) {
		ctx, span := trace.Start(ctx, "WatchTaskAbortingForReceiver.mapper")
		defer span.End()

		taskUUIDs := transform.Slice(changes, func(in changestream.ChangeEvent) string {
			return in.Changed()
		})

		// The namespace watched, only triggers when a task status is
		// set to ABORTING. Find which tasks are for the receiver
		// provided to the watcher.
		taskIDs, err := w.st.GetTaskIDsByUUIDsFilteredByReceiverUUID(ctx, receiverUUID, taskUUIDs)
		if err != nil {
			return nil, errors.Errorf("task aborted watcher mapper %q: %w", receiverUUID, err)
		}
		return taskIDs, nil
	}

	return w.watcherFactory.NewNamespaceMapperWatcher(
		ctx,
		initialQuery,
		fmt.Sprintf("aborting status task watcher for %q", receiverUUID),
		mapper,
		eventsource.NamespaceFilter(w.st.NamespaceForTaskAbortingWatcher(), changestream.Changed),
	)
}
