// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// Translate Suneido code to JavaScript
class
	{
	CallClass(src, globalName = 'eval', library = '')
		{
		ast = Suneido.Parse(src)
		_jst = Object(output: '', indent: 0, nextnum: 0, curName: '', :globalName,
			:library)
		JsTranslateConstant(ast, isConst:)
		return _jst.output
		}
	Get(member)
		{
		if _jst.Member?(member)
			return _jst[member]
		return
		}

	// output stuff
	Indented(block)
		{
		++_jst.indent
		block()
		--_jst.indent
		}
	Print(s)
		{
		if _jst.output[-1] is '\n'
			_jst.output $= ' '.Repeat(4/*=spaces*/ *_jst.indent)
		_jst.output $= s
		}
	Println(s = '')
		{
		.Print(s $ '\n')
		}
	PrintPos()
		{
		return [pos: _jst.output.Size(), indent: _jst.indent]
		}
	At(pos, block)
		{
		tail = _jst.output[pos.pos ..]
		indent = _jst.indent
		_jst.indent = pos.indent
		_jst.output = _jst.output[.. pos.pos]
		block()
		_jst.output $= tail
		_jst.indent = indent
		}

	Const(ast)
		{
		if false isnt i = _const.FindIf({ .compare(it, ast) })
			{
			.Print('$const[' $ i $ ']')
			return
			}
		// for complex constants (not boolean, integer, or string)
		// we save the ast to generate later
		// and just print a reference to it
		i = _const.Size()
		.Print('$const[' $ i $ ']')
		_const.Add(ast)
		}

	compare(const1, const2)
		{
		if Type(const1) isnt 'AstNode' or Type(const2) isnt 'AstNode'
			return const1 is const2

		return AstSearch.Compare2(const1, const2, noSpecial?:)
		}

	MAX_JS_SAFE_INTEGER: 9007199254740991
	Value(val, isConst = false)
		{
		switch
			{
		case Boolean?(val):
			.Print(Display(val))
		case Number?(val):
			isConst ? .constNumber(val) : .number(val)
		case String?(val):
			.Print('"' $ val.Escape() $ '"')
		case Date?(val):
			isConst ? .Print("su.mkdate('" $ Display(val)[1..] $ "')") : .Const(val)
			}
		}

	number(val)
		{
		if val is val.Int() and
			-.MAX_JS_SAFE_INTEGER < val and val <= .MAX_JS_SAFE_INTEGER
			.Print(Display(val))
		else
			.Const(val)
		}

	constNumber(val)
		{
		if val.Int?() and not IsInf?(val)
			.Print(val)
		else
			.Print('su.mknum("' $ val $ '")')
		}

	NextNum() // used by JsTranslateFunction for function id for block return
		{
		return ++_jst.nextnum
		}

	CompilerNameSuffix: #{
		FUNCTION: '$f'
		BLOCK: '$b'
		METHOD: '$m'
		CLASS: '$c'
		// for gSuneido parse
		Function: '$f'
		Block: '$b'
		Method: '$m'
		Class: '$c'
		}
	MethodSeparaor: '#'
	NameBegin(memberName, def)
		{
		if _jst.curName is ''
			_jst.curName = _jst.globalName;
		else if memberName is false
			_jst.curName $= def;
		else
			_jst.curName $= .MethodSeparaor $ memberName;
		}
	NameEnd()
		{
		i = Max(_jst.curName.FindLast('$'), _jst.curName.FindLast(.MethodSeparaor))
		_jst.curName = i is false ? '' : _jst.curName[..i]
		}
	ConvertQMark(name)
		{
		if name.Suffix?('?')
			return name[..-1] $ "_Q"
		return name
		}
	}
