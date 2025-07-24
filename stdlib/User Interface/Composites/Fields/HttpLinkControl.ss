// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
LinkControl
	{
	Name: 'HttpLink'
	Prefix: 'http://'
	Status: 'e.g. www.axonsoft.com, double click to view in your browser'
	New(@args)
		{
		super(@args)
		.AddContextMenuItem("", "")
		.AddContextMenuItem("Go To Website", .GoToLink)
		}
	MergePrefix(adr)
		{
		if adr.Prefix?(.Prefix) or adr.Prefix?('https://')
			return adr
		return .Prefix $ adr
		}
	}
