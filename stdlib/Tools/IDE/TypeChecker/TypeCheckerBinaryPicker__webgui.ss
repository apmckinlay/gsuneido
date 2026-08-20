// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// LEGACY-TYPECHECK-BINARY: false path means there is no binary, so render nothing
function(path)
	{
	return path is false
		? #Skip
		: [#Field, set: path, name: #TypeCheckerBinary, width: 50]
	}
