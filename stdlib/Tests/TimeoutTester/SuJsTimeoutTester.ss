// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
SuJsTester
	{
	CallClass(name, user, title = false)
		{
		.Start(Url.Encode('SuJsTimeoutTester', [:name, :user]), id: name, :title,
			timeout: 10)
		}

	ExtraRoutes: #(
		['GET',		'/SuJsTimeoutTester$',			'SuJsTimeoutTester.Page'])

	Page(env)
		{
		return .BasePage(env, Object('SuJsTimeoutTester.Run', args: env.queryvalues),
			JsLoadRuntime.GetUrl("su_code_bundle.js"))
		}

	Run(name, user)
		{
		type = name.BeforeFirst('_')
		types = GetContributions('TimeoutTesterTypes')
		timeoutTester = types.FindOne({ it.name is type })
		Suneido.User = 'timeout' $ user
		sessionName = Suneido.User $ '@local(jsS)'
		Thread.Name(sessionName)
		Database.SessionId(sessionName)

		.log("Client started for " $ type $ ' as ' $ Suneido.User)
		PersistentWindow.Load(timeoutTester.set)
		TimeoutClient.OpenScreen(timeoutTester.name, timeoutTester)
		.log("attempting to add client for " $ type $ ' as ' $ Suneido.User)
		ServerEval('TimeoutTester.AddStartedClient', name)

		SuRenderBackend().RegisterBeforeDisconnectFn({ ServerSuneido.Set(name, true) })

		if timeoutTester.GetDefault(#access?, true)
			{
			// +1.2 min so that Access timeout can be triggered first in timeout testers
			// Access timeout need 1 + 1 mins to trigger
			extra = 1.2
			SuRenderBackend().SetTimeoutMin(
				Database.Info().GetDefault(#timeoutMin, 240/*= 4 hrs*/) + extra)
			}
		Delay(4.MinutesInMs() /*=should have timed out*/,
			{ SuRenderBackend().DumpStatus(name $ ' not timeout after 4 mins') })
		}

	log(s)
		{
		Rlog('timeout', s $ '\r\n')
		}
	}