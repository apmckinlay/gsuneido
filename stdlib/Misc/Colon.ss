// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// methods for working with name:value strings
class
	{
	Name(s)
		{
		return s.BeforeFirst(':')
		}
	Value(s)
		{
		return s.AfterFirst(':').Trim()
		}
	From(name, value)
		{
		return name $ ': ' $ value
		}
	}