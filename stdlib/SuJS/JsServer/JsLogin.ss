// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		if Sys.Win32?()
			return .notSupport

		userAgent = env.GetDefault('user_agent', 'UNKNOWN')
		context = [
			toYear: Date().Year(),
			domain: '@axonsoft.com',
			favIcon: '/Res?name=win-favicon.ico',
			manifest: '/Res?name=WindowsManifest.json',
			logo: '/Res?name=login-axonswoosh.svg',
			touchIcon: '',
			extraOnLoad: '']

		if userAgent.Has?('Mac OS X')
			{
			context.favIcon = '/Res?name=apple-favicon.ico'
			context.touchIcon = HtmlString(.touchIcon)
			context.manifest = '/Res?name=AppleManifest.json'
			}

		if .NoAuth?()
			{
			result = .loginSuccess('default',
				env.GetDefault('x_forwarded_for', env.remote_user),
				env.host,
				userAgent,
				preAuth:)
			context.extraOnLoad = HtmlString('doLogin(' $ result[2] $ ')')
			result[2] = '<!DOCTYPE html>' $ Razor(JsLoginTemplate, context)
			return result
			}

		return '<!DOCTYPE html>' $ Razor(JsLoginTemplate, context)
		}

	userTable: 'users'
	NoAuth?()
		{
		return not TableExists?(.userTable) or QueryEmpty?(.userTable)
		}

	touchIcon : `<link rel="apple-touch-icon" ` $
		`type="image/png" sizes="180x180" href="/Res?name=apple-touch-icon.png">`

	notSupport: '<!DOCTYPE html>
		<html>
			<body>
				suneido.js is not supported on Win32 server (or in standalone mode).
			</body>
		</html>'

	Auth(env)
		{
		.delay()
		remote = env.GetDefault('x_forwarded_for', env.remote_user)
		host = env.host
		try
			data = Json.Decode(env.body)
		catch (e)
			{
			SuneidoLog('ERROR: (CAUGHT) Issue during login attempt: ' $ e, params: env,
				caughtMsg: 'Stopped failed login - User sent back to login screen')
			return .error("Invalid Request, please log in again")
			}
		if not .validParams?(data)
			return .error("Invalid Request, please log in again")
		user = data.GetDefault(#user, '')
		userAgent = data.GetDefault(#user_agent, 'UNKNOWN')
		password = data.GetDefault(#password, '')
		forgotPassword = data.GetDefault(#forgotPassword, false)
		preauth = data.GetDefault(#preauth, false)

		if false isnt result = .preAuth(env, user, remote, host, userAgent, preauth)
			return result

		if '' isnt msg = .loginClass.ForgotPassword(user, forgotPassword)
			return .error(msg)

		if false is userRec = .loginClass.GetUserRec(user)
			{
			.delay(1.SecondsInMs()) // slow down brute force attack
			return .error('Incorrect user name', extra: #(focus: 'user'))
			}

		if false isnt result = .checkTwoFA(userRec, data, .hasCookie?(user, env),
			Object(:remote, :host))
			return result

		if .authUser(userRec, password) is false
			{
			.delay(1.SecondsInMs()) // slow down brute force attack
			return .error('Incorrect password', extra: #(focus: 'password'))
			}

		return .loginSuccess(user, remote, host, userAgent)
		}

	checkTwoFA(userRec, data, hasCookie?, extra)
		{
		result = .loginClass.TwoFA(userRec, data, hasCookie?, :extra)
		if String?(result)
			return .error(result)
		else if Object?(result)
			return Json.Encode([step2:].Merge(result))
		return result
		}

	preAuth(env, user, remote, host, userAgent, preauth = false)
		{
		if preauth is false
			return false
		if false is .validatePreauthToken(env)
			return .error("Session has expired, please log in again")
		else if user is ''
			return .error("Missing username, please log in again")
		else
			return .loginSuccess(user, remote, host, userAgent, preAuth:)
		}

	allowedParams: #("user", "preauth", "user_agent", "book", "code",
		"password", "newuserreset", "forgotPassword", "email", "remote", "host", "from",
		"step2", "sessionId", "msg")
	validParams?(data)
		{
		if not Object?(data)
			return false
		return data.Members().Subset?(.allowedParams)
		}

	CreatePreauthToken()
		{
		token = Base64.Encode(Display(Timestamp()))
		ServerSuneido.Add(#SuPreauthTokens, true, token)
		return token
		}

	validatePreauthToken(env)
		{
		if not env.queryvalues.Member?(#token)
			return false

		token = env.queryvalues.token
		if valid = ServerSuneido.GetAt(#SuPreauthTokens, token)
			ServerSuneido.DeleteAt(#SuPreauthTokens, token)
		return valid
		}

	authUser(userRec, password)
		{
		return password.Lower() is userRec.passhash.ToHex()
		}

	hasCookie?(user, env)
		{
		if env.Member?('cookie')
			{
			key = StringXor(user, .key).ToHex()
			return env.cookie.Has?(key $ '=')
			}
		return false
		}

	loginSuccess(user, remote, host, userAgent, extraCookies = #(), preAuth = false,
		tfaEmail= false)
		{
		if '' isnt err = .loginClass.AfterLoginSuccess(user)
			return .error(err)

		ob = SuSessionManager.CreateLoginToken(user, remote, host, userAgent, tfaEmail)
		bundles = [
			[tag: 'link', rel: "stylesheet",
				href: JsLoadRuntime.GetUrl("codemirror.css")],
			[tag: 'link', rel: "stylesheet",
				href: JsLoadRuntime.GetUrl("foldgutter.css")],
			[tag: 'script', src: JsLoadRuntime.GetUrl("codemirror_bundle.js")],
			[tag: 'script', src: JsLoadRuntime.GetUrl("su_code_bundle.js")],
			[tag: 'script', src: JsLoadRuntime.GetUrl("su_bundle.min.js")]]
		bundles.Append(IconFontHelper.GetFontStyles(JsLoadRuntime.GetUrl).
			Map({ [tag: 'style', innerText: it] }))
		set = .loginClass.GetInitSet(user)
		return ['OK',
			['Set-Cookie': Object(ob.token $ '=' $ ob.key $ '; Path=/; Secure; HttpOnly').
				MergeUnion(extraCookies)],
			Json.Encode([sources: bundles,
				onload: JsTranslate(
					'function () {
						SuInitClient(set: "' $ set $ '", token: "' $ ob.token $ '", ' $
						'preAuth: ' $ Display(preAuth) $ ')}', 'eval') $ ".call()"])]
		}

	TwoFA(env)
		{
		.delay()
		data = Json.Decode(env.body)
		user = data.GetDefault(#user, '')
		userAgent = data.GetDefault(#user_agent, 'UNKNOWN')

		if not .validParams?(data)
			return .error("Invalid Request")
		if not data.Member?(#remote) or not data.Member?(#host)
			return .error('Invalid session')

		if true isnt result = .twoFAAuth(data, user)
			return result

		password = data.GetDefault(#password, '')
		if false is userRec = .loginClass.GetUserRec(user)
			{
			.delay(1.SecondsInMs()) // slow down brute force attack
			return .error('Incorrect user name', extra: #(focus: 'user', back:))
			}

		if .authUser(userRec, password) is false and
			not .loginClass.AllowAuthFailure?(user)
			return .error('Invalid user name or password',
				extra: #(focus: 'password', back:))

		if user is 'default'
			SuneidoLog('WARNING: user logged in as default',
				params: data.Project(#host, #remote))

		return .loginSuccess(user, data.remote, data.host, userAgent,
			extraCookies: [.buildCookie(user)], tfaEmail: data.GetDefault(#email, false))
		}

	twoFAAuth(data, user)
		{
		sessionId = data.GetDefault(#sessionId, '')
		otp = data.GetDefault(#code, '').Trim()

		if TwoFAManager.Auth(user, sessionId, otp) is false
			return .error('Invalid code', extra: #(focus: 'code', form: '3'))

		TwoFAManager.Invalidate(otp)
		return true
		}

	key: 'suneido.js'
	buildCookie(user)
		{
		key = StringXor(user, .key).ToHex()
		expires = Date().Plus(days: 14).InternetFormat()
		return key $ '=' $ Base64.Encode(expires) $ '; Expires=' $ expires $
			'; Path=/login_submit; Secure; HttpOnly'
		}

	error(msg, extra = #())
		{
		return Json.Encode([err: msg].Merge(extra))
		}

	delay(ms = 500)
		{
		Thread.Sleep(ms)
		}

	getter_loginClass()
		{
		return LastContribution('JsLogin')
		}
	}
