// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Get(env)
		{
		args = Url.SplitQuery(env.query)
		if args.Member?("edit")
			return WikiEdit("edit", args.edit, env.remote_user)
		else if args.Member?("append")
			return WikiEdit("append", args.append, env.remote_user)
		else if args.Member?("find")
			return WikiFind(args.find)
		else if args.Member?("semantic")
			return WikiFind.Semantic(args.semantic)
		else
			return WikiView(args.Member?(0) ? args[0] : "StartPage")
		}
	Post(env)
		{
		args = Url.SplitQuery(env.query)
		if args.Member?("unlock")
			{
			WikiUnlock(args.unlock)
			return '<html><head>
				<meta http-equiv="Refresh" content="1;URL=Wiki?' $ args.unlock $ '">
				</head><body>Thank you!</body></html>'
			}
		else
			return WikiSave(args[0], env.body)
		}
	}