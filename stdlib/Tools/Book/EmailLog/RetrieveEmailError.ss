// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ErrorRegEx: '^(ERROR|ERRATIC): (BookSendEmail|ForwardSendMessage)'

	CallClass(log)
		{
		context = Object().Set_default('')
		context.date = log.sulog_timestamp.ShortDateTime()
		context.errMsg = log.sulog_message.Replace(.ErrorRegEx $ "(\W)*", "")

		locals = log.sulog_locals
		if locals.Empty?()
			{
			// locals default value is #()
			// params default value is ''
			locals = log.GetDefault('sulog_params', '')
			if locals is ''
				return context
			}

		context.from = locals.GetDefault('from', '')
		context.to = locals.GetDefault('to', '')

		if locals.Member?('ok')
			context.errMsg $= ' - ' $ locals.ok
		else if locals.Member?('response')
			context.errMsg $= ' - ' $ locals.response.GetDefault('content', '')

		if locals.Member?('message')
			context.message = locals.message
		else if locals.Member?('mime')
			context.message = locals.mime.MimeBase_payload

		return context
		}
	}