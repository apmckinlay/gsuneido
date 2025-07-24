// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return .Member?(#params) and Object?(.params) and .params.Member?(#Source)
		? .params.Source
		: ''
	}
