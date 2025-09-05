// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(data, mime, history = false)
		{
		if TestRunner.RunningTests?()
			throw 'Test is sending the real email'

		info = data.info
		from = data.from.Replace('[ ,]*$', '')
		mime.From(info.sesFrom)
		mime.Reply_To(from)
		mime.To(data.to)

		.sendToService(info.service, from, data.to, info, mime,
			history: { history(data, mime) })
		}

	sendToService(ip, from, to, info, mime, history = false)
		{
		xhr = SuUI.MakeWebObject('XMLHttpRequest')
		xhr.AddEventListener('readystatechange', { |event/*unused*/|
			if xhr.readyState is 4/*=DONE*/
				{
				if xhr.status is HttpResponseCodes.OK
					{
					Print('email sent')
					EmailAttachmentComponent.CloseOverlay()
					history()
					}
				else if String(xhr.status)[0] in ('4','5') or xhr.status is 0
					{
					EmailAttachmentComponent.CloseOverlay()
					if ip is info.service
						.sendToService(info.serviceOther, from, to, info, mime, history)
					else
						.HandleError(0, xhr.status $ ' ' $ xhr.response,
							ipAddrs: [info.service, info.serviceOther])
					}
				}
			})

		xhr.Open('POST', 'https://' $ ip $ '/email')
		xhr.SetRequestHeader('X-Suneido-From', from)
		xhr.SetRequestHeader('X-Suneido-To', to)
		xhr.SetRequestHeader('Authorization', info.authentication.Authorization)
		xhr.SetRequestHeader('Content-Type', 'text/plain')
		xhr.Send(mime.ToString())
		}

	FormatAddressesForSend(to)
		{
		to = to.Tr(';', ',').Replace('[\(\[]', '{').Replace('[\)\]]', '}')
		tos = to.Split(',').Map(#Trim)
		return tos.Map!(.CleanupDisplayName).Join(',')
		}

	CleanupDisplayName(address)
		{
		if not address.Has?('<')
			return address
		displayName = address.BeforeLast('<').Tr(';:<>"\',')
		return .encode(displayName) $ '<' $ address.AfterLast('<')
		}

	encode(display)
		{
		if display =~ '^[\x20-\x7f]*$'
			return display
		display = display.Trim().ToUtf8()
		return '=?UTF-8?B?' $ Base64.Encode(display).RightTrim('=') $ '?='
		}

	CreateMime(subject, message, filename, attachFileName)
		{
		return MimeMultiPart().Subject(subject).
			Attach(MimeText(message)).
			AttachFile(attachFileName, :attachFileName, fileContent: filename)
		}

	ServiceError: 'Email service is temporarily unavailable\n' $
		'Please try again in 5 minutes '
	HandleError(hwnd, result, error_msg = '', logError = false, ipAddrs = #())
		{
		if .InvalidEmailAddressResult?(result)
			error_msg = result.AfterFirst(' ')
		else if result =~ '^200 ' and not ipAddrs.Empty?()
			error_msg = 'The following services are blocked by your network settings:\n' $
				ipAddrs.Map({ '    https://' $ it $ '\n' }).Join() $
				'Please contact your system administrator'
		else if result is .ServiceError
			error_msg = result
		else if error_msg is ''
			error_msg = 'There was a problem sending the email.'
		if logError is false
			.AlertErrorMessage(error_msg, hwnd)
		else
			logError(error_msg)
		}

	AlertErrorMessage(error_msg, hwnd /*unused*/)
		{
		SuRender().Event(false, 'Alert',
			Object(error_msg, 'Email Attachment: Not Sent', flags: MB.ICONERROR))
		}

	InvalidEmailAddressResult?(result)
		{
		return result =~ "^400 invalid email address format" or
			result =~ "^403 blacklisted address"
		}
	}
