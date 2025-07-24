// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// abstract base class for cached functions, wrapping LruCache
// e.g. used by Contributions, QueryKeys
// derived classes must define Func
// and may optionally define CacheSize or Init
// Init is called once before setting up the LruCache
// returns a defensive Copy of mutable objects
// WARNING: this copy is shallow so nested objects could be modified
// it is safer and more efficient to make the return value from Func read-only.
class
	{
	CacheSize: 100
	HashArgs?: false
	OkForResetAll?: true
	Expires?: false
	CallClass(@args)
		{
		name = 'Memoize_' $ Name(this)
		if Suneido.Member?(ov = Name(this) $ 'Override')
			result = Suneido[ov]
		else
			{
			cache = Suneido.GetInit(name,
				{
				Assert(name isnt 'Memoize_', msg: "incorrect Memoize setup")
				if name is "Memoize_Memoize"
					throw "Memoize must be used as a base class, not called directly"

				.Init();
				f = .Expires?
					? .lruCacheExpire
					: LruCache
				f(.Func, .CacheSize,
					okForResetAll?: .OkForResetAll?, expirySeconds: .ExpirySeconds)
				})
			result = cache.Get(@args)
			}
		if Object?(result) and not result.Readonly?()
			result = result.Copy() // won't actually copy unless updated
		return result
		}

	ExpirySeconds: 60
	lruCacheExpire: class
		{
		New(.func, cache_size, .okForResetAll?, .expirySeconds)
			{
			.lru = LruCache(.get, cache_size, okForResetAll?: okForResetAll?)
			}
		get(@args)
			{
			return Object(value: (.func)(@args),
				expiry: Timestamp().Plus(seconds: .expirySeconds))
			}
		Get(@args)
			{
			result = .lru.Get(@args)
			if result.expiry < Timestamp()
				result.Merge(.get(@args))
			return result.value
			}
		Reset()
			{
			.lru.Reset()
			}
		ResetExpire()
			{
			if .okForResetAll?
				.lru.Reset()
			}
		}
	Init()
		{
		}
	ResetCache()
		{
		name = 'Memoize_' $ Name(this)
		if false isnt cache = Suneido.GetDefault(name, false)
			cache.Reset()
		}
	ResetAllExpire()
		{
		for mem in Suneido.Members().Copy()
			if Type(Suneido[mem]) is 'Instance' and Suneido[mem].Base?(.lruCacheExpire)
				Suneido[mem].ResetExpire()
		}
	}
