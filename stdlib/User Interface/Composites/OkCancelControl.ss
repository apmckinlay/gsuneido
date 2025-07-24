// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// this is just the buttons, see OkCancel for the dialog function
HorzEqualControl
	{
	Name: 'OkCancel'
	Xstretch: 0

	New()
		{
		super('Fill', #(Button 'OK'), 'Skip', #(Button 'Cancel'))
		}
	}

