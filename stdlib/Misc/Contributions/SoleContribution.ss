// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// returns a single contribution for a name
// throws if there are no contributions, or if there are more than one
function (name)
	{
	contribs = Contributions(name)
	if contribs.Empty?()
		throw "SoleContribution: no definition for <lib>_" $ name
	if contribs.Size() > 1
		throw "SoleContribution: multiple definitions for <lib>_" $ name
	return contribs[0]
	}