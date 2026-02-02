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

		if temp?
			.DeleteTask(filename)

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

	Trigger(filename, saveName)
		{
		// filename must be a file in temp
		Assert(Paths.Basename(filename) is filename)

		.AddTask(Suneido.User, filename, saveName)
		SuRenderBackend().RecordAction(false, 'SuDownloadFile', [
			target: Base64.Encode(filename.Xor(EncryptControlKey())),
			:saveName])
		}

	AddTask(user, filename, saveName)
		{
		if Sys.Client?()
			ServerEval('JsDownload.AddTask', user, filename, saveName)

		.ensureTasks()
		Suneido.SuJsDownloadTasks[user][filename] = Object(t: Date(), :saveName)
		}

	DeleteTask(filename)
		{
		if Sys.Client?()
			ServerEval('JsDownload.DeleteTask', filename)

		.ensureTasks()
		for user in Suneido.SuJsDownloadTasks.Members().Copy()
			Suneido.SuJsDownloadTasks[user].Delete(filename)
		}

	ensureTasks()
		{
		if not Suneido.Member?(#SuJsDownloadTasks)
			Suneido.SuJsDownloadTasks = Object().Set_default(Object())
		}

	CheckTask(user)
		{
		Assert(not Sys.Client?())

		.ensureTasks()
		result = Object()
		before = Date().Minus(seconds: 30)
		defVal = Object(t: Date.End())
		for filename in Suneido.SuJsDownloadTasks[user].Members().Copy()
			{
			ob = Suneido.SuJsDownloadTasks[user].GetDefault(filename, defVal)
			if ob.t < before
				result[filename] = ob
			}
		return result
		}

	// in thread
	WarnIfOutstanding(result)
		{
		if not Sys.SuneidoJs?()
			return

		if false isnt alert = Suneido.GetDefault(#JsDownloadAlertWindow, false)
			Defer({
				if not alert.Destroyed?()
					alert.Ctrl.Update(result)
				})
		else if result.NotEmpty?()
			Defer({ JsDownloadAlertControl(result) })
		}
	}
