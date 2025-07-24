// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
LinkControl
	{
	Name: 'MailLink'
	Prefix: 'mailto:'
	Status: 'e.g. jones@mail.com, double click to send an email'
	New(width = 25, set = false, mandatory = false, tabover = false, hidden = false)
		{
		super(:width, :set, :mandatory, :tabover, :hidden)
		.AddContextMenuItem("", "")
		.AddContextMenuItem("Send Email To", .GoToLink)
		}
	ValidLink?(s)
		{
		return ValidEmailAddresses?(s)
		}
	}
