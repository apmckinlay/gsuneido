// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	routes: (
		['Get', 	'/TestConnection', 		function (@unused) { return 'Okay' }],
		['Get', 	'/Res', 				'GetBookRes'],
		['GET',		`/runtime`,				'JsLoadRuntime'],
		['GET',		'/load$',				'JsLoadRecord'],
		['GET', 	'/suneidoapp', 			'JsSuneidoAPP.Handle'],
		['GET',		'/download',			'JsDownload'],
		['POST',	'/upload',				'JsUpload'],
		['GET',		'/attachment',			'JsViewAttachment'],
		['GET', 	'/$',					'JsLogin'],
		['POST',	'/login_submit$',		'JsLogin.Auth'],
		['POST',	'/twoFA_submit$',		'JsLogin.TwoFA'],
		['GET', 	'/connect$',			function (env) {
						WebSocketHandler(env, env.socket, JsWebSocketServer) }],
		['GET', 	'/robots.txt', 			'JsRobotsTxt']
	)
	CallClass()
		{
		return GetContributions('RackRoutes').
			Filter({ it.GetDefault(#public, false) }).
			Append(.routes)
		}
	}