// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		return Suneido.GetInit(.name(), { new this })
		}
	Reset() // can be called on class or instance
		{
		Suneido.Delete(.name())
		return // no return value
		}
	name()
		{
		c = Instance?(this) ? .Base() : this
		return Name(c)
		}
	ResetAll()
		{
		for item in Suneido.Members().Copy()
			if Instance?(Suneido[item]) and Suneido[item].Base?(Singleton)
				Suneido[item].Reset()
		}
	}