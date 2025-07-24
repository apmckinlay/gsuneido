// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
JsTranslate
	{
	CallClass(ast, outerName = false)
		{
		.NameBegin(outerName, .CompilerNameSuffix[ast.type])
		_className = .Get(#curName)
		_library = .Get(#library)
		(new this)(ast)
		.NameEnd()
		}
	Call(ast)
		{
		Assert(ast.type is: 'Class')
		if ast.base is 'class'
			{
			baseName = false
			base = 'su.root_class'
			}
		else
			{
			baseName = ast.base
			base = 'su.global("' $ baseName $ '")'
			}
		.Println('(function () {') // start of IIFE
		.Indented()
			{
			baseName = baseName is false ? 'false' : '"' $ baseName $ '"'
			.Println('let $super = ' $ baseName $ ';')
			.Println('let $c = Object.create(' $ base $ ');')
			.Println('$c.$setClassInfo("' $
				_library $ '", "' $ _className $ '");')
			.members(ast)
			.Println('Object.freeze($c);')
			.Println('return $c;')
			}
		.Println('}())') // end of IIFE
		}
	members(members)
		{
		for (i = 0; false isnt kv = members[i]; ++i)
			{
			// need to use put() to handle overriding read-only super class
			.Print('$c.put("' $ .convertName(kv.key) $ '", ')

			if Type(kv.value) is 'AstNode' and kv.value.type is 'Function'
				JsTranslateFunction(kv.value, kv.key, inMethod:)
			else
				JsTranslateConstant(kv.value, kv.key)
			.Println(');')
			}
		}
	Privatize(owner, name)
		{
		if _inMethod and owner.type is 'Ident' and owner.name is 'this'
			return .PrivatizeName(name)
		return name
		}

	PrivatizeName(name)
		{
		if name.Prefix?("getter_")
			{
			s = name.RemovePrefix("getter_")
			if s.Size() isnt 0 and not s.Capitalized?()
				return "Getter_" $ _className $ "_" $ s
			}
		if name.Prefix?("get_")
			{
			s = name.RemovePrefix("get_")
			if s.Size() isnt 0 and not s.Capitalized?()
				return "Get_" $ _className $ "_" $ s
			}
		if not name.Capitalized?()
			return _className $ "_" $ name;
		return name;
		}

	convertName(name)
		{
		if name.Capitalized?()
			return name
		if name.Prefix?("getter_")
			return 'Getter_' $ _className $ '_' $ name[7..] /*= 'Getter_'.Size() */
		return _className $ '_' $ name
		}
	}
