// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Get(member, defaultVal = false)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.Get", member, defaultVal)
		.ensureAllowed(member)
		return Suneido.GetDefault(member, defaultVal)
		}

	Set(member, value)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.Set", member, value)
		.ensureAllowed(member)
		Suneido[member] = value
		return
		}

	DeleteMember(member)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.DeleteMember", member)
		.ensureAllowed(member)
		Suneido.Delete(member)
		return
		}

	DeleteAt(member, at)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.DeleteAt", member, at)
		.ensureAllowed(member)
		.validateObject(member)
		Suneido[member].Delete(at)
		}

	RemoveFromMember(member, value)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.RemoveFromMember", member, value)
		.ensureAllowed(member)
		.validateObject(member)
		Suneido[member].Remove(value)
		}

	AddToMember(member, value)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.AddToMember", member, value)
		.ensureAllowed(member)
		.validateObject(member)
		Suneido[member].Add(value)
		}

	Add(member, value, at)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.Add", member, value, at)
		.ensureAllowed(member)
		.validateObject(member)
		Suneido[member][at] = value
		return
		}

	validateObject(member)
		{
		if not Suneido.Member?(member)
			Suneido[member] = Object()
		if not Object?(Suneido[member])
			throw "Member " $ member $ " is not an object"
		}

	GetAt(member, at, defaultVal = false)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.GetAt", member, at, defaultVal)
		.ensureAllowed(member)
		.validateObject(member)
		return Suneido[member].GetDefault(at, defaultVal)
		}

	HasMember?(member)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.HasMember?", member)

		.ensureAllowed(member)
		return Suneido.Member?(member)
		}

	ensureAllowed(mem)
		{
		list = GetContributions(#ServerSuneidoWhitelist)
		Assert(list.exact.Member?(mem) or list.pattern.Any?({ mem =~ it }),
			"Found disallowed member access - " $ mem)
		}
	}
