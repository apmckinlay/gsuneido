// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
SuBookSendEmail
	{
	ErrorMsg: 'unable to send email'
	CallClass(hwnd, from, to, mime, error_msg = '', quiet? = false, pdfNames = [],
		logError = false)
		{
		if TestRunner.RunningTests?()
			throw 'Test is sending the real email'
		result = true
		info = BookEmailInfo()
		from = .CleanupDisplayName(from)
		to = .FormatAddressesForSend(to)

		mimeMultiBuffered? = mime.Base?(MimeMultiBuffered)
		if mimeMultiBuffered?
			_bufferedSend = true

		if .stopSending?(hwnd, to, mime, mimeMultiBuffered?, quiet?)
			return false

		emailFunc = .sendViaSES?() ? .SendEmailSES : .sendEmailSmtp
		Working('Sending email...', :quiet?)
			{
			result = emailFunc(from, to, mime, info)
			}
		if result isnt true and not quiet?
			.HandleError(hwnd, result, error_msg, logError, NetworkService.IPAddrs())
		if result is true
			OptContribution('ProcessEmailSuccess', function(@unused) {})(from, to,
				mime, :pdfNames, quiet?: .setQuiet(quiet?))

		return result is true // TODO return error message
		}

	stopSending?(hwnd, to, mime, mimeMultiBuffered?, quiet?)
		{
		if quiet? or not Sys.Client?() or .fromContactAxon?(to)
			return false

		if not mimeMultiBuffered? or mime.GetAttachedFiles().Empty?()
			return false

		if '' is copyTo = OpenImageSettings.Copyto()
			return false

		copyFolder = Paths.ToLocal(Paths.Combine(copyTo, CopyFileAndAttach.SubFolder()))
		existing = CopyFileAndAttach.EnsureDirExists(copyFolder, copyTo, quiet?: false)
		if copyFolder isnt existing
			{
			if Object?(existing)
				Alert(existing.msg, title: 'Send Email', :hwnd, flags: MB.ICONWARNING)

			return true
			}

		return false
		}

	fromContactAxon?(to)
		{
		return to is BookEmailInfo().to
		}

	setQuiet(quiet?)
		{
		return quiet?
			? true
			: Sys.Client?()
				? false
				: true
		}

	sendViaSES?()
		{
		return not NetworkService.IPAddrs().Empty?()
		}

	sendEmailSmtp(from, to, mime, info) // returns true or error message
		{
		.setMimeSmtp(mime, from, to)
		ok = SmtpSendMsg(
			info.server,
			from,
			to,
			mime.ToString(),
			user: info.GetDefault(#user, false),
			password: info.GetDefault(#password, false),
			helo_arg: info.helo_arg)
		if String?(ok)
			{
			.logError("SmtpSendMsg failed - response code: " $ ok)
			return .ErrorMsg
			}
		return ok
		}

	setMimeSmtp(mime, from, to)
		{
		mime.From(from)
		mime.To(to)
		mime.Date()
		mime.Message_ID()
		}

	SendEmailSES(from, to, mime, info /*unused*/) // returns true or .ErrorMsg
		{
		if not NetworkService.RegisteredForService?()
			return NetworkService.NotRegisteredMessage

		if not String?(mime)
			.setMimeSES(mime, from, to)

		try
			{
			result = NetworkService.RunWithService()
				{ |ip|
				ok = .emailSES(from, to, mime, ip)
				ok is .ErrorMsg ? false : true
				}
			if result is false
				ok = .ServiceError
			}
		catch(err)
			{
			SuneidoLog("Error sending email: " $ err)
			// if there is an unexpected error with NetworkService, try the old way
			ip = NetworkService.IPAddress()
			ok = .emailSES(from, to, mime, ip)
			if ok is .ErrorMsg
				ok = .emailSES(from, to, mime, NetworkService.OtherIP(ip))
			}
		return ok
		}

	setMimeSES(mime, from, to)
		{
		from = from.Replace('[ ,]*$', '')
		mime.From(AmazonSES.SourceEmail(from))
		mime.Reply_To(from)
		mime.To(to)
		}

	emailSES(from, to, mime, ip) // returns true or .ErrorMsg
		{
		try
			{
			mimeString = String?(mime) ? mime : mime.ToString()
			if true is result = .forwardSendMsg(ip, from, to, mimeString)
				return true
			else if .InvalidEmailAddressResult?(result) or result =~ '^200 '
				return result

			.logError(NetworkService.LogMsg(ip, NetworkService.OtherIP(ip)),
				Object(:from, :to, :result, message: mimeString))
			return .ErrorMsg
			}
		catch (err)
			{
			if err.Prefix?(`File: can't open`)
				return .fileMissingError(err)
			mimeString = String?(mime) ? mime : mime.ToString()
			.logError(err $ NetworkService.LogMsg(ip2: NetworkService.OtherIP(ip)),
				Object(:from, :to, message: mimeString))
			return .ErrorMsg
			}
		}

	fileMissingError(err)
		{
		filePath = err.AfterFirst("can't open '").BeforeFirst("' in mode")
		fileName = Paths.Basename(filePath)
		filePathOnly = filePath.Replace(fileName, '')
		msg = ''
		if false is CheckDirExists(filePathOnly)
			msg = 'folder does not exist: ' $ filePathOnly
		else
			msg = .fileExists?(filePath)
				? fileName $ ', it exists, but was being held open and cannot attach'
				: fileName $ ', it does not exist'
		SuneidoLog('ERROR: Failed to attach ' $ msg)
		return .ErrorMsg
		}

	// extracted for tests
	fileExists?(filePath)
		{
		return FileExists?(filePath)
		}
	forwardSendMsg(ip, from, to, mimeStr)
		{
		return ForwardSendMsg(ip, from, to, mimeStr)
		}

	AlertErrorMessage(error_msg, hwnd)
		{
		Alert(error_msg, title: 'Send Email', :hwnd, flags: MB.ICONERROR)
		}

	logError(msg, params = #())
		{
		SuneidoLog('ERRATIC: BookSendEmail - ' $ msg, :params)
		}

	CreateMime(subject, message, filename, attachFileName)
		{
		return MimeMultiBuffered().
			Subject(subject).
			Attach(MimeText(message)).
			AttachFile(filename, :attachFileName)
		}
	}
