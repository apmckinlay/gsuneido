// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(data = '', .xPrompt = "Left", .yPrompt = "Top", width = false)
		{
		super(:data, :width)
		}
	GetSize(data = '')
		{
		.Data = .buildStr(data)
		return super.GetSize(.Data)
		}
	Print(x, y, w, h, data = '')
		{
		.Data = .buildStr(data)
		super.Print(x, y, w, h, .Data)
		}
	ExportCSV(data = '')
		{
		.Data = .buildStr(data)
		return .CSVExportString(.Data.Tr(','))
		}
	buildStr(data = '')
		{
		split = CoordControl.SplitCoord(data)
		if split.x is '' and split.y is ''
			return ''
		return Opt(.xPrompt $ ': ', split.x) $ ', ' $ Opt(.yPrompt $ ': ', split.y)
		}
	}
