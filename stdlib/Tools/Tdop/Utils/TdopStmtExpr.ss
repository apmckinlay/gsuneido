// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(@args)
	{
	prevStmt = _setStmtnest(0)
	e = _expr(@args)
	_setStmtnest(prevStmt)
	return e
	}
