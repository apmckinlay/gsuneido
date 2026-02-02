// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
DelayBase
	{
	CallClass(delayMs, block)
		{
		if not Sys.MainThread?()
			throw "ERROR Delay can only be called from the main GUI thread"
		if delayMs < 100 /*=minDelay*/
			throw 'ERROR Delay minimum is 100 (ms)'
		return new this(delayMs, block)
		}
	}