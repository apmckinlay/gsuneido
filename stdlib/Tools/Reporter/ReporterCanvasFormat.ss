// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
QueryFormat
	{
	Name: 'ReporterCanvas'
	New(@args)
		{
		super(@.setup(args))
		}

	setup(args)
		{
		if not args.Member?(#Sf)
			{
			r = ReporterModel(args[0], defaultMode: 'form').Report(_report.Params)
			r[0].Delete(0) // 'ReporterCanvasFormat' member
			return r[0]
			}

		return args
		}

	Output()
		{
		items = DrawControl.BuildItems(.Data.reporterCanvas)
		cols = .Columns.Map({ it.text })
		ForEachDAF(items)
			{
			if cols.Find(it.GetField()) is false
				it.NoPaint? = true
			}
		return Object("DrawCanvasRow", CanvasGroup(items),
			canvasWidth: .Data.reporterCanvas.GetDefault(#canvasWidth, false),
			canvasHeight: .Data.reporterCanvas.GetDefault(#canvasHeight, false))
		}

	AfterOutput(data /*unused*/)
		{
		return 'pg'
		}
	}
