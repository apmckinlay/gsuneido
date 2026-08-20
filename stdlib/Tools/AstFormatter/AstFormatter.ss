// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// ref: Wadler, P., & Kilmer, J. (2002). A prettier printer.
AstFmtStmt
	{
	// argWrap
	//		fill = pack arguments, breaking only when needed
	//		onePerLine = every argument on each line
	// bodyStyle
	//		preserve = keep the original source code's shape
	//		fit = collapse any body to one line
	//		expand = every body breaks
	CallClass(src, width = 90, tabWidth = 4, argWrap = #fill, bodyStyle = #preserve)
		{
		Assert(Compilable?(src))
		Assert(Number?(width) and Number?(tabWidth))
		Assert(#(fill, onePerLine).Has?(argWrap))
		Assert(#(preserve, fit, expand).Has?(bodyStyle))

		// body can change layout upstream; iterate to the fixpoint
		out = (new this(width, tabWidth, argWrap, bodyStyle)).Process(src)
		for ..2
			{
			next = (new this(width, tabWidth, argWrap, bodyStyle)).Process(out)
			if next is out
				break
			out = next
			}

		// we must ensure the output code always compiles; else hard error
		Assert(Compilable?(out))
		Assert(AstFmtEquals?(src, out))
		return out
		}

	New(.Width = 90, .TabWidth = 4, .ArgWrap = #fill, .BodyStyle = #preserve)
		{
		Assert(Number?(.Width) and Number?(.TabWidth))
		Assert(#(fill, onePerLine).Has?(.ArgWrap))
		Assert(#(preserve, fit, expand).Has?(.BodyStyle))
		}

	maxpos: 9_999_999
	Process(src)
		{
		.Src = src
		ast = Suneido.Parse(src)
		if Type(ast) isnt #AstNode // simple constant record e.g. a number
			return .Rtrim(src) $ '\n'
		.Cm = AstFmtComments(src)
		.Style = AstFmtStyle(src, .Cm)
		curr = [i: 0, done: 0, blank: false]
		ctx = [constant: false, noSym: false, cond: false, singleLine: false,
			blockArg: false, tight: false, method: false, lastStmt: false, bracket: false,
			cols: false, chain: false, guard: false, parentEnd: .maxpos]
		doc = .Cat(.Fmt(ast, ctx, curr), .Cm.Leading(curr, .maxpos))
		return .Rtrim(AstFmtRender(.Width, .TabWidth).Render(doc)) $ '\n'
		}

	Default(type, node, ctx/*unused*/, curr)
		{
		if Type(node) is #AstNode and node.pos isnt false and node.end not in (0, false)
			return .verbatimNode(node, curr)
		throw "AstFormatter: cannot format " $ Display(type)
		}

	verbatimNode(node, curr)
		{
		.Cm.SkipTo(curr, node.end)
		return .Tok(.Src[node.pos .. node.end])
		}
	}

