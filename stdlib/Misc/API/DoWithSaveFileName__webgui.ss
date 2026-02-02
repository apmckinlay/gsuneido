// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
function (filter = "", hwnd/*unused*/ = false, flags/*unused*/ = false,
	title = "Save", ext = '', file = "", initialDir/*unused*/ = "",
	block = false, alert = '')
	{
	if ext is ''
		ext = filter.AfterLast('.')
	filename = SuGetTempSaveName(ext)
	if block isnt false
		try
			{
			block(filename)
			baseName = Paths.Basename(filename)
			saveName = file isnt '' ? file : baseName
			if not saveName.Has?('.')
				saveName $= Opt('.', ext)
			JsDownload.Trigger(baseName, saveName)
			}
		catch (err)
			{
			Alert((alert is '' ? 'Unable to save file' : alert) $ ': ' $ err,
				title, 0, MB.ICONWARNING)
			return false
			}

	return true
	}