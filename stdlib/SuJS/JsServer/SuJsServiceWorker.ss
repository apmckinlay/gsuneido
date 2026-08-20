// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Get(@unused)
		{
		return .CallClass()
		}

	Func()
		{
		if false is template = SuJsLoadRuntime.QueryRec('sw.js')
			return ['404 Not Found', [], '']

		// Only check Object? because type checker complains.
		if not Object?(paths = GetContributions('RackRoutes'))
			paths = #()

		paths = paths.Filter({ it.GetDefault(#rally, false) }).
			Map({ it[1] }).
			Map(Display).
			Join(', ')

		content = template.text.Replace('(?q)const REWRITE_PATHS = []',
			'const REWRITE_PATHS = [' $ paths $ ']')
		return ['200 OK',
			[Content_Type: 'application/javascript', Cache_Control: 'no-cache'],
			content]
		}
	}