// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Margin: false
	Init()
		{
		.SetOption('rulers', [[
			column: .Margin isnt false ? .Margin : CheckCode.MaxLineLength,
			color: ToCssColor(.GetSchemeColor('longLineMargin')),
			lineStyle: 'solid']])
		}
	}