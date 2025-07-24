// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// This should normally be used instead of calling Unload directly.
// It handles calling observers to reset caches etc.
class
	{
	CallClass(name)
		{
//		name = LibraryTags.RemoveTagsFromName(name)
		Unload(name)
		if Sys.Client?()
			ServerEval('LibUnload', name)
		if Suneido.Member?('LibUnload_observers')
			for observer in Suneido.LibUnload_observers.Copy()
				observer(name)
		}
	AddObserver(key, observer)
		{
		if not Suneido.Member?('LibUnload_observers')
			Suneido.LibUnload_observers = Object()
		Suneido.LibUnload_observers.Add(observer, at: key)
		}
	RemoveObserver(key)
		{
		if Suneido.Member?('LibUnload_observers')
			Suneido.LibUnload_observers.Delete(key)
		}
	}