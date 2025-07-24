// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// like SoleContribution but with a default value
// note: if the contributions are functions,
// you probably want to make the default a function
// e.g. OptContribution(name, function (@args) { <value> })(...)
function (name, def)
	{
	contribs = Contributions(name)
	if contribs.Empty?()
		return def
	if contribs.Size() > 1
		throw "OptContribution: multiple definitions for <lib>_" $ name
	return contribs[0]
	}