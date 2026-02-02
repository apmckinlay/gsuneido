// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'skip'
	ComponentName: 'Skip'
	small: 2
	medium: 7
	line: 23  // approx height of one vertical line
	New(amount = 10, small = false, medium = false, line = false)
		{
		if small
			amount = .small
		else if medium
			amount = .medium
		else if line
			amount = .line
		.ComponentArgs = Object(amount)
		}
	GetReadOnly() // read-only not applicable to skip
		{
		return true
		}
	SetReadOnly(unused) { }
	}