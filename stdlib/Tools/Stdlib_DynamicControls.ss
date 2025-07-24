// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
#(
	Address: function(@args)
		{
		return AddressControl.Layout(@args)
		},
	ShowHide: function(@args)
		{
		showFn = args.Extract(1)
		return true is (String?(showFn) ? Global(showFn) : showFn)()
			? args
			: #()
		},
	NameAbbrev: function(@args)
		{
		return NameAbbrevControl.Layout(@args)
		}
)
