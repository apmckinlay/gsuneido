// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MaxStrSize: 1000
	Separator: "<#booklogsSeparator#>"
	CallClass(s, params = '', name = 'book', systemLog = false)
		{
		if .logDisabled?(s, systemLog)
			return

		params = LogFormatEntry(params, maxStrSize: .MaxStrSize)
		if Object?(params)
			params = Display(params)

		msg = '%m' $ .clientMemory() $ '\t%s\t' $ Display(s.Ellipsis(.MaxStrSize)) $
			(params isnt '' ? '\t' $ params : '') $
			.Separator $ SuneidoVersion()

		logPath = Paths.ParentOf(GetContributions('LogPaths').GetDefault('booklog', ''))
		Rlog(Opt(logPath, '/') $ name, msg)
		}

	logDisabled?(s, systemLog)
		{
		return ((systemLog isnt true and Suneido.User is 'default') or s is '' or
			TestRunner.RunningTests?())
		}

	clientMemory()
		{
		if not Sys.Client?()
			return ''

		return '/' $ Rlog.MemoryArena(false) $ ' : Res: ' $ ResourceCounts().Sum()
		}
	}