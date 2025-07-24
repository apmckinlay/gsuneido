// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// TODO: make F1 work somehow
Controller
	{
	New(book, path)
		{
		super(#(LinkButton 'Help' size: '+2'))
		.book = book
		.path = path
		}
	On_Help()
		{
		OpenBook(.book, path: .path)
		}
	}