// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
// Note: if user types a name with no extension
// 		 AND that name exists THEN no default extension is added
function (filter = "", hwnd = false, flags = false,
	title = "Save", ext = '', file = "", initialDir = "")
	{
	if filter is ""
		filter = "All Files (*.*)\x00*.*"
	filter $= '\x00\x00' // ensure filter is terminated with two nuls
	if hwnd is false
		try hwnd = _hwnd
	if flags is false
		flags = OFN.OVERWRITEPROMPT | OFN.PATHMUSTEXIST |
			OFN.HIDEREADONLY | OFN.NOCHANGEDIR
	title = TranslateLanguage(title)
	ofn = Object(structSize: OPENFILENAME.Size(), :title,
		:file, maxFile: 8000, :filter, :flags, defExt: ext, :initialDir)
	if hwnd isnt false
		ofn.hwndOwner = hwnd
	ok = CenterDialog(hwnd)
		{ GetSaveFileName(ofn) }
	return ok ? ofn.file.BeforeFirst('\x00') : ""
	}
