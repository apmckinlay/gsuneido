// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
CodeViewAddon
	{
	Name:		ClassOutline
	Inject:		topRight
	Controls(addonControls)
		{ addonControls.Add(#(WndPane, (Scroll, ClassOutline)), at: 10) }

	Init()
		{
		if .AddonControl isnt false
			.Set()
		}

	IdleAfterChange()
		{ .Set() }

	Set()
		{ .AddonControl.Set(.Get()) }

	Invalidate()
		{ .Set() }

	Addon_RedirMethods()
		{ return #(ClassOutline_SelectItem, Outline_Highlight) }

	ClassOutline_SelectItem(member)
		{
		.SendToAddons(#GoToDef)
		.SendToAddons(#Goto, member)
		.Defer({ .Send("SetFocus") })
		}

	Outline_Highlight(cl, name)
		{ .AddonControl.Outline_Highlight(cl, name) }
	}
