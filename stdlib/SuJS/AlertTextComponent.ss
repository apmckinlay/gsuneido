// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
StaticComponent
	{
	line: 1
	ConvertText(text)
		{
		metrics = SuRender().GetTextMetrics(.El, 'M')
		xmax = 400
		yThreshold = 400

		lines = .getLines(text, xmax)
		line = lines.Size()
		if metrics.height * line > yThreshold
			{
			xmax = 800
			lines = .getLines(text, xmax)
			line = lines.Size()
			}
		.line = line
		xmin = 200
		lines.Each()
			{
			xmin = Max(xmin, SuRender().GetTextMetrics(.El, it).width)
			}
		.Orig_xmin = xmin
		.Orig_ymin = metrics.height * line
		return lines.Map(XmlEntityEncode).Join('<br>').Replace(' ', '\&nbsp')
		}

	getLines(text, xmax)
		{
		lines = Object()
		for line in text.Lines()
			lines.Add(@.BestFit(line, xmax))
		return lines
		}

	}
