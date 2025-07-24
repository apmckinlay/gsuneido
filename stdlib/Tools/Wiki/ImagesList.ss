// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		env = Url.Decode(env.query)
		file = env.AfterLast('/')
		path = env.BeforeLast('/')
		folder = path.AfterLast('/').Capitalize()
		path $= '/'

		dir = .getImageFiles(path)
		if dir.Empty?()
			return "<p>no photos</p>"

		return .buildImageFilesHtml(path, folder, dir, file)
		}

	getImageFiles(path)
		{
		dir = Object()
		for f in Dir(path $ "*.*").Sort!()
			if f =~ "(?i)\.(gif|jpg|jpeg|png)$"
				dir.Add(f)
		return dir
		}

	buildImageFilesHtml(path, folder, dir, file)
		{
		if false is i = dir.Find(file)
			file = dir[i = 0]

		html = ""
		html $= Xml('p', Xml('b', folder),
			align: 'center', style: 'margin-top: 0; margin-bottom: 0')

		next = Xml('b', 'next &gt;')
		prev = Xml('b', '&lt; prev')
		html $= Xml('p',
			(i <= 0 ? prev :
				Xml('a', prev, target: '_top', href: "ImagesPage?" $
					path $ dir[i - 1])) $ '&nbsp;&nbsp;\n' $
			(i + 1 >= dir.Size() ? next :
				Xml('a', next, target: '_top', href: "ImagesPage?" $
					path $ dir[i + 1])) $ '\n',
			align: 'center', style: 'margin-top: 0; margin-bottom: .5em')
		html $= '<font size=2 face=Arial>'
		for f in dir
			{
			base = f.Extract("[^.]*").Capitalize()
			html $= f is file
				? Xml('b', base)
				: Xml('a', base, target: '_top', href: "ImagesPage?" $ path $ f)
			html $= '<br />\n'
			}
		html $= '</font>'
		return Xml('html', Xml('body', html, bgcolor: 'lightblue'))
		}
	}