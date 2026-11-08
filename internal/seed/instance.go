package seed

import "errors"

// ErrNoInstance is a template with no instance to choose from.
var ErrNoInstance = errors.New("seed: a template has no instance")

// Instance gives the number of the instance of a template that an attempt takes, as the platform
// chooses it: below(count) of the place of the attempt's seed with the names template, the task's
// ID and instance (the text of the format, section 23; R60 spec). Nothing else draws from that
// place. count is the number of the template's instances, numbered from 0.
func Instance(attemptSeed uint64, task string, count int) (int, error) {
	if count <= 0 {
		return 0, ErrNoInstance
	}
	s, err := Place(attemptSeed, "template", task, "instance")
	if err != nil {
		return 0, err
	}
	return int(s.Below(uint64(count))), nil
}
