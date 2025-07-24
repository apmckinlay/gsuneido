// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
StaticComponent
	{
	line: 1
	ConvertText(text)
		{
		xmin = .Orig_xmin isnt 0 ? .Orig_xmin : 700
		lines = Object()
		for line in text.Lines()
			lines.Add(@.BestFit(line, xmin))
		.line = lines.Size()
		return lines.Join('<br>').Replace(' ', '\&nbsp')
		}

	Recalc()
		{
		metrics = SuRender().GetTextMetrics(.El, 'M')
		.Xmin = .Orig_xmin isnt 0 ? .Orig_xmin : 700
		.Ymin = .Orig_ymin isnt 0 ? .Orig_ymin : metrics.height * .line
		.SetMinSize()
		if .Orig_xmin isnt 0
			.El.SetStyle('width', .Orig_xmin $ 'px')
		if .Orig_ymin isnt 0
			.El.SetStyle('height', .Orig_ymin $ 'px')
		}
	}
