// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.source, .astTree = false)
		{
		if .astTree is false
			.astTree = Tdop(.source)
		.init()
		}

	init()
		{
		.comments = Object()
		.blanks = Object()
		scan = Scanner(.source)
		while scan isnt scan.Next()
			{
			if scan.Type() is #COMMENT
				{
				rec = Object(
					text: scan.Text(),
					pos: scan.Position() - scan.Text().Size() + 1)
				.comments.Add(rec, at: rec.pos)
				}
			else if scan.Type() in (#NEWLINE, #WHITESPACE)
				{
				rec = Object(
					type: scan.Type(),
					text: scan.Text(),
					pos: scan.Position() - scan.Text().Size() + 1)
				.blanks.Add(rec, at: rec.pos)
				}
			}
		}

	Getter_Root()
		{
		return .astTree
		}

	GetNewWriter(rootNode = false)
		{
		return AstWriter(.source, .astTree, .comments, .blanks, rootNode)
		}

	Default(@args)
		{
		f = args[0]
		if not .Root.Member?(f)
			throw "method not found: " $ args[0]
		return .Root[f](@+1args)
		}
	}