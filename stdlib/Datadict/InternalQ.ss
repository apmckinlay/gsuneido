// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function(field)
	{
	return Datadict(field, #(Internal)).GetDefault(#Internal, false)
	}