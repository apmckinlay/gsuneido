// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
// TODO: make it keyboard accessible up/down/enter
Controller
	{
	ComponentName: 'BookSearchResults'
	New(parent, fieldHwnd, results)
		{
		super(['Border', results, border: 5])
		.parent = parent
		.ComponentArgs.fieldHwnd = fieldHwnd
		}
	Commands: ( ("Close", "Escape") )
	Goto(address) // from links
		{
		.parent.Send(#Goto, address)
		.Inactivate()
		}
	closing?: false
	Inactivate()
		{
		if .closing? or .Destroyed?()
			return
		.closing? = true
		.Window.CLOSE()
		}
	}