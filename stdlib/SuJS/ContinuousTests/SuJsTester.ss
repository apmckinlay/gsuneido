// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Start(endpoint, id, title = 'Tester', timeout = 5)
		{
		WebBrowserControl( 'http://localhost:' $ .port() $ '/' $ endpoint, :title)
		expire = Date().Plus(minutes: timeout)
		.wait(expire, id)
		}

	wait(expire, id)
		{
		if ServerSuneido.Get(id, false) is true or Date() > expire
			ExitClient()
		Delay(100/*=.1 sec*/, { .wait(expire, id) })
		}

	port()
		{
		return ServerPort() + 100 + 1/*=extra*/
		}

	ExtraRoutes: #()
	RunSuJsHttpServer()
		{
		SuCode().CodeBundle // build the code bundle
		SuSessionLog.Ensure()
		Thread({
			RunSuJSHttpServer(port: .port(), extraRoutes: .ExtraRoutes)})
		}

	BasePage(env, run, codeBundleUrl)
		{
		user = 'default'

		remote = env.GetDefault('x_forwarded_for', env.remote_user)
		host = env.host
		userAgent = 'UNKNOWN'
		ob = SuSessionManager.CreateLoginToken(user, remote, host, userAgent, :run)

		stylesheets = [
			JsLoadRuntime.GetUrl("codemirror.css"),
			JsLoadRuntime.GetUrl("foldgutter.css")]
		scripts = [
			JsLoadRuntime.GetUrl("codemirror_bundle.js"),
			codeBundleUrl,
			JsLoadRuntime.GetUrl("su_bundle.min.js")]
		styles = IconFontHelper.GetFontStyles(JsLoadRuntime.GetUrl).Map(HtmlString)

		preAuth = true
		set = LastContribution('JsLogin').GetInitSet(user)
		onload = HtmlString(JsTranslate(
			'function () {
				SuInitClient(set: "' $ set $ '", token: "' $ ob.token $ '", ' $
				'preAuth: ' $ Display(preAuth) $ ')}', 'eval') $ ".call()")

		context = [:stylesheets, :scripts, :styles, :onload]
		return ['OK',
			['Set-Cookie': Object(ob.token $ '=' $ ob.key $ '; Path=/; Secure; HttpOnly')]
			'<!DOCTYPE html>' $ Razor(SuJsWebTesterTemplate, context)]
		}
	}
