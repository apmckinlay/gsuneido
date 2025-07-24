// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// returns the last contribution in Libraries() order
// throws if there are no contributions
function (name)
	{
	contribs = Contributions(name)
	if contribs.Empty?()
		throw "LastContribution: no definition for <lib>_" $ name
	return contribs.Last()
	}