// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		if .invalidToken?(env)
			return Object('Unauthorized', [], 'Your session is invalid or expired')

		if not env.queryvalues.Member?(0)
			return ['BadRequest', [], 'Invalid request, missing file name']

		filename = Base64.Decode(env.queryvalues[0]).Xor(EncryptControlKey())
		preview? = env.queryvalues.GetDefault('preview', false)
		saveName = env.queryvalues.GetDefault('saveName', false)
		.Download(env, filename, preview?, saveName)
		}

	invalidToken?(env)
		{
		return JsSessionToken.Validate(env) is false
		}

	Download(env, filename, preview? = false, saveName = false)
		{
		// Todo: only allow downloading from files in temp or attachment or UserData (for Videos)
		temp? = Paths.Basename(filename) is filename
		dest = temp? ? Paths.Combine('temp', filename) : filename
		if not FileExists?(dest)
			return Object('404 Not Found', Object(), 'not found')
		headers = .buildHeaders(filename, preview?, saveName)
		result = SendFileToSocket(env.socket, dest, headers)
		if temp?
			{
			if true isnt deleteResult = DeleteFile(dest)
				SuneidoLog('ERRATIC: Could not clean up file from temp',
					params: Object(:deleteResult, :dest), calls:)
			}
		return result
		}

	buildHeaders(filename, preview?, saveName)
		{
		headers = Object()
		if preview? in (false, 'false')
			headers['Content-Disposition'] =
				'attachment; filename="' $
					Paths.Basename(saveName is false ? filename : String(saveName)) $ '"'
		if false isnt type = MimeTypes.GetDefault(filename.AfterLast('.').Lower(), false)
			{
			if preview? isnt false and type.Prefix?('text/')
				type = 'text/plain'
			headers['Content_Type'] = type
			}
		return headers
		}
	}
