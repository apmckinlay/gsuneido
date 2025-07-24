// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	alertTitle: "Drag and Drop"
	ImageDropFiles(hDrop)
		{
		if .GetReorderOnly()
			return

		if .GetImageReadOnly()
			{
			.AlertInfo(.alertTitle, "This field is currently protected")
			return
			}

		files = DragQueryFileList(hDrop)
		for file in files
			if false is .checkFile(file)
				return

		if files.Size() isnt 1
			{
			if 0 is .Send("ImageDropFileList", files)
				.AlertInfo(.alertTitle, "Please drag a single file at a time")
			return
			}
		.Attach(files[0])
		}
	checkFile(file)
		{
		if .checkInvalidChar(file) is false
			return false
		if false is dir = .checkDir(file)
			return false
		if dir[0].Suffix?('/') // is a folder
			{
			.alert("This field only accepts files")
			return false
			}
		if ExecutableExtension?(dir[0])
			{
			.alert(ExecutableExtension?.InvalidTypeMsg)
			return false
			}
		return true
		}

	checkDir(file)
		{
		try
			dir = .dir(file)
		catch (e)
			{
			msg = 'The system cannot find the path specified'
			if not e.Has?(msg)
				{
				msg = 'There was a problem attaching the file'
				.log(e)
				}
			.alert(msg $ ': ' $ file)
			return false
			}
		if dir.Size() < 1
			{
			// filename with unicode could be changed to ascii through win32 api,
			// which can cause file not accessible
			.alert("The following file does not exist:\n\n" $ file $
				'\n\nPlease reselect files or ' $
				'rename the file without non-standard characters')
			return false
			}
		if dir.Size() > 1
			{
			.alert("The following file matched to more than 1 file:\n\n" $ file)
			return false
			}
		return dir
		}

	checkInvalidChar(file)
		{
		if not OpenFileName.ContainsInvalidChar?(file)
			return true

		.alert("The following file contains unsupported characters:\n\n" $ file $
			'\n\nPlease reselect files')
		return false
		}
	// extracted for testing
	dir(file)
		{
		if AttachmentS3Bucket() is ''
			return Dir(file)
		else
			return Object(file)
		}
	log(error)
		{
		prefix = error.Has?(`\\tsclient`) ? 'INFO' : 'ERROR'
		SuneidoLog(prefix $ ': (CAUGHT) Unexpected attachment failure: ' $ error,
			caughtMsg: 'user alerted: There was a problem attaching the file')
		}
	alert(msg)
		{
		.AlertInfo(.alertTitle, msg)
		}
	}
