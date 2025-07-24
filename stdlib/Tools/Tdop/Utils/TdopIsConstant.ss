// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(node)
		{
		if node.Token in (TDOPTOKEN.NUMBER, TDOPTOKEN.DATE,
			TDOPTOKEN.SYMBOL, TDOPTOKEN.TRUE, TDOPTOKEN.FALSE, TDOPTOKEN.OBJECT,
			TDOPTOKEN.RECORD, TDOPTOKEN.CLASSDEF, TDOPTOKEN.FUNCTIONDEF, TDOPTOKEN.DLLDEF,
			TDOPTOKEN.STRUCTDEF, TDOPTOKEN.CALLBACKDEF)
			return true
		if node.Match(TDOPTOKEN.UNARYOP)
			return node.Children[0].Token in (TDOPTOKEN.ADD, TDOPTOKEN.SUB) and
				node.Children[1].Match(TDOPTOKEN.NUMBER)
		if .stringOrStringCat(node)
			return true
		return false
		}

	stringOrStringCat(node)
		{
		if node.Match(TDOPTOKEN.STRING)
			return true
		return node.Match(TDOPTOKEN.BINARYOP) and .stringOrStringCat(node.Children[0]) and
			node.Children[1].Match(TDOPTOKEN.CAT) and .stringOrStringCat(node.Children[2])
		}
	}
