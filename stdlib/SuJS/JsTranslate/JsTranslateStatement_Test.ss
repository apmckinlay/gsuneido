// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: (
		('return true',				'return true;')

		('throw "error"',			'throw exception("error");')

		('forever ;',				'while (true) {
										}')

		('forever stmt',			'while (true) {
										stmt;
									}')

		('forever
			{ stmt }',				'while (true) {
										stmt;
									}')

		('forever
			{ if (c) break }',		'while (true) {
										if (toBool(c)) {
											break;
										}
									}')

		('forever
			{ if (c) continue }',	'while (true) {
										if (toBool(c)) {
											continue;
										}
									}')

		('while (a < b) stmt',		'while (lt(a, b)) {
										stmt;
									}')

		('for (x in ob) f(x)',
			'for ($tmp0 = iter(ob); $tmp0 !== (x = next($tmp0)); ) {
				call(f, x);
			}')

		('if (a) ;',				'if (toBool(a)) {
									}')

		('if (a) stmt',				'if (toBool(a)) {
										stmt;
									}')

		('if (a) b
		  else c',					'if (toBool(a)) {
										b;
									} else {
										c;
									}')

		('if (a) ;
		  else c',					'if (toBool(a)) {
									} else {
										c;
									}')

		('if (a) b
		  else ;',					'if (toBool(a)) {
										b;
									}')

		('if (a) b
		  else if (c) d
		  else e',					'if (toBool(a)) {
										b;
									} else if (toBool(c)) {
										d;
									} else {
										e;
									}')

		('for (a; b < c; d)
			e',						'for (a; lt(b, c); d) {
										e;
									}')

		('for (a, a2; b < c; d, d2)
			e',						'for (a, a2; lt(b, c); d, d2) {
										e;
									}')
		('for (; ; )
			e',						'for (; ; ) {
										e;
									}')
		('for (a; ; )
			e',						'for (a; ; ) {
										e;
									}')
		('for (; b; )
			e',						'for (; toBool(b); ) {
										e;
									}')
		('for (; ; c)
			e',						'for (; ; c) {
										e;
									}')
		('for i in a..b c',			'$tmp0 = (a);
$tmp1 = (b);
for (i = $tmp0; lt(i, $tmp1); i = inc(i)) {
c;
}')
		('for ..b c',				'$tmp1 = (0);
$tmp2 = (b);
for ($tmp0 = $tmp1; lt($tmp0, $tmp2); $tmp0 = inc($tmp0)) {
c;
}')
		('for m,v in ob e',
			'for ($tmp1 = iter(invoke(ob, "Assocs")); $tmp1 !== ($tmp0 = next($tmp1)); ) {
m = $tmp0.get(0);
v = $tmp0.get(1);
{
e;
}
}')


		('do stmt while (a < b)',	'do {
										stmt;
									} while (lt(a, b));')
		('switch { case 1,2: a }',	'switch (true) {
									case 1:
									case 2:
										a;
										break;
									default:
										throw exception("unhandled switch case");
									}')
		('switch (a) { case 1: return b; default: c }',
									'switch (a) {
									case 1:
										return b;
									default:
										c;
									}')
		('switch (a) { case 1: default: }',
									'switch (a) {
									case 1:
										break;
									default:
										break;
									}')
		('try x catch (e) y',
									`try {
										x;
									} catch (_e) {
										var e = catchMatch(_e, "");
										y;
									}`)
		('try x catch (e, "foo") y',
									`try {
										x;
									} catch (_e) {
										var e = catchMatch(_e, "foo");
										y;
									}`)
		('try x',
									`try {
										x;
									} catch (_e) {
									}`)
		('return 1, a',
									`return [1, a];`,
									#20250227)
		('a, b, c = fn()',
									`[a, b, c] = (call(fn));`,
									#20250227)
		)
	Test_one()
		{
		for c in .cases
			{
			if c.Member?(2) and BuiltDate() < c[2]
				continue
			.test(c[0], c[1])
			}
		}
	test(src, expected)
		{
		dst = JsTranslate('function () {\n' $ src $ '\n}')
		dst = dst.Replace('su\.')
		lines = dst.Lines()
		i = 1 + lines.Find("    let $f = function () {") + 1 // skip maxargs
		if lines[i].Has?(' var ')
			++i
		j = lines[i..].Find('    };')
		dst = lines[i :: j].Join('\n')
		Assert(dst like: expected)
		}
	Test_block()
		{
		s = JsTranslate('function () { f = function () { return } }')
		Assert(s hasnt: 'throw su.blockreturn')
		Assert(s hasnt: 'su.rethrowBlockReturn(_e)')
		s = JsTranslate('function () { b = { return } }')
		Assert(s has: 'throw su.blockreturn')
		s = JsTranslate('function () { b = { return }; try b() }')
		Assert(s has: 'su.rethrowBlockReturn(_e)')
		Assert(s has: 'throw su.blockreturn')
		for tok in #(':break', ':continue')
			{
			s = JsTranslate('function () { b = { ' $ tok[1..] $ ' } }')
			Assert(s has: 'throw su.exception("block' $ tok $ '")')

			s = JsTranslate('function () { b = { forever { ' $ tok[1..] $ ' } } }')
			Assert(s hasnt: 'throw su.exception("block' $ tok $ '")')
			Assert(s has: tok[1..] $ ';')

			s = JsTranslate('function () { forever { b = { ' $ tok[1..] $ ' } } }')
			Assert(s has: 'throw su.exception("block' $ tok $ '")')
			}
		}
	}
