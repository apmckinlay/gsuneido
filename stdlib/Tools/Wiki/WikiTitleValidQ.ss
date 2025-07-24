// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (page)
	{
	return not String?(page) or page !~ `\A[A-Z][a-z0-9]+[A-Z][a-zA-Z0-9]*\Z`
		? "ERROR: invalid page name: " $ page
		: ''
	}