// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// abstract base class for cached functions with no arguments
// like Memoize but you only need to cache a single value (don't need LruCache)
// derived classes must define Func
class
	{
	CallClass()
		{
		name = 'MemoizeSingle_' $ Name(this)
		if name is "MemoizeSingle_MemoizeSingle"
			throw "MemoizeSingle must be used as a base class, not called directly"
		if Suneido.Member?(ov = Name(this) $ 'Override')
			result = Suneido[ov]
		else
			result = Suneido.GetInit(name, .Func)
		if Object?(result) and not result.Readonly?()
			result = result.Copy() // won't actually copy unless updated (copy-on-write)
		return result
		}
	ResetCache()
		{
		name = 'MemoizeSingle_' $ Name(this)
		Suneido.Delete(name)
		}
	ClearAll()
		{
		for m in Suneido.Members().Copy()
			if String?(m) and m.Prefix?('MemoizeSingle_')
				Suneido.Delete(m)
		}
	}