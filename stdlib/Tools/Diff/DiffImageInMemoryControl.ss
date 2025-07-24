// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
DiffImageControl
	{
	url1: false
	url2: false
	New(text1, text2, .Title1 = 'Before', .Title2 = 'After')
		{
		super('', '',
			ImageViewer.Style(.url1 = InMemory.Add(text1))
			ImageViewer.Style(.url2 = InMemory.Add(text2)))
		}

	Destroy()
		{
		if .url1 isnt false
			InMemory.Remove(.url1)
		if .url2 isnt false
			InMemory.Remove(.url2)
		super.Destroy()
		}
	}
