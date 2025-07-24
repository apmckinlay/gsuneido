// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (node, indent = "", ast = false)
	{
	if ast is false
		ast = node
	if Type(node) isnt #AstNode
		Print(indent $ Display(node))
	else
		{
		extra = ""
		if false isnt pos = node.pos
			extra $= " " $ pos
		if false isnt end = node.end
			extra $= " " $ end
		s = Display(node)
		if s.Size() > 20
			s = ""
		if false isnt e = ast.Extra(node)
			s = Display(e) $ " " $ s
		Print(indent $ node.type.Upper() $ extra, s.Tr('\r\n', ' '))
		for (i = 0; false isnt c = node.children[i]; ++i)
			AstTraverse(c, indent $ "\t", ast) // RECURSE
		}
	}