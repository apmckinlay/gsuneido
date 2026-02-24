// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_Addon_highlight_occurrences
	{
	Init() {}
	Styling() { return [] }
	UpdateUI() {}

	ComponentAddon(componentAddons)
		{
		componentAddons.Addon_highlight_occurrences_Component =
			Object(fore: .GetSchemeColor('occurrence'))
		}
	}