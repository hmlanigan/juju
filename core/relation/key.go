// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation

type Key string

func (k Key) String() string {
	return string(k)
}
