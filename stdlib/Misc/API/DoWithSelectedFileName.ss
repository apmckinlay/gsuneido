// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(selectedFileName, block = false)
		{
		tempFileName = Sys.SuneidoJs?()
			? SujsAdapter.CallGlobal('SuGetTempSaveName')
			: selectedFileName

		if block isnt false
			Finally({ block(tempFileName) }, { .Finish(selectedFileName, tempFileName) })

		return tempFileName
		}

	Finish(selectedFileName, tempFileName)
		{
		if not Sys.SuneidoJs?()
			return

		JsDownload.Trigger(Paths.Basename(tempFileName), selectedFileName)
		}
	}
