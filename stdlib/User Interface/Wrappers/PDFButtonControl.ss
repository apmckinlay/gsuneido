// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MenuButtonControl
	{
	savePrompt: 'Save to file...'
	New(buttonText = 'PDF', command = false)
		{
		super(buttonText, Object(.savePrompt, 'Email as attachment...'), :command)
		}
	}
