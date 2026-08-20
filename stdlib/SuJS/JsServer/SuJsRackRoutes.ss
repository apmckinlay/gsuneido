// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	routes: (
		['Get', 	'/TestConnection', 		function (@unused) { return 'Okay' }],
		['Get', 	'/Res', 				'GetBookRes'],
		['GET',		`/runtime`,				'SuJsLoadRuntime'],
		['GET',		'/load$',				'SuJsLoadRecord'],
		['GET', 	'/suneidoapp', 			'SuJsSuneidoAPP.Handle'],
		['GET',		'/download',			'SuJsDownload'],
		['POST',	'/upload',				'SuJsUpload'],
		['GET',		'/attachment',			'SuJsViewAttachment'],
		['GET', 	'/$',					'SuJsLogin'],
		['POST',	'/login_submit$',		'SuJsLogin.Auth'],
		['POST',	'/twoFA_submit$',		'SuJsLogin.TwoFA'],
		['GET', 	'/connect$',			function (env) {
						WebSocketHandler(env, env.socket, SuJsWebSocketServer) }],
		['GET', 	'/robots.txt', 			'SuJsRobotsTxt']
		['GET', 	'/sw.js', 				'SuJsServiceWorker.Get']
		['GET',		'/rally',				'SuJsRally']
		['POST',	'/rally',				'SuJsRally']
	)
	CallClass()
		{
		return GetContributions('RackRoutes').
			Filter({ it.GetDefault(#public, false) }).
			Append(.routes)
		}
	}
