// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
// This class is for calling Suneido.js' special built-in functions
// Referencing them directly in Suneido IDE will cause compile errors
// - GetCurrentWindow()
// - GetCurrentDocument()
class
	{
	GetCurrentWindow()
		{
		return Global(#GetCurrentWindow)()
		}
	GetCurrentDocument()
		{
		return Global(#GetCurrentDocument)()
		}
	GetContentWindow(iframe)
		{
		return Global(#GetContentWindow)(iframe)
		}
	GetCodeMirror()
		{
		return Global(#GetCodeMirror)()
		}

	Open(url)
		{
		try
			.GetCurrentWindow().Open(url)
		catch (e/*unused*/)
			SuRender().Event(false, 'Alert', Object(msg:"Cannot open url: " $ url,
				title:"Invalid URL", flags: MB.ICONERROR))
		}

	Default(@args)
		{
		return Global(args[0])(@args[1..])
		}
	}
