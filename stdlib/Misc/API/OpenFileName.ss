// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
class
	{
	CallClass(filter = "", hwnd = false, flags = false,
		multi = false, title = "Open", file = "", initialDir = "")
		{
		if filter is ""
			filter = "All Files (*.*)\x00*.*"
		filter $= '\x00\x00' // ensure filter is terminated with two nuls
		if hwnd is false
			try hwnd = _hwnd
		// Contrary to official MSDN documentation at http://goo.gl/xkmEAu, it
		// appears the OFN_NOCHANGEDIR flag has some effect in GetOpenFileName(...)
		// (see comments at same URL).
		flags = .flags(flags, multi)
		title = TranslateLanguage(title)
		ofn = Object(structSize: OPENFILENAME.Size(), :title,
			:file, maxFile: 8000, :filter, :flags, :initialDir)
		if hwnd isnt false
			ofn.hwndOwner = hwnd
		ok = CenterDialog(hwnd)
			{ GetOpenFileName(ofn) }

		if not ok
			return multi ? #() : ""

		files = multi ? .multipleFiles(ofn.file) : Object(ofn.file.BeforeFirst('\x00'))
		validFiles = .validateFiles(files, hwnd)
		return multi ? validFiles : validFiles.GetDefault(0, "")
		}

	flags(flags, multi)
		{
		if flags is false
			flags = OFN.FILEMUSTEXIST | OFN.PATHMUSTEXIST |
				OFN.HIDEREADONLY | OFN.NOCHANGEDIR
		if multi is true
			flags |= OFN.ALLOWMULTISELECT | OFN.EXPLORER
		return flags
		}

	multipleFiles(s)
		{
		list = s.Split('\x00').Remove('')
		if list.Size() > 1
			{
			dir = list.PopFirst()
			for i in list.Members()
				list[i] = Paths.Combine(dir, list[i])
			}
		return list
		}

	validateFiles(files, hwnd)
		{
		validFiles = Object()
		invalidFiles = Object()
		for file in files
			{
			if .ContainsInvalidChar?(file)
				invalidFiles.Add(file)
			else
				validFiles.Add(file)
			}

		if invalidFiles.NotEmpty?()
			Alert('Following files contain unsupported characters ' $
				'and will be ignored\r\n' $ invalidFiles.Join('\r\n'),
				'Invalid file name', hwnd, MB.ICONWARNING)

		return validFiles
		}

	ContainsInvalidChar?(file)
		{
		return file.Has?('?') or file.Tr('\x01-\x7f') isnt ''
		}
	}
