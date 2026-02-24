// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_Addon_show_modified_lines
	{
	init()
		{
		}

	buildMarker(type, top/*unused*/ = '001', bottom/*unused*/ = '001', level = 0)
		{
		return [:level,
			marker: [#diffMarker, back: .GetSchemeColor(type), type: .markerType]]
		}
	}