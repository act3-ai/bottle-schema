// Package selectors provides the ability to work with a collection of label selectors
package selectors

import (
	"k8s.io/apimachinery/pkg/labels"
)

// LabelSelectorSet represents a set of selectors
type LabelSelectorSet []labels.Selector

// TODO maybe Everything and Nothing should be constants or variables and not functions

// Everything matches everything
func Everything() LabelSelectorSet {
	return nil
}

// Nothing matches nothing
func Nothing() LabelSelectorSet {
	return []labels.Selector{}
}

// Parse will parse an array of strings to create a LabelSelectorSet
// If selectors is nil then we match everything this matches everything
// If selectors is empty then we match nothing (no selector will ever match since there are none)
func Parse(selectors []string) (LabelSelectorSet, error) {
	if selectors == nil {
		// We match everything (as if no selectors where provided to reduce the set down)
		return Everything(), nil
	}

	lss := make(LabelSelectorSet, len(selectors))
	for i, selector := range selectors {
		sel, err := labels.Parse(selector)
		if err != nil {
			return nil, err
		}
		lss[i] = sel
	}
	return lss, nil
}

// Matches will return true if any of the selectors match (this implements the OR condition)
func (lss LabelSelectorSet) Matches(l labels.Labels) bool {
	// not an optimization
	if lss == nil {
		// Everything
		return true
	}

	for _, ls := range lss {
		if ls.Matches(l) {
			return true
		}
	}
	return false
}

// MatchAny returns true if the selector set matches any of the label sets
func (lss LabelSelectorSet) MatchAny(lbls []labels.Labels) bool {
	// optimization
	if lss == nil {
		// Everything
		return true
	}

	// optimization
	if len(lss) == 0 {
		return false
	}

	for _, lbl := range lbls {
		if lss.Matches(lbl) {
			return true
		}
	}
	// TODO this can possibly be done more efficiently.  This is the common case (not matching anything).
	return false
}

// MatchAll returns true if the selector set matches all of the label sets
func (lss LabelSelectorSet) MatchAll(lbls []labels.Labels) bool {
	// optimization
	if lss == nil {
		return true
	}

	for _, lbl := range lbls {
		if !lss.Matches(lbl) {
			return false
		}
	}
	// TODO this can possibly be done more efficiently.
	return true
}

// LabelsFromSets converts from a slice of labels.Set to a its interface (labels.Labels).
func LabelsFromSets(s []labels.Set) []labels.Labels {
	// Preserve nil
	if s == nil {
		return nil
	}
	s2 := make([]labels.Labels, len(s))
	for i, v := range s {
		s2[i] = v
	}
	return s2
}
