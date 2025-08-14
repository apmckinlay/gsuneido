// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		if JsSessionToken.Validate(env) is false
			return #('BadRequest', 'Your session is invalid or expired')
		if not env.queryvalues.Member?(0)
			return #('BadRequest', 'missing argument')
		fileName = String(env.queryvalues[0]).Tr(Paths.InvalidChars)
		if ExecutableExtension?(fileName)
			return ['BadRequest', ExecutableExtension?.InvalidTypeMsg]

		if env.queryvalues.GetDefault('s3', 'false') is 'true'
			{
			url = OptContribution(
				"Attachment_PresignedUrl", {|@unused| false })(fileName, method: 'PUT')
			Assert(url isnt false)
			return url
			}

		EnsureDir('temp')
		ts = Display(Timestamp())[1..]
		saveName = Paths.Combine(Paths.Combine('temp', ts), fileName)
		EnsureDirectories(saveName)
		try
			File(saveName, 'w')
				{ |f|
				env.socket.CopyTo(f, env.GetDefault(#content_length, false))
				}
		catch (e)
			{
			SuneidoLog('INFO: JsUpload: ' $ e, params: [:fileName, :saveName])
			return #('BadRequest', 'Upload file failed')
			}
		return saveName
		}
	}