// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(source, root, writeMgr, registerDebuggerFn)
		{
		writer = writeMgr.GetNewWriter(root)

		_helper = Object(:source, :writer, :registerDebuggerFn, :writeMgr, skip: Object())
		map = .createMap()
		TdopTraverse2(root, { .handleNode(it, map) })

		return writer.ToString()
		}

	createMap()
		{
		map = Object()
		map[TDOPTOKEN.FUNCTIONDEF] = { |node| .handleStmtList(node[5/*=LIST{stmt}*/]) }
		map[TDOPTOKEN.BLOCK] = { |node| .handleStmtList(node[4/*=LIST{stmt}*/]) }
		map[TDOPTOKEN.SWITCHSTMT] = { |node|
			node[5/*=LIST{CASE_ELEM}*/].Children.
				Each({ .handleStmtList(it[3/*=LIST{stmt}*/]) }) }
		map[TDOPTOKEN.IFSTMT] = { |node|
			.handleStmts(node[4/*=stmts*/])
			.handleStmts(node[6/*=stmts*/]) }
		map[TDOPTOKEN.WHILESTMT] = { |node| .handleStmts(node[4/*=stmts*/]) }
		map[TDOPTOKEN.DOSTMT] = { |node| .handleStmts(node[1/*=stmts*/]) }
		map[TDOPTOKEN.FORSTMT] = { |node| .handleStmts(node[8/*=stmts*/]) }
		map[TDOPTOKEN.FORINSTMT] = { |node| .handleStmts(node[6/*=stmts*/]) }
		map[TDOPTOKEN.FOREVERSTMT] = { |node| .handleStmts(node[1/*=stmts*/]) }
		map[TDOPTOKEN.TRYSTMT] = { |node| .handleStmts(node[1/*=stmts*/]) }
		map[TDOPTOKEN.CATCHSTMT] = { |node| .handleStmts(node[2/*=stmts*/]) }
		return map
		}

	handleNode(node, map)
		{
		if node.Length is 0 or _helper.skip.Has?(node)
			return false

		if false isnt fn = map.GetDefault(node.Token, false)
			fn(node)

		return true
		}

	handleStmtList(listNode)
		{
		for i in .. listNode.ChildrenSize()
			{
			if listNode[i].Length is 0 or .isSuperNew?(listNode[i])
				continue

			num = (_helper.registerDebuggerFn)(listNode[i])
			_helper.writer.Add(listNode, i, 'SteppingDebugger(' $ num $ ');\r\n')
			}
		return true
		}

	isSuperNew?(node)
		{
		if not node.Match(TDOPTOKEN.STMT)
			return false
		node = node[0]

		if not node.Match(TDOPTOKEN.CALL)
			return false

		if not node[0].Match(TDOPTOKEN.MEMBEROP)
			return false

		return node[0][0].Match(TDOPTOKEN.SUPER) and node[0][2].Value is 'New'
		}

	handleStmts(stmtsNode)
		{
		if stmtsNode.Length is 0
			return

		if stmtsNode.Match(TDOPTOKEN.STMTS)
			{
			.handleStmtList(stmtsNode[1])
			return
			}

		num = (_helper.registerDebuggerFn)(stmtsNode)
		replace = .CallClass(_helper.source, stmtsNode,
			_helper.writeMgr, _helper.registerDebuggerFn)
		_helper.writer.Replace(stmtsNode, '{ SteppingDebugger(' $ num $ '); ' $
			replace $ ' }')
		_helper.skip.Add(stmtsNode)
		}
	}
