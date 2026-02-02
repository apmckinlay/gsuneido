// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MultiAutoChooseControl
	{
	Name: "EmailAddresses"
	AlternateJoinChar: ';'
	New(width = 50, height = 3, mandatory = false)
		{
		super(.list, :width, :height, :mandatory, allowOther: ValidEmailAddress?)
		EmailAddresses.Ensure()
		}
	list(prefix)
		{
		return EmailAddresses.GetAddrs(prefix, limit: 20)
		}
	KillFocus()
		{
		super.KillFocus()

		.Get().
			Tr(';', ',').
			Split(',').
			Map!(#Trim).
			Filter({ it isnt '' and ValidEmailAddress?(it) }).
			Each(EmailAddresses.OutputAddr)
		}

	ValidData?(@args)
		{
		args.list = .list
		args.allowOther = ValidEmailAddress?
		super.ValidData?(@args)
		}
	}