// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (ctrl)
	{
	if ctrl.Method?(#TopDown)
		ctrl.TopDown(#Startup)
	else if ctrl.Method?(#Startup)
		ctrl.Startup()
	}