// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function(hwnd, value, filter, status, usesubfolder = false)
	{
	settings = LastContribution(#OpenImageSetting)()
	if settings is false or settings.GetDefault('stopExecutable?', false) is false
		return ToolDialog(hwnd,
			Object(OpenImageSelect, value, filter, status, usesubfolder))

	selectedfile = OpenFileName(:hwnd, title: 'Select an Attachment', :filter,
		file: value, attachment?:)
	if ExecutableExtension?(selectedfile)
		{
		Alert(ExecutableExtension?.InvalidTypeMsg, "Select Attachment",
			hwnd, MB.ICONINFORMATION)
		return false
		}

	return OpenImageSelect.ResultFile(selectedfile,	settings.normally_linkcopy,	true)
	}