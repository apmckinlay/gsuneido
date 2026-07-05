// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// finds the global classes a record constructs, so the type-check bundle can
// pull in their sources the checker needs them to give `new X(...)` the
// nominal type Instance(X) instead of Unknown.
//
// modeled on Qc_globalRefs, but where that record counts every capitalized ident the
// same, this separates unambiguous construction (`new X(...)`) from a bare
// call (`X(...)`, which might be a class or a plain function). keeping them
// apart lets the gatherer follow only construction edges; sparse and (assuming) acyclic
// instead of the whole coupling graph
class
	{
	CallClass(code)
		{
		try
			ast = Suneido.Parse(code)
		catch
			return Object(constructed: Object(), called: Object())

		constructed = Object()
		called = Object()
		.traverse(ast, constructed, called)

		return Object(:constructed, :called)
		}

	traverse(ast, constructed, called)
		{
		if Type(ast) isnt 'AstNode'
			return
		if ast.type is 'Call'
			.classifyCall(ast, constructed, called)
		for (i = 0; false isnt c = ast.children[i]; ++i)
			.traverse(c, constructed, called) // recursive
		}

	// `new X(...)` parses to a Call on Mem{expr: Ident X, mem: "*new*"};
	// a bare `X(...)` is a Call straight on Ident X; a static method call
	// `X.Foo(...)` is a Call on Mem{expr: Ident X, mem: "Foo"}. all three
	// edges pull X into the bundle so the checker can resolve construction,
	// CallClass, and static return types like Json.Encode. value-method
	// calls (v.Foo(), super.New(), new on a non-ident) are left out
	classifyCall(call, constructed, called)
		{
		fn = call.func
		if Type(fn) isnt 'AstNode'
			return
		if fn.type is 'Mem'
			{
			// both `new X(...)` and `X.Foo(...)` have a capitalized Ident
			// receiver; anything else on a Mem is a value-method call.
			if fn.expr.type isnt 'Ident' or not fn.expr.name.Capitalized?()
				return
			if .newCall?(fn)
				.bump(constructed, fn.expr.name)
			else if not fn.expr.name.Upper?()
				.bump(called, fn.expr.name)   // static call: Json.Encode(...)
			}
		else if fn.type is 'Ident' and fn.name.Capitalized?() and
			not fn.name.Upper?()
			.bump(called, fn.name)
		}

	newCall?(mem)
		{
		m = mem.mem
		return Type(m) is 'AstNode' and m.type is 'Constant' and m.value is '*new*'
		}

	bump(counts, name)
		{
		counts[name] = counts.GetDefault(name, 0) + 1
		}
	}
