// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
DelayBase
	{
	CallClass(block)
		{
		if Sys.MainThread?()
			throw "ERROR RunOnGui can only be used from other threads"
		return new this(0, block)
		}
	}