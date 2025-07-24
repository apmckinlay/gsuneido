// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (env)
	{
	env = env.query
	Xml('html'
		Xml('frameset',
			Xml('frame', scrolling: 'auto', src: 'ImagesList?' $ env) $
				Xml('frame', scrolling: 'auto', src: 'ImagePage?' $ env),
			frameborder: 'no',cols: '20%,80%'))
	}
