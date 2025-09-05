// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Start(endpoint, id, title = 'Tester', timeout = 5)
		{
		w = WebBrowserControl( 'http://localhost:' $ .port() $ '/' $ endpoint, :title)
		w.Ctrl.OnLoad({ .Log('Host client loaded for ' $ id $
			' (' $ w.Ctrl.LocationURL $ ')') })
		w.Ctrl.OnNavComplete({ .Log('Host client nav completed for ' $ id $
			' (' $ w.Ctrl.LocationURL $ ')') })
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
		routes = Object(
			['POST',	'/SuJsTesterLog$', 		Name(this) $ '.LogFromWeb'])
		Thread({
			RunSuJSHttpServer(port: .port(), extraRoutes: routes.Append(.ExtraRoutes))})
		}

	BasePage(env, run, codeBundleUrl)
		{
		user = 'default'
		id = env.queryvalues.GetDefault(#name, '')

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
		onload = HtmlString(
			'const id = "' $ Opt('[', id, ']: ') $ '";' $
			.errorHandler $
			JsTranslate(
				'function () {
					SuInitClient(set: "' $ set $ '", token: "' $ ob.token $ '", ' $
					'preAuth: ' $ Display(preAuth) $ ')}', 'eval') $ ".call()")

		context = [:stylesheets, :scripts, :styles, :onload]
		return ['OK',
			['Set-Cookie': Object(ob.token $ '=' $ ob.key $ '; Path=/; Secure; HttpOnly')]
			'<!DOCTYPE html>' $ Razor(SuJsWebTesterTemplate, context)]
		}

	errorHandler: `
const methods = ["log", "warn", "error", "info", "debug"];
const originalMethods = {};
methods.forEach(method => {
	originalMethods[method] = console[method];

	console[method] = function (...args) {
		sendLog({
			method: method,
			args: args.map(
				a => (typeof a === "object" ? JSON.stringify(a) : a)).join(" ")});

		// Call original method to keep normal console output
		originalMethods[method].apply(console, args);
	};
});

// Generic error handler for runtime errors
window.onerror = function (message, source, lineno, colno, error) {
	sendLog({
		type: "runtime-error",
		message: message,
		source: source,
		line: lineno,
		column: colno,
		stack: error?.stack || null
	});
};

// Handler for unhandled Promise rejections
window.addEventListener("unhandledrejection", function (event) {
	sendLog({
		type: "unhandled-rejection",
		message: event.reason?.message || String(event.reason),
		stack: event.reason?.stack || null
	});
});

// Function to send logs to /SuJsTesterLog
function sendLog(data) {
	fetch("/SuJsTesterLog", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: id + JSON.stringify(data)
	}).catch(err => {});
}
`

	LogFromWeb(env)
		{
		.Log(env.body)
		return ''
		}

	LogFileName: false
	Log(s)
		{
		if .LogFileName isnt false
			Rlog(.LogFileName, s $ '\r\n')
		}
	}
