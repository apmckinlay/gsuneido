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

		SujsAdapter.CallOnRenderBackend('RecordAction', false, 'SuDownloadFile',
			[target: Base64.Encode(Paths.Basename(tempFileName).
			Xor(EncryptControlKey())),
			saveName: selectedFileName])
		}
	}
