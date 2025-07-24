// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		name = env.path.AfterFirst('runtime/')
		headers = Object('Cache-Control': 'max-age=1209600'/*=two weeks*/)

		if false is rec = .getRec(name)
			{
			SuneidoLog("ERRATIC: JsLoadRuntime - Can't find " $ Display(name))
			return ['404' /*= not found */, #(), '']
			}

		if NotModified?(env, rec.GetDefault('lib_modified', rec.lib_committed))
			return ['304', headers, '']

		return ['200', headers, rec.text]
		}

	getRec(name)
		{
		switch (name)
			{
		case 'su_code_bundle.js':
			return [text: SuCode().CodeBundle.code,
				lib_modified: Max(
					SuCode().CodeBundle.code_built.lib_modified,
					SuCode().CodeBundle.code_built.lib_committed),
				hash: SuCode().CodeBundle.code_built.hash]
		case 'su_bundle.js', 'su_bundle.min.js', 'su_bundle.min.js.map',
			'codemirror.css', 'foldgutter.css', 'codemirror_bundle.js',
			'suneido.ttf', 'suneido2.ttf':
			return Query1Cached('imagebook', path: '/res', :name)
		default:
			return false
			}
		}

	GetUrl(name)
		{
		if false is rec = .getRec(name)
			return "/runtime/" $ name

		if rec.Member?(#hash)
			return "/runtime/" $ name $ '?id=' $ rec.hash

		return "/runtime/" $ name $ '?date=' $
			rec.GetDefault('lib_modified', rec.lib_committed).Format("yyyyMMddHHmmss")
		}

	Import(dir) // Called manually
		{
		path = Paths.Combine(dir, 'runtime')
		for file in #('su_bundle.js', 'su_bundle.min.js', 'su_bundle.min.js.map',
			'su_global_builtins.json')
			ImportSvcTableText(
				Paths.Combine(path, file), 'imagebook', '/res', quiet:, skipSize:)
		}
	}
