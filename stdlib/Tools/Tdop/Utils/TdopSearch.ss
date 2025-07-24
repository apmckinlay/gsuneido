// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(target, pattern, pos = 0, prev = false)
		{
		TdopTraverse2(target,
			{ |node|
			if true is keepTraverseDown = .needCompare(node, pattern, pos, prev)
				if ((false isnt res = .compare(node, pattern)) and
					(.checkResultWithinRange(node, pos, prev) is true))
					return Object(Object(node.Position - 1, node.Length, :node)).
						Append(res)

			keepTraverseDown
			},
			reverse: prev)
		return false
		}

	compare(node1, node2)
		{
		res = Object()
		if .matchWildcard?(node1, node2)
			return res.Add(Object(node1.Position - 1, node1.Length, node: node1))

		if not (.equal?(node1, node2) and .equalValue?(node1, node2))
			return false

		for (i = 0; i < node1.ChildrenSize(); i++)
			{
			if false is childRes = .compare(node1.Children[i], node2.Children[i])
				return false
			res.Append(childRes)
			}

		return res
		}

	matchWildcard?(node1, node2)
		{
		switch
			{
		case TdopIsExpression(node1):
			return .isWildcard?(node2)
		case node1.Token is TDOPTOKEN.KEYARG:
			return node2.Match(TDOPTOKEN.ARG) and
				.isWildcard?(node2[0])
		default:
			return false
			}
		}

	isWildcard?(node)
		{
		return node.Match(TDOPTOKEN.IDENTIFIER) and node.Value =~ '^[[:alpha:]]$'
		}

	equal?(node1, node2)
		{
		if node1.Token isnt node2.Token or
			(node1.Position is -1 and node2.Position isnt -1) or
			(node1.Position isnt -1 and node2.Position is -1) or
			node1.ChildrenSize() isnt node2.ChildrenSize()
			return false
		return true
		}

	equalValue?(node1, node2)
		{
		if node1.Member?(#Value) and not node2.Member?(#Value) or
			not node1.Member?(#Value) and node2.Member?(#Value) or
			node1.Member?(#Value) and node1.Value isnt node2.Value
			return false
		return true
		}

	needCompare(node, pattern, pos, prev)
		{
		return node.Position isnt -1 and node.Length >= pattern.Length and
			(prev is false
				? node.Position - 1 + node.Length > pos
				: node.Position - 1 < pos)

		}

	checkResultWithinRange(node, pos, prev)
		{
		if prev is false
			return node.Position - 1 >= pos
		else
			return node.Position - 1 + node.Length <= pos
		}
	}