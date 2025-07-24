// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (err = "Transaction: block commit failed")
	{
	msg = "Your changes conflict with another user. Unable to save."
	SuneidoLog("ERROR: (CAUGHT) " $ err,
		caughtMsg: "user alerted: " $ msg)
	Alert(msg, title: 'Error', flags: MB.ICONERROR)
	}