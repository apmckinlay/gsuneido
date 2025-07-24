// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// Evaluate constant expressions at compile time.
// Returns a value if folded, otherwise the original ast.
// Only returns values of boolean, number, or string
class
	{
	CallClass(ast)
		{
		switch ast.type
			{
		case 'Constant':
			return  Type(ast.value) isnt 'AstNode' ? ast.value : ast
		case 'Unary':
			return .unary(ast)
		case 'Nary':
			return .nary(ast)
		case 'Binary':
			return .binary(ast)
		default:
			return ast
			}
		}
	unary(ast)
		{
		arg = ast.expr
		val = AstFoldExpr(arg) // recurse
		if Same?(val, arg)
			return ast
		switch ast.op
			{
		case 'Add':		return +val
		case 'Sub':		return -val
		case 'Not':		return not val
		case 'BitNot':	return ~val
		case 'Div':		return 1/val
		case 'LParen':	return val
			}
		}
	nary(ast)
		{
		if ast.op in ('And', 'Or')
			return .andor(ast)

		values = Object()
		for i in .. ast.size
			{
			if Same?(ast[i], val = AstFoldExpr(ast[i]))
				return ast
			values.Add(val)
			}
		res = values[0]
		fn = .fns[ast.op]
		for (i = 1; i < values.Size(); i++)
			res = fn(res, values[i])
		return res
		}
	binary(ast)
		{
		leftast = ast.lhs
		rightast = ast.rhs
		left = AstFoldExpr(leftast) // recurse
		right = AstFoldExpr(rightast) // recurse
		if Same?(left, leftast) or Same?(right, rightast)
			return ast
		return (.fns[ast.op])(left, right)
		}
	fns: #(
		// nary
		'Cat':			function (left, right) { return left $ right },
		'Add':			function (left, right) { return left + right },
		'Sub':			function (left, right) { return left - right },
		'Mul':			function (left, right) { return left * right },
		'BitOr':		function (left, right) { return left | right },
		'BitAnd':		function (left, right) { return left & right },
		'BitXor':		function (left, right) { return left ^ right },
		// binary
		'Mod':			function (left, right) { return left % right },
		'LShift':		function (left, right) { return left << right },
		'RShift':		function (left, right) { return left >> right },
		'Is':			function (left, right) { return left is right },
		'Isnt':			function (left, right) { return left isnt right },
		'Lt':			function (left, right) { return left < right },
		'Lte':			function (left, right) { return left <= right },
		'Gt':			function (left, right) { return left > right },
		'Gte':			function (left, right) { return left >= right },
		'Match':		function (left, right) { return left =~ right },
		'MatchNot':		function (left, right) { return left !~ right })
	andor(ast)
		{
		for i in .. ast.size
			{
			val = AstFoldExpr(ast[i])
			if ast.op is 'And' and val is false
				return false
			if ast.op is 'Or' and val is true
				return true
			}
		return ast
		}
	}
