// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: VirtualListViewControl

	Getter_Model()
		{
		return .GetModel()
		}

	Getter_Grid()
		{
		return .GetGrid()
		}

	Getter_Thumb()
		{
		return .GetViewControls().thumb
		}

	Getter_Header()
		{
		return .GetViewControls().header
		}

	Getter_ExpandBar()
		{
		return .GetViewControls().expandBar
		}

	Getter_Scroll()
		{
		return .GetViewControls().scroll
		}
	}