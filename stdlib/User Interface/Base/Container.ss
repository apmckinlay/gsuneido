// Copyright (C) 2001 Suneido Software Corp.  All rights reserved worldwide.
// base class for controls which contain other controls
Control
	{
	SetEnabled(enabled)
		{
		Assert(Boolean?(enabled))
		.Broadcast("SetEnabled", enabled)
		super.SetEnabled(enabled)
		}
	GetEnabled()
		{
		for (child in .GetChildren())
			if (not child.GetEnabled())
				return false
		return true
		}
	SetValid(valid)
		{
		.Broadcast("SetValid", valid)
		}
	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.Broadcast("SetVisible", visible)
		super.SetVisible(visible)
		}
	SetReadOnly(readOnly)
		{
		Assert(Boolean?(readOnly))
		.Broadcast("SetReadOnly", readOnly)
		super.SetReadOnly(readOnly)
		}
	GetReadOnly()
		{
		for (child in .GetChildren())
			if (not child.GetReadOnly())
				return false
		return true
		}
	HasFocus?()
		{
		for (child in .GetChildren())
			if (child.HasFocus?())
				return true
		return false
		}
	SetFocus()
		{
		// call SetFocus on first child that supports SetFocus
		for child in .GetChildren()
			if child.MethodClass("SetFocus") isnt Control
				{
				child.SetFocus()
				break
				}
		}
	Broadcast(@args)
		{
		method = args[0]
		for (child in .GetChildren())
			if (child.Method?(method))
				child[method](@+1 args)
		}
	Update()
		{
		for (child in .GetChildren())
			child.Update()
		super.Update()
		}
	Clear()
		{
		for (child in .GetChildren())
			child.Destroy()
		}
	Destroy()
		{
		.Clear()
		super.Destroy()
		}
	}
