// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	// WARNING: this is an example of warning comments in todo panel
	// ERROR: this is an example of error comments in todo panel
	warningIndicator: 	false
	errorIndicator: 	false
	warnLevel: 			88
	errorLevel: 		89
	Init()
		{
		.warningIndicator = .IndicatorIdx(level: .warnLevel)
		.errorIndicator = .IndicatorIdx(level: .errorLevel)
		}

	Styling()
		{
		warnColor = .GetSchemeColor('warning')
		errorColor = .GetSchemeColor('error')
		return [
			[level: .warnLevel, indicator: [INDIC.TEXTFORE, fore: warnColor]],
			[level: .errorLevel, indicator: [INDIC.TEXTFORE, fore: errorColor]]]
		}

	UpdateUI()
		{
		.ClearIndicator(.warningIndicator)
		.ClearIndicator(.errorIndicator)
		pos = 0
		for line in .Get().Lines()
			{
			lineLength = line.Size()
			if line =~ '(?i)warning\s*(:|-)'
				.SetIndicator(.warningIndicator, pos, line.Size())
			if line =~ '(?i)error\s*(:|-)'
				.SetIndicator(.errorIndicator, pos, line.Size())
			pos += lineLength + 2 // \r + \n
			}
		}
	}
