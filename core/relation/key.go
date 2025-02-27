// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation

type Key string

// ParseKey returns a relation Key.
//func ParseKey(value string) (Key, error) {
//	if value == "" {
//		return "", fmt.Errorf("id %q: %w", value, errors.NotValid)
//	}
//	epNames := strings.Split(value, " ")
//	switch len(epNames) {
//	case 1:
//		return Key(value), nil
//	case 2:
//		sort.Sort(epNames)
//		endpointNames := []string{}
//		for _, ep := range epNames {
//			endpointNames = append(endpointNames, ep)
//		}
//		return Key(strings.Join(endpointNames, " ")), nil
//	default:
//		return "", fmt.Errorf("id %q: %w", value, errors.NotValid)
//	}
//}

//func relationKey(endpoints []relation.Endpoint) string {
//	eps := epSlice{}
//	for _, ep := range endpoints {
//		eps = append(eps, ep)
//	}
//	sort.Sort(eps)
//	endpointNames := []string{}
//	for _, ep := range eps {
//		endpointNames = append(endpointNames, ep.String())
//	}
//	return strings.Join(endpointNames, " ")
//}
