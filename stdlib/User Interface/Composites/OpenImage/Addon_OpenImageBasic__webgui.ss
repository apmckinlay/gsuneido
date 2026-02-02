// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
_Addon_OpenImageBasic
	{
	saveAsPrompt: 'Download...'

	On_Download()
		{
		if AttachmentS3Bucket() is ''
			{
			.On_Save_As()
			return
			}

		curFilename = .FullPath()
		if false is .ProcessFile(curFilename)
			return

		if false is url = OptContribution(
			"Attachment_PresignedUrl", {|@unused| false })(curFilename, download?:)
			{
			SuneidoLog('ERROR: (CAUGHT) Cannot generate presigned url', [:curFilename])
			Alert('Unable to download file', 'Download', 0, MB.ICONWARNING)
			return
			}

		SuRenderBackend().RecordAction(false, 'SuJsExecute', [#open, url])
		}
	}
