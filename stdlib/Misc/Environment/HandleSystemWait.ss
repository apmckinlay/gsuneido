// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass (cmd, wait? = false)
		{
		return .linux?()
			? wait? ? cmd : cmd $ ' &'
			: wait? ? 'start /w ' $ cmd : 'start ' $ cmd
		}

	linux?() // overridden by tests
		{
		return Sys.Linux?()
		}
	}

