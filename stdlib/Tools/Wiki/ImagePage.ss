// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// FIXME: file:/// does not seem to being working anymore
function (env)
	{
	imagePath = Paths.ToStd(Url.Decode(env.query))
	if imagePath.Suffix?('/')
		{
		f = false
		for f in Dir(imagePath $ '*.*').Sort!()
			if f =~ "(?i)\.(gif|jpg|jpeg|png)$"
				break
		if not String?(f)
			return Xml('html', Xml('body', Xml('p', 'no images in ' $ imagePath)))
		imagePath $= f
		}
	return Xml('html',
		Xml('body',
			Xml('p',
				Xml('img', src: 'file:///' $ imagePath, height: '95%'), align: 'center'),
			bgcolor: 'black'))
	}