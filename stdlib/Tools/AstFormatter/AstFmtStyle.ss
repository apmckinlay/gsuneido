// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.src, .cm)
		{
		}

	Unbrace(node, guard)
		{
		if false is stmt = .strippable(node)
			return false
		return guard isnt false and .captures?(stmt, guard) ? false : stmt
		}

	strippable(node)
		{
		if Type(node) isnt #AstNode or node.type isnt #Compound or node.size isnt 1
			return false
		stmt = node[0]
		if Type(stmt) isnt #AstNode or stmt.type is #Compound
			return false
		if stmt.pos is false or stmt.end in (0, false) or node.end in (0, false)
			return false
		if .cm.CommentIn?(stmt.end, node.end)
			return false
		return stmt
		}

	ElseChain(f, cur)
		{
		return false isnt (inner = .strippable(f)) and inner.type is #If and
			not .cm.CommentIn?(cur.done, inner.pos)
			? inner
			: false
		}

	captures?(stmt, guard)
		{
		forever
			{
			if Type(stmt) isnt #AstNode
				return false
			switch stmt.type
				{
			case #If:
				if stmt.f isnt false
					stmt = stmt.f
				else if guard is #Else
					return true
				else
					stmt = stmt.t
			case #TryCatch:
				if stmt.catch isnt false
					stmt = stmt.catch
				else if guard is #Catch
					return true
				else
					stmt = stmt.try
			case #While, #Forever, #ForIn, #For:
				stmt = stmt.body
			case #Compound:
				if false is stmt = .strippable(stmt)
					return false
			default:
				return false
				}
			}
		}

	// string literals: 'c' for single characters, #word for identifier-like
	// strings, "..." for everything else; swap the quote kind rather than
	// escape; backquotes (raw/regex), multiline strings, and anything already
	// needing escapes stay as written
	Quote(x, s, sym = true) // x = the value, s = the literal as written
		{
		if s[0] is '`' or s.Has?('\n') or not .Plain?(x)
			return s
		if x.Size() is 1
			return x is "'" ? '"' $ x $ '"' : "'" $ x $ "'"
		if sym is true and .BareWord?(x) and x[0] isnt '_'
			return '#' $ x
		if not x.Has?('"')
			return '"' $ x $ '"'
		if not x.Has?("'")
			return "'" $ x $ "'"
		return s
		}

	Plain?(x)
		{
		return x.Tr(" -~") is "" and not x.Has?('\\')
		}

	BareWord?(x)
		{
		return String?(x) and x.Identifier?() and
			x not in (#true, #false, #function, #class, #dll, #struct, #callback)
		}

	Shorthand?(arg)
		{
		return arg.name isnt false and arg.expr.type is #Ident and
			arg.expr.name.LocalName?() and arg.expr.name is arg.name
		}

	CatchVar(node)
		{
		if false is var = node.catchvar
			return false
		return var.type is #Ident and var.name is #unused and node.catchpat is false
			? false
			: var
		}

	SuperNew(node)
		{
		return #Mem is (fn = node.func).type and #Ident is (expr = fn.expr).type and
			expr.name is #super and #Constant is (meth = fn.mem).type and
			meth.value is #New
		}

	NewThis(node)
		{
		return #Mem is (fn = node.func).type and #Ident is (expr = fn.expr).type and
			expr.name is #this and #Constant is (meth = fn.mem).type and
			meth.value is "*new*"
		}

	UseBrackets(node)
		{
		if node.func.type isnt #Ident
			return false
		if node.func.name is #Record
			return true
		if node.func.name isnt #Object
			return false
		arg = node[0]
		return arg isnt false and arg.name is false
		}

	Vertical?(node)
		{
		if false is node.children[1]
			return false
		if node.pos in (0, false)
			return false
		prev = node.pos
		for (i = 0; false isnt m = node[i]; ++i)
			{
			if m.pos in (0, false) or m.end in (0, false)
				return false
			if not .src[prev .. m.pos].Has?('\n')
				return false
			prev = m.end
			}
		return true
		}

	AlignWidth(node)
		{
		w = 0
		for (i = 0; false isnt m = node[i]; ++i)
			{
			if m.named isnt true or not .BareWord?(m.key)
				return false
			if Type(m.value) is #AstNode
				return false
			w = Max(w, m.key.Size())
			}
		return w
		}

	SwitchExpr(e)
		{
		while Type(e) is #AstNode and e.type is #Unary and e.op is #LParen and
			e.expr.type in (#Ident, #Constant, #Mem, #Call)
			e = e.expr
		return e
		}

	BoolConst?(n, v)
		{
		return Type(n) is #AstNode and n.type is #Constant and n.value is v
		}

	ZeroIdx?(n)
		{
		return n isnt false and Type(n) is #AstNode and n.type is #Constant and
			n.value is 0
		}

	Negatable?(expr)
		{
		return (expr.type is #Binary and expr.op in (#Is, #Isnt, #Match, #MatchNot)) or
			expr.type is #In or (expr.type is #Unary and expr.op is #Not)
		}

	// debug statements go at the left margin
	Debug?(expr)
		{
		return expr.type is #Call and expr.func.type is #Ident and
			expr.func.name in (#Print, #TracePrint, #ServerPrint, #StackTrace, #Inspect,
				#TraceCallStack)
		}

	MemberSrc(m, keyed = false)
		{
		if m.pos in (0, false) or m.end in (0, false)
			return false
		s = .ValueSrc(.src[m.pos .. m.end])
		if keyed is false
			return s
		return String?(m.key) and m.key.Identifier?()
			? .ValueSrc(s.AfterFirst(':'))
			: false
		}

	ValueSrc(s)
		{
		s = s.Trim()
		if s.Suffix?(',')
			s = s[..-1].Trim()
		return s
		}

	Simple(expr)
		{
		return expr is false or expr.type in (#Constant, #Ident)
		}

	ArgComma?(a, b)
		{
		// 'f(0, :a)' without its comma reparses as 'f(0: a)' (this is a bug if unhandled)
		if a.name isnt false or b.name is false or b.pos in (0, false) or
			.Shorthand?(b) or not String?(b.name) or not b.name.Identifier?()
			return true
		i = b.pos - 1
		while i >= 0 and .src[i] in (' ', '\t', '\n', '\r')
			--i
		return i < 0 or .src[i] is ','
		}

	BracketTable?(node)
		{
		return false isnt node.pos and .src[node.pos] is '['
		}

	ExplicitThis?(e)
		{
		return e.pos isnt false and .src[e.pos .. e.pos + #this.Size()] is #this
		}

	Fnlen(fn)
		{
		if fn.type is #Ident
			return fn.name.Size()
		if fn.type is #Mem and fn.mem.type is #Constant and String?(fn.mem.value) and
			fn.mem.value.Identifier?()
			{
			if fn.expr.type is #Ident and fn.expr.name is #this
				return 1 + fn.mem.value.Size()
			if false isnt n = .Fnlen(fn.expr)
				return n + 1 + fn.mem.value.Size()
			}
		return false
		}

	// precedence-based spacing
	// ref: compile/expression.go (highest binding first)
	Nprec: (Or: 4, And: 5, BitOr: 7, BitXor: 8, BitAnd: 9, Cat: 13, Add: 13, Mul: 14)
	Tight?(e, min = 0)
		{
		if e.type is #Binary and e.op is #Mod // Mod is Binary, not Nary: same tier as Mul
			return .Nprec.Mul > min and .Simple(e.lhs) and .Simple(e.rhs)
		if e.type isnt #Nary or e.op not in (#Mul, #Add) or .Nprec[e.op] <= min
			return false
		for (i = 0; false isnt x = e[i]; ++i)
			if not .Simple(x) and
				not (x.type is #Unary and x.op in (#Div, #Sub) and .Simple(x.expr))
				return false
		return true
		}

	OkToSingleLine(node)
		{
		if node.pos isnt false and .src[node.pos .. node.end].Has?('\n')
			return false
		return .OkToSingleLine2(node)
		}

	OkToSingleLine2(node)
		{
		if Type(node) isnt #AstNode
			return true
		if node.type in (#If, #Switch, #TryCatch, #For, #ForIn, #Forever, #While,
			#DoWhile)
			return false
		children = node.children
		for (i = 0; false isnt c = children[i]; ++i)
			if not .OkToSingleLine2(c)
				return false
		return true
		}
	}
