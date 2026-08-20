// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.src, .cm) { }

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
	// escape; backquotes (raw/regex) and anything already needing escapes
	// stay as written
	Quote(x, s, sym = true) // x = the value, s = the literal as written
		{
		if s[0] is '`'
			return s
		if not .Plain?(x)
			return .requote(x, s)
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

	// escape-carrying or multiline literals: single and double quoted strings
	// process escapes identically, so when the inner text has no quote
	// characters the delimiters swap freely and the escapes move verbatim
	requote(x, s)
		{
		inner = s[1..-1]
		if inner.Has?('"') or inner.Has?("'")
			return s
		return x.Size() is 1 ? "'" $ inner $ "'" : '"' $ inner $ '"'
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

//	NewThis(node)
//		{
//		return #Mem is (fn = node.func).type and #Ident is (expr = fn.expr).type and
//			expr.name is #this and #Constant is (meth = fn.mem).type and
//			meth.value is "*new*"
//		}
	// [...] only replaces a Record/Object call when the literal keeps the same
	// type: brackets are a Record when empty or every member is named, an Object
	// when any member is positional. Converting an all-positional (or mixed)
	// Record, or an all-named/empty Object, would flip the type, so those stay
	// as calls.
	UseBrackets(node)
		{
		if node.func.type isnt #Ident
			return false
		if node.func.name is #Record
			return not .anyPositional?(node)
		if node.func.name is #Object
			return .anyPositional?(node)
		return false
		}

	anyPositional?(node)
		{
		for (i = 0; false isnt arg = node[i]; ++i)
			if arg.name is false
				return true
		return false
		}

	Vertical?(node)
		{
		if false is node.children[1]
			return false
		if false is prev = .openPos(node)
			return false
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

	// call arguments written one per line in the source stay one per line
	VerticalArgs?(node, last)
		{
		if last < 1
			return false
		for (i = 1; i <= last; ++i)
			if not .ArgBrokeInSrc?(node, i)
				return false
		return true
		}

	// did the source break between two consecutive call arguments?
	// Argument nodes and composite exprs carry no pos/end - recover spans by
	// descending to the outermost positioned descendants
	ArgBrokeInSrc?(node, i)
		{
		if false is e = .spanEnd(node[i-1].expr)
			return false
		if false is s = .spanStart(node[i].expr)
			return false
		// a block argument always gets a line of its own, so the break before it
		// is the formatter's, not the author's - counting it would make the next
		// pass see a vertical argument list that the pass before it created
		if node[i].expr.type is #Block
			return false
		return .src[e..s].Has?('\n')
		}

	spanStart(e)
		{
		if Type(e) isnt #AstNode
			return false
		if e.pos not in (0, false)
			return e.pos
		if e.type is #Mem and e.dotpos not in (0, false)
			return .dotStart(e)
		for (i = 0; false isnt c = e.children[i]; ++i)
			if false isnt r = .spanStart(c)
				return r
		return false
		}

	// dotpos is the member name position: an implicit-this member's
	// synthesized ident sits there too, so the dot before it is lost
	dotStart(e)
		{
		s = .spanStart(e.expr)
		if s isnt false and s < e.dotpos
			return s
		i = e.dotpos - 1
		while i > 0 and .src[i] in (' ', '\t', '\n', '\r')
			--i
		return i
		}

	spanEnd(e)
		{
		if Type(e) isnt #AstNode
			return false
		if e.end not in (0, false)
			return e.end
		// implicit-this member .name has no end; its synthesized ident sits at
		// dotpos, so recover the name's end from dotpos (mirror of dotStart)
		if e.type is #Mem and e.dotpos not in (0, false) and Type(e.mem) is #AstNode and
			e.mem.type is #Constant and String?(e.mem.value)
			return e.dotpos + e.mem.value.Size()
		n = 0
		while false isnt e.children[n]
			++n
		for (i = n - 1; i >= 0; --i)
			if false isnt r = .spanEnd(e.children[i])
				return e.type in (#RangeTo, #RangeLen, #Mem) ? .closeBracket(r) : r
		return false
		}

	// range and subscript nodes have no pos/end of their own, so the recovered
	// end is the last child's, which stops before the closing ]; include it,
	// else a verbatim slice emits an unbalanced bracket
	closeBracket(r)
		{
		i = r
		while i < .src.Size() and .src[i] in (' ', '\t', '\r', '\n', '.', ':', '[')
			++i
		return i < .src.Size() and .src[i] is ']' ? i + 1 : r
		}

	// did the source break between the open delimiter and the first argument?
	// scan back to the delimiter, over any named-argument key
	LeadBreak?(arg)
		{
		if false is p = .spanStart(arg.expr)
			return false
		i = p - 1
		while i > 0 and .src[i] not in ('\n', '(', '[')
			--i
		return .src[i] is '\n'
		}

	// a nested table node has no pos of its own - recover the open delimiter
	// by scanning back over the whitespace preceding the first member
	openPos(node)
		{
		if node.pos isnt 0 and node.pos isnt false
			return node.pos
		if false is m = node[0]
			return false
		if m.pos in (0, false)
			return false
		i = m.pos - 1
		while i > 0 and .src[i] in (' ', '\t', '\n', '\r')
			--i
		return i
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

	// per-member pad widths that align contiguous runs of scalar 'key: value'
	// class members into a column; methods and complex values break a run
	ClassAlign(node)
		{
		pads = Object()
		i = 0
		while false isnt node[i]
			{
			if not .classAlignable?(node[i])
				{
				pads[i] = false
				++i
				continue
				}
			// gather the run of alignable members and its widest key
			run = Object()
			w = 0
			while false isnt node[i] and .classAlignable?(node[i])
				{
				w = Max(w, node[i].key.Size())
				run.Add(i)
				++i
				}
			for j in run
				pads[j] = w
			}
		return pads
		}

	classAlignable?(m)
		{
		return .BareWord?(m.key) and m.value isnt true and Type(m.value) isnt #AstNode
		}

	// a $ chain containing a multiline string keeps its authored layout: the
	// string resets to column 0, so structural indent is meaningless there
	MultilineCatSpan(node)
		{
		if node.op isnt #Cat or not .anyMlString?(node)
			return false
		if false is s = .spanStart(node[0])
			return false
		if false is e = .spanEnd(node[node.size - 1])
			return false
		return [pos: s, end: e]
		}

	anyMlString?(node)
		{
		for (i = 0; false isnt c = node[i]; ++i)
			if .mlString?(c)
				return true
		return false
		}

	mlString?(c)
		{
		return c.type is #Constant and String?(c.value) and c.pos not in (0, false) and
			c.end not in (0, false) and .src[c.pos .. c.end].Has?('\n')
		}

	// a chain of two or more dotted method calls breaks after each dot when it
	// cannot fit on one line
	ChainBreak?(node)
		{
		return .DottedCall?(node) and .DottedCall?(node.func.expr)
		}

	DottedCall?(e)
		{
		if Type(e) isnt #AstNode or e.type isnt #Call or e.func.type isnt #Mem or
			e.func.mem.type isnt #Constant or not String?(e.func.mem.value) or
			not e.func.mem.value.Identifier?()
			return false
		// a trailing block argument hangs its body outside the chain group
		return e.size is 0 or e[e.size - 1].name isnt #block
		}

	// a vertical table whose rows are uniform flat literals (same keys in the
	// same order, all scalar values) aligns as a grid: column widths, or false
	TableColumns(node, width, tabWidth)
		{
		if false is node[2] // a grid needs at least 3 rows
			return false
		keys = false
		cols = Object()
		for (i = 0; false isnt m = node[i]; ++i)
			{
			if false is row = .rowCells(m)
				return false
			if keys is false
				keys = row.keys
			else if keys isnt row.keys
				return false
			for (j = 0; j < row.cells.Size(); ++j)
				cols[j] = Max(cols.GetDefault(j, 0), row.cells[j].Size())
			}
		return .tableFits?(node, cols, width, tabWidth) ? cols : false
		}

	rowCells(m)
		{
		if not .tableRow?(m)
			return false
		keys = Object()
		cells = Object()
		for (i = 0; false isnt c = m.value[i]; ++i)
			{
			if not .scalarCell?(c) or (cell = .CellText(c)).Has?('\n')
				return false
			keys.Add(c.key)
			cells.Add(cell)
			}
		return [:keys, :cells]
		}

	tableRow?(m)
		{
		return m.named isnt true and m.pos not in (0, false) and
			m.end not in (0, false) and Type(m.value) is #AstNode and
			m.value.type is #Object and false isnt m.value.children[0] and
			not .cm.CommentIn?(m.pos, m.end) and not .Vertical?(m.value)
		}

	scalarCell?(c)
		{
		return c.named is true and c.value isnt true and Type(c.value) isnt #AstNode and
			.BareWord?(c.key)
		}

	CellText(c)
		{
		return c.key $ ": " $ .cellValue(c)
		}

	cellValue(c)
		{
		x = c.value
		if .BareWord?(x)
			return x
		s = .MemberSrc(c, keyed:)
		if String?(x)
			return .Quote(x, s is false ? Display(x) : s, sym: false)
		return Number?(x) ? .Num(x, s) : Display(x)
		}

	// 1234 -> 1234
	// 12345 -> 12_345
	// 0xff -> 0x ff
	// 0xABCDABCD -> 0xABCD_ABCD
	Num(x, s)
		{
		hexLower = 2
		hexUpper = 4
		decimalGroupSize = 3
		if s is false
			s = Display(x)
		if s =~ "^0[xX][0-9a-fA-F_]+$"
			return .digitGroups(s, hexLower, hexUpper)
		if s =~ "^-?[0-9_]+$"
			return .digitGroups(s, s[0] is '-' ? 1 : 0, decimalGroupSize)
		return s
		}

	digitGroups(s, pre, n)
		{
		digitSplitThresh = 4
		digits = s[pre..].Tr('_')
		if digits.Size() <= digitSplitThresh
			return s[..pre] $ digits
		grouped = ""
		for (i = digits.Size(); i > n; i -= n)
			grouped = '_' $ digits[i-n .. i] $ grouped
		return s[..pre] $ digits[..i] $ grouped
		}

	tableFits?(node, cols, width, tabWidth)
		{
		w = 3 // '(' ')' and the row comma
		for c in cols
			w += c + 2
		i = node[0].pos - 1
		while i >= 0 and .src[i] isnt '\n'
			{
			w += .src[i] is '\t' ? tabWidth : 1
			--i
			}
		return w <= width
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

	// a debug statement written at the left margin stays there; one written
	// with the code is laid out like any other call
	Debug?(expr)
		{
		if expr.type isnt #Call or expr.func.type isnt #Ident
			return false
		if expr.func.name not in (#Print, #TracePrint, #ServerPrint, #StackTrace,
			#Inspect, #TraceCallStack)
			return false
		p = .spanStart(expr)
		return p isnt false and (p is 0 or .src[p-1] is '\n')
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

	ExplicitThis?(node)
		{
		// implicit .foo synthesizes its this ident at the member position
		e = node.expr
		return e.pos isnt false and e.pos isnt node.dotpos
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

	// new X(...) must not convert to X["*new*"](...)
	NewCall(node)
		{
		fn = node.func
		return fn.type is #Mem and fn.mem.type is #Constant and fn.mem.value is "*new*"
			? fn.expr
			: false
		}

	RecordBrackets?(node)
		{
		for (i = 0; false isnt m = node[i]; ++i)
			if m.named isnt true
				return false
		return true
		}
	}
