// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query)
		{
		WithQuery(query)
			{|q|
			return .tree("", q.Tree())
			}
//		q = Suneido.ParseQuery(query)
//		return .tree("", q)
		}
	tree(indent, q)
		{
		s = indent $ q.nrows $ ' ' $ (q.fixcost + q.varcost) $ ' ' $ q.string $ '\n'
		if q.type is 'view'
			return s $ .tree(indent $= '\t', q.source)
		switch q.nchild
			{
		case 0:
			return s
		case 1:
			return .tree(indent, q.source) $ s
		case 2:
			indent $= '\t'
			return .tree(indent, q.source1) $ s $ .tree(indent, q.source2)
			}
		}
	}