// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(url)
		{
		res = SuneidoAPP(url)
		return res.Prefix?("<") ? .Convert(res) : res
		}

	Convert(res)
		{
		res = res.
			Replace(`src=(['"])suneido:`, `src=\1/suneidoapp`).
			Replace(`url\((['"]?)suneido:`, `url(\1/suneidoapp`).
			Replace(`href=('(suneido:)?(.+?)'|"(suneido:)?(.+?)")`,
				{ |s|
				link = s[6/*=remove href='"*/..-1].RemovePrefix('suneido:')
				'href=' $ (link.Prefix?('http')
					? s[-1] $ link $ s[-1] $ ` target="_blank" rel="noopener noreferrer"`
					: `"javascript:suIframeSend('` $ Base64.Encode(link) $ `');void(0)"`)
				})

		return res
		}

	invalidRequest: ('/eval', '/from')
	Handle(env)
		{
		url = Url.Decode(env.path.RemovePrefix('/suneidoapp'))
		if .invalidRequest.Any?({ url.Prefix?(it) })
			return ['400', [], '']

		data = SuneidoAPP('suneido:' $ url)
		headers = url =~ SuneidoAPP.Images $ '$'
			? ['Cache-Control': 'public, max-age=604800, immutable']
			: []
		return ['200', headers, data]
		}
	}
