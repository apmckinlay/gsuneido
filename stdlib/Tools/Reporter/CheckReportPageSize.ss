// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(paramRpt, devMode)
		{
		if false is report = Query1("params", report: paramRpt)
			return ''

		if not report.params.Member?(#reporterCanvas)
			return ''

		canvas = report.params.reporterCanvas
		paperDefault = ReportPagePaperSpecs.Default
		width = canvas.GetDefault(#canvasWidth, paperDefault.w)
		height = canvas.GetDefault(#canvasHeight, paperDefault.h)
		if width isnt devMode.width or height isnt devMode.height
			return 'Design Tab > Canvas Size/Orientation (' $
				ReportPagePaperSpecs.GetPageSelection(width, height) $ '/' $
				canvas.GetDefault(#orientation, paperDefault.orientation) $
				') does not match Page Setup Size/Orientation'

		return ''
		}
	}