package ident

const (
	FixtureNextParent = "fx.nextParent"
	FixtureMaybeSkip  = "fx.maybeSkip"
	// FixtureNotFoundName returns the name of a resource that does not exist,
	// next to the resource with the name it is passed.
	FixtureNotFoundName = "fx.notFoundName"
	// FixtureUserSettableID returns a new user settable ID: the IDGenerator's when set,
	// else the fallback it is passed.
	FixtureUserSettableID = "fx.userSettableID"
)
