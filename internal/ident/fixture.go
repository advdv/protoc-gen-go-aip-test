package ident

const (
	FixtureNextParent = "fx.nextParent"
	FixtureMaybeSkip  = "fx.maybeSkip"
	// FixtureNormalize returns a copy of a resource that the config's Normalize has changed,
	// or the resource itself when Normalize is nil, for the tests to compare.
	FixtureNormalize = "fx.normalize"
	// FixtureNormalizeAll is FixtureNormalize for each resource of a slice.
	FixtureNormalizeAll = "fx.normalizeAll"
	// FixtureListed checks that the resources listed under a parent are the ones the test
	// created plus the config's Preexisting others, all under the parent, and returns the
	// listed resources that the test created, for it to compare with the ones it created.
	FixtureListed = "fx.listed"
	// FixtureNotFoundName returns the name of a resource that does not exist,
	// next to the resource with the name it is passed.
	FixtureNotFoundName = "fx.notFoundName"
	// FixtureUserSettableID returns a new user settable ID: the IDGenerator's when set,
	// else the fallback it is passed.
	FixtureUserSettableID = "fx.userSettableID"
)
