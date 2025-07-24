// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(filename, text, fullFileName = false)
		{
		if not fullFileName
			filename = ServerPath().Dir(filename)
		if Sys.Linux?()
			text = text.Replace('\r\n', '\n')
		PutFile(filename $ .ScriptExt(), text)
		if Sys.Linux?()
			Spawn(P.WAIT, 'chmod', '+x', filename)
		}
	ScriptExt()
		{
		return Sys.Windows?() ? '.bat' : ''
		}
	ScriptName(file)
		{
		return Paths.ToLocal('./' $ file) $ .ScriptExt()
		}
	}
