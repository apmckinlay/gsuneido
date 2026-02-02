// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'OpenFileName'

	New(filter = '', hDrop = false, multi = false, attachment? = false)
		{
		s3? = attachment? and AttachmentS3Bucket() isnt ''
		magickJS = s3?
			? OptContribution('ExternalJsCDN', '')
			: ''
		.ComponentArgs = [.convertFilter(filter), hDrop, multi, .getLimit(), s3?,
			magickJS]
		}

	getLimit()
		{
		limit = OptContribution('AttachmentSizeLimit', function () { return 0 })()
		return limit not in ('', 0) ? limit : 10 /*mb*/
		}

	convertFilter(filter)
		{
		if filter is ''
			return ''
		filter = filter.Split('\x00').GetDefault(1, '')
		if filter is '*.*'
			return ''
		return filter.Tr(';', ',').Tr('*')
		}

	UploadFinished(saveNames)
		{
		.Window.Result(saveNames)
		}

	FileSizeOverLimit()
		{
		.AlertError('Open File', 'The selected file is too large to upload (over ' $
			Display(.getLimit()) $ 'MB)')
		.Window.Result(false)
		}

	FileNameTooLong(name, limit)
		{
		.AlertError('Open File', 'The selected file is too long ' $
			'(please limit to ' $ limit $ ' characters) : ' $ name)
		.Window.Result(false)
		}

	InvalidExtenstion(name)
		{
		.AlertError('Open File', ExecutableExtension?.InvalidTypeMsg $ '\r\n- ' $ name)
		.Window.Result(false)
		}
	}
