// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// export html pages of book to single large file
class
	{
	New(.book, file, title = "")
		{
		dir = file.BeforeLast('.') $ " (images)"
		.images = '\=' $ Paths.Basename(dir)
		File(file, "w")
			{|f|
			.f = f
			hdr = '<html>
				<head>
				<title>' $ title $ '</title>
				<style>
				h1,h2,h3,h4,h5,h6 { margin-bottom: 6pt }
				p,li { margin-top: 0; margin-bottom: 6pt }
				pre { margin-left: 4ex }
				dt { font-weight: bold }
				dd { margin-bottom: 6pt }
				table { margin-bottom: 6pt }
				</style>
				</head>
				</head>' $ HtmlWrap('', book)
			f.Writeline(hdr $ (hdr.Suffix?('<body>') ? '' : '<body>'))
			.process("")
			f.Writeline("</body></html>")
			}
		EnsureDir(dir)
		dir $= '/'
		QueryApply(book $ " where path =~ '^/res\>'")
			{|x|
			PutFile(dir $ x.name, x.text)
			}
		}
	process(path) // recursive
		{
		QueryApply(.book $ " where path = " $ Display(path) $
			" sort order, name")
			{|x|
			name = x.path $ "/" $ x.name
			if (x.text.Prefix?('<'))
				.f.Writeline(Asup(.imageLinks(x.text)))
			else if x.path !~ "^/res\>" and x.text.Eval().Prefix?('<') // needs Eval
				.f.Writeline(Asup(.imageLinks(x.text.Eval()))) // needs Eval
			.process(name) // do children (if any)
			}
		}
	imageLinks(text)
		{
		return text.Replace("suneido:/" $ .book $ "/res", .images)
		}
	}