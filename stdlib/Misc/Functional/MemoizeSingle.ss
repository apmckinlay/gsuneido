// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// abstract base class for cached functions with no arguments
// like Memoize but you only need to cache a single value (don't need LruCache)
// derived classes must define Func
class
	{
	CallClass()
		{
		name = 'MemoizeSingle_' $ Name(this)
		if name is 'MemoizeSingle_MemoizeSingle'
			throw 'MemoizeSingle must be used as a base class, not called directly'
		cache = .Cache()
		if cache.Member?(ov = Name(this) $ 'Override')
			result = cache[ov]
		else
			result = cache.GetInit(name, .Func)
		if Object?(result) and not result.Readonly?()
			result = result.Copy() // won't actually copy unless updated (copy-on-write)
		return result
		}

	cacheMember: 'MemoizeSingleCache'
	Cache()
		{
		return Suneido.GetInit(.cacheMember, Object)
		}

	ResetCache()
		{
		.Cache().Delete('MemoizeSingle_' $ Name(this))
		}

	ClearAll()
		{
		Suneido[.cacheMember] = Object()
		}
	}
