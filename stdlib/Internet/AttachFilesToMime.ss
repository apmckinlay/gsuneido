// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(mime, attachments, allowLinks? = false)
		{
		result = EmailAttachment.ValidateAttachments(attachments)
		contrib = GetContributions("AddEmailAttachmentLinks")
		sendAsLink? = .sendAsLink?(result, allowLinks?, contrib)
		if not .allowSend?(result, sendAsLink?)
			{
			Alert(result, 'Email', 0, MB.ICONINFORMATION)
			return false
			}
		if sendAsLink?
			{
			contrib = Global(contrib.Last())
			mime.Attach(MimeText(contrib.FileLinks(contrib.UploadFiles(attachments))))
			}
		else
			attachments.Each(mime.AttachFile)
		}

	sendAsLink?(result, allowLinks?, contrib)
		{
		return result is true
			? false
			: allowLinks? and not contrib.Empty?()
		}

	allowSend?(result, sendAsLink?)
		{
		if result is true
			return true
		if result.Has?('File size')
			return sendAsLink?
		return false
		}
	}