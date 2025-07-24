// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// saves the data and returns a SuneidoApp url for it
	// NOTE: you may need to add an extension to have the file type recognized
	Add(data)
		{
		t = Display(Timestamp()).Tr('#.')
		Suneido.GetInit('InMemory', { Object() })[t] = data
		return 'suneido:/inmemory/' $ t
		}

	// called by SuneidoApp to retrieve the data
	Get(url)
		{
		// ignore any extension added for file type
		if url.Has?('.')
			url = url.BeforeLast('.')
		return Suneido.InMemory[url.Tr('^0-9')]
		}

	// remove the url/data
	Remove(url)
		{
		Suneido.InMemory.Delete(url.Tr('^0-9'))
		return
		}
	}