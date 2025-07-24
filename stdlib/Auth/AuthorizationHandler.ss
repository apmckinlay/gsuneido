// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Handles authorizing with token on command line or default to Login
// Also provides AddTokenToCmdLine
class
	{
	CallClass(cmdline)
		{
		// always remove token from command line, whether needed or not
		file = cmdline.Extract(pat = `^t@(\S+)`) // t@<filename>
		token = file isnt false
			? .getFile(file)
			: cmdline.Extract(pat = `^t:(\S+)`) // t:<token>
		if token isnt false
			cmdline = cmdline.Replace(pat).Trim()
		if .authorized?()
			{
			if token isnt false
				.alert("unnecessary command line authorization")
			return cmdline
			}
		if token is false
			{
			return cmdline.Has?('_UpdateClientCopies') and cmdline.Has?('shortcutName')
				? cmdline
				: "Login(origCmd:" $ Display(cmdline) $ ")"
			}
		Suneido.defaultLoginReason = 'token supplied'
		if .authorize(.decodeToken(token))
			return cmdline
		.fatal("command line authorization failed")
		}

	getFile(file)
		{
		return (false is s = GetFile(file)) ? "" : s
		}

	authorized?()
		{
		try
			ServerEval('Date')
		catch (unused, "not authorized")
			return false
		return true
		}

	decodeToken(token)
		{
		try
			return Base64.Decode(token)
		catch
			return ""
		}

	authorize(token) // overridden by test
		{
		return Database.Auth(token)
		}

	alert(msg) // overridden by test
		{
		Alert(msg)
		}

	fatal(msg) // overridden by test
		{
		Fatal(msg)
		}

	AddTokenToCmdLine(cmd)
		{
		return ' t:' $ Base64.Encode(.token()) $ Opt(' ', cmd.Trim())
		}

	token() // overridden by test
		{
		return Database.Token()
		}
	}
