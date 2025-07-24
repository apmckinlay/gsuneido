// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if .Function isnt "count" and .Value is ""
		return "Value field required"
	if .Rows isnt "" and .Columns isnt "" and .Rows is .Columns
		return "Rows field must be different from Columns field"
	return ""
	}