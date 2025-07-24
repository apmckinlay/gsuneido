// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(set, isHttps? = false, userAgent = 'unknown', preAuth = false)
		{
		.initFont()
		Suneido.isHttps? = isHttps?

		if false isnt run = Suneido.GetDefault(#run, false)
			Object?(run) ? Global(run[0])(@run.args) : (Global(run))()
		else
			{
			if LastContribution('JsLogin').PostLogin(:set, :userAgent, :preAuth) is false
				.exit()

			PersistentWindow.Load(set)
			if not JsLogin.NoAuth?()
				Login.PostLoginPlugins(origCmd: '')
			try Query1('postinit').text.Eval() // needs Eval
			}
		}

	initFont()
		{
		SetGuiFont([fontPtSize: StdFontsSize.DefaultSize,
			lfFaceName: StdFonts.Ui(),
			lfHeight: -12,
			lfItalic: 0,
			lfWeight: FW.NORMAL])
		Suneido.stdfont = Suneido.logfont
		}

	FromServer(info)
		{
		name = info.user $ '@' $ info.remote $ '<' $ info.token $ '>(jsS)'
		Thread.Name(name)
		Database.SessionId(name)
		.setSuneido(info)
		.localeSettings()
		}

	setSuneido(info)
		{
		Sys.Init(suneidojs:)
		user = info.user
		Suneido.JsConnectionHost = info.host
		Suneido.Merge(info)

		if user is 'default'
			{
			Suneido.User = Suneido.User_Loaded = "default"
			Suneido.user_roles = #('admin')
			}
		else
			{
			if false is Login.SetUser(user)
				.exit()
			}

		Suneido.start_time = Date()
		Suneido.Language = #(name: "english", charset: "DEFAULT", dict: "en_US")

		// for TranslateLanguage Cache
		Suneido.CacheLanguage = ""
		Suneido.pdc = NULL
		}

	localeSettings()
		{
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")
		Settings.Set('SystemShortDateFormat', "yyyy-MM-dd")
		Settings.Set('LongDateFormat', "dddd, MMMM dd, yyyy")
		Settings.Set('TimeFormat', "h:mm tt")
		Settings.Set('ThousandsSeparator', ",")
		}

	exit()
		{
		SuRenderBackend().Terminate(reason: '')
		}
	}
