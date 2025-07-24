// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
function (subject, message, emailFrom, emailTo, notify, logError = false)
	{
	if false is logError and
		false isnt logError = OptContribution('BookSendEmail_LogError', false)
		logError = logError(Object(:subject, :message), emailTo, notify)

	BookSendEmail(0, emailFrom, emailTo, MimeText(message).Subject(subject), :logError)
	}
