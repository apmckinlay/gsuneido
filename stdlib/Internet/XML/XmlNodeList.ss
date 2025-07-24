// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(children)
		{
		.list = Object().Add(@children) // or .Copy() ?
		}
	Getter_(name)
		{
		if .list.Member?(name) // handle [i]
			return .list[name]
		result = Object()
		.list.Each
			{ result.Add(@(it[name].List())) }
		return XmlNodeList(result)
		}
	List()
		{
		return .list
		}
	Text()
		{
		return .list.Map(#Text).Join()
		}
	ToString()
		{
		return .list.Map(#ToString).Join()
		}
	}