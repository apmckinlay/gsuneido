// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Get(member, defaultVal = false)
		{
		if Sys.Client?()
			return ServerEval("ServerSuneido.Get", member, defaultVal)
		return Suneido.GetDefault(member, defaultVal)
		}

	Set(member, value)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.Set', member, value)
		Suneido[member] = value
		return
		}

	DeleteMember(member)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.DeleteMember', member)
		Suneido.Delete(member)
		return
		}

	DeleteAt(member, at)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.DeleteAt', member, at)
		.validateObject(member)
		Suneido[member].Delete(at)
		}

	Add(member, value, at)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.Add', member, value, at)
		.validateObject(member)
		Suneido[member][at] = value
		return
		}
	validateObject(member)
		{
		if not Suneido.Member?(member)
			Suneido[member] = Object()
		if not Object?(Suneido[member])
			throw 'Member ' $ member $ ' is not an object'
		}

	GetAt(member, at, defaultVal = false)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.GetAt', member, at, defaultVal)
		.validateObject(member)
		return Suneido[member].GetDefault(at, defaultVal)
		}

	HasMember?(member)
		{
		if Sys.Client?()
			return ServerEval('ServerSuneido.HasMember?', member)

		return Suneido.Member?(member)
		}
	}