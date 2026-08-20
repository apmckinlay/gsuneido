// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(org, fmt)
		{
		return .Eq(Suneido.Parse(org), Suneido.Parse(fmt))
		}

	Eq(org, fmt)
		{
		org = .unwrap(org)
		fmt = .unwrap(fmt)
		if Type(org) isnt #AstNode or Type(fmt) isnt #AstNode
			return .valueEq(org, fmt)
		if false isnt r = .rewriteEq(org, fmt)
			return r.result
		if org.type isnt fmt.type
			return false
		return .sameTypeEq(org, fmt)
		}

	// cross-type rewrites the formatter can make; false if none applies,
	// else [result:] with the comparison outcome
	rewriteEq(org, fmt)
		{
		if org.type is #Trinary and fmt.type isnt #Trinary
			return [result: .boolTrinaryEq(org, fmt)]
		// while true / for (;;) are rewritten to forever
		if fmt.type is #Forever and org.type isnt #Forever
			return [result: .foreverEq(org, fmt)]
		return .rewriteEq2(org, fmt)
		}

	rewriteEq2(org, fmt)
		{
		// x++ / x-- statements are rewritten to x += 1 / x -= 1
		if org.type is #Unary and fmt.type is #Binary
			return [result: .incDecEq(org, fmt)]
		// a long string may be split into a $ chain of pieces
		if org.type is #Constant and fmt.type is #Nary and String?(org.children[0])
			return [result: .splitStringEq(org, fmt)]
		return false
		}

	incDecEq(org, fmt)
		{
		if org.op not in (#PostInc, #PostDec) or
			fmt.op isnt (org.op is #PostInc ? #AddEq : #SubEq)
			return false
		return Type(fmt.rhs) is #AstNode and fmt.rhs.type is #Constant and
			fmt.rhs.value is 1 and .Eq(org.expr, fmt.lhs)
		}

	sameTypeEq(org, fmt)
		{
		switch org.type
			{
		case #Ident:
			return org.name is fmt.name
		case #Unary, #Binary:
			return .opEq(org, fmt)
		case #Nary:
			return .naryEq(org, fmt)
		case #Call:
			return .callEq(org, fmt)
		case #RangeTo, #RangeLen:
			return .rangeEq(org, fmt)
		default:
			return .compoundEq(org, fmt)
			}
		}

	// x[0 .. hi] === x[.. hi] and x[0 :: len] === x[:: len]:
	// a written 0 start-index is equivalent to an omitted one
	rangeEq(org, fmt)
		{
		if not .Eq(org.expr, fmt.expr) or not .zeroFromEq(org.from, fmt.from)
			return false
		len? = org.type is #RangeLen
		return .Eq(len? ? org.len : org.to, len? ? fmt.len : fmt.to)
		}

	zeroFromEq(x, y)
		{
		return .zeroFrom?(x) and .zeroFrom?(y) ? true : .Eq(x, y)
		}

	zeroFrom?(n)
		{
		n = .unwrap(n)
		return n is false or
			(Type(n) is #AstNode and n.type is #Constant and n.value is 0)
		}

	opEq(org, fmt)
		{
		return org.op is fmt.op and .childrenEq(org, fmt)
		}

	naryEq(org, fmt)
		{
		return org.op is fmt.op and
			(org.op is #Cat ? .catEq(org, fmt) : .childrenEq(org, fmt))
		}

	compoundEq(org, fmt)
		{
		switch org.type
			{
		case #Function, #Block:
			return .funcEq(org, fmt)
		case #TryCatch:
			return .tryCatchEq(org, fmt)
		case #ForIn:
			return .forInEq(org, fmt)
		case #Return:
			// return throw: the flag is not a child, compare it explicitly
			return org.throw is fmt.throw and .childrenEq(org, fmt)
		default:
			return .childrenEq(org, fmt)
			}
		}

	funcEq(org, fmt)
		{
		return .paramsEq(org.params, fmt.params) and .childrenEq(org, fmt)
		}

	forInEq(org, fmt)
		{
		return org.var is fmt.var and org.var2 is fmt.var2 and .childrenEq(org, fmt)
		}

	tryCatchEq(org, fmt)
		{
		return org.catchpat is fmt.catchpat and .catchVarEq(org, fmt) and
			.Eq(org.try, fmt.try) and .Eq(org.catch, fmt.catch)
		}

	foreverEq(org, fmt)
		{
		if org.type is #While
			return .trueCond?(org.cond) and .Eq(org.body, fmt.body)
		if org.type is #For
			return org.init.Empty?() and org.cond is false and org.inc.Empty?() and
				.Eq(org.body, fmt.body)
		return false
		}

	trueCond?(c)
		{
		c = .unwrap(c)
		return Type(c) is #AstNode and c.type is #Constant and c.value is true
		}

	unwrap(n)
		{
		while Type(n) is #AstNode and n.type is #Compound and false isnt n.children[0] and
			false is n.children[1]
			n = n.children[0]
		// (x) parens are value-transparent; switch (x) => switch x drops them
		while Type(n) is #AstNode and n.type is #Unary and n.op is #LParen
			n = n.expr
		return n
		}

	valueEq(org, fmt)
		{
		if Type(org) is Type(fmt)
			return org is fmt
		// splitting a string inside a folded constant unfolds it into a
		// node; evaluate the node back to a value and compare
		node = Type(org) is #AstNode ? org : fmt
		val = Type(org) is #AstNode ? fmt : org
		r = .evalc(node)
		return r.ok is true and r.val is val
		}

	evalc(n)
		{
		if Type(n) isnt #AstNode
			return [ok:, val: n]
		if n.type is #Constant
			return [ok:, val: n.children[0]]
		if n.type is #Nary and n.op is #Cat
			return .evalCat(n)
		if n.type in (#Object, #Record)
			return .evalObject(n)
		return [ok: false]
		}

	evalCat(n)
		{
		s = ""
		for (i = 0; false isnt c = n.children[i]; ++i)
			{
			r = .evalc(c)
			if r.ok isnt true or not String?(r.val)
				return [ok: false]
			s $= r.val
			}
		return [ok:, val: s]
		}

	evalObject(n)
		{
		ob = n.type is #Record ? [] : Object()
		for (i = 0; false isnt m = n.children[i]; ++i)
			{
			r = .evalc(m.value)
			if r.ok isnt true
				return [ok: false]
			if m.named
				ob[m.key] = r.val
			else
				ob.Add(r.val)
			}
		return [ok:, val: ob]
		}

	// cond ? true : false is rewritten to cond; cond ? false : true to its negation
	boolTrinaryEq(org, fmt)
		{
		t = org.t
		f = org.f
		if not .boolConst?(t) or not .boolConst?(f) or t.value is f.value
			return false
		return t.value is true ? .Eq(org.cond, fmt) : .negEq(org.cond, fmt)
		}

	boolConst?(n)
		{
		return Type(n) is #AstNode and n.type is #Constant and Boolean?(n.value)
		}

	// fmt must be the negation of cond: swapped is/isnt and =~/!~,
	// cancelled double negation, or a wrapping not (covers not in)
	negEq(cond, fmt)
		{
		cond = .unwrap(cond)
		if .notNode?(cond)
			return .Eq(cond.expr, fmt)
		if Type(fmt) isnt #AstNode
			return false
		if .notNode?(fmt)
			return .Eq(cond, fmt.expr)
		return .swappedCompareEq(cond, fmt)
		}

	notNode?(n)
		{
		return Type(n) is #AstNode and n.type is #Unary and n.op is #Not
		}

	swappedCompareEq(cond, fmt)
		{
		swap = #(Is: Isnt, Isnt: Is, Match: MatchNot, MatchNot: Match)
		if Type(cond) isnt #AstNode or cond.type isnt #Binary or fmt.type isnt #Binary
			return false
		if not swap.Member?(cond.op) or fmt.op isnt swap[cond.op]
			return false
		return .Eq(cond.lhs, fmt.lhs) and .Eq(cond.rhs, fmt.rhs)
		}

	splitStringEq(org, fmt)
		{
		joined = ""
		for (i = 0; false isnt c = fmt.children[i]; ++i)
			{
			if Type(c) isnt #AstNode or c.type isnt #Constant or
				not String?(c.children[0])
				return false
			joined $= c.children[0]
			}
		return org.children[0] is joined
		}

	// compare $ chains with adjacent string constants merged, so a
	// re-split chain still compares equal
	catEq(org, fmt)
		{
		norg = .catNorm(org)
		nfmt = .catNorm(fmt)
		if norg.Size() isnt nfmt.Size()
			return false
		for (i = 0; i < norg.Size(); ++i)
			if String?(norg[i]) or String?(nfmt[i])
				{
				if norg[i] isnt nfmt[i]
					return false
				}
			else if not .Eq(norg[i], nfmt[i])
				return false
		return true
		}

	catNorm(node)
		{
		list = Object()
		for (i = 0; false isnt c = node.children[i]; ++i)
			{
			v = .constString(c)
			if v isnt false and list.Size() > 0 and String?(list.Last())
				list[list.Size() - 1] $= v
			else
				list.Add(v is false ? c : v)
			}
		return list
		}

	// the string value if c is a string constant, else false
	constString(c)
		{
		return Type(c) is #AstNode and c.type is #Constant and String?(c.children[0])
			? c.children[0]
			: false
		}

	// argument names are not in children, compare args explicitly
	callEq(org, fmt)
		{
		if org.size isnt fmt.size or not .Eq(org.func, fmt.func)
			return false
		for (i = 0; i < org.size; ++i)
			if org[i].name isnt fmt[i].name or not .Eq(org[i].expr, fmt[i].expr)
				return false
		return true
		}

	paramsEq(org, fmt)
		{
		for (i = 0;; ++i)
			{
			x = org[i]
			y = fmt[i]
			if x is false and y is false
				return true
			if not .paramEq(x, y)
				return false
			}
		}

	paramEq(x, y)
		{
		if x is false or y is false or x.name isnt y.name or x.hasdef isnt y.hasdef
			return false
		// hasdef matches here, so only the default value is left to compare
		return x.hasdef is false or .Eq(x.defval, y.defval)
		}

	catchVarEq(org, fmt)
		{
		if .Eq(org.catchvar, fmt.catchvar)
			return true
		// catch (unused) is rewritten to plain catch
		return fmt.catchvar is false and Type(v = org.catchvar) is #AstNode and
			v.type is #Ident and v.name is #unused and org.catchpat is false
		}

	childrenEq(org, fmt)
		{
		corg = org.children
		cfmt = fmt.children
		for (i = 0;; ++i)
			{
			x = corg[i]
			y = cfmt[i]
			if x is false and y is false
				return true
			if not .Eq(x, y)
				return false
			}
		}
	}
