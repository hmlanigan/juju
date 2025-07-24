// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

// Lifer represents an entity with a life.
type Lifer interface {
	Life() Life
}
