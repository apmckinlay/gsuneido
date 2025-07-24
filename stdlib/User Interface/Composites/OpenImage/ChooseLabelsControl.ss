// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Getter_DialogControl()
		{
		return Object('LabelChooser', Addon_Labels.GetLabels(),	.Name,
			Addon_Labels.GetType())
		}
	}