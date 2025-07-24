// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.test('class { }',
			'(function () {
				let $super = false;
				let $c = Object.create(su.root_class);
				$c.$setClassInfo("", "eval");
				Object.freeze($c);
				return $c;
			}())')
		.test('Base { }',
			'(function () {
				let $super = "Base";
				let $c = Object.create(su.global("Base"));
				$c.$setClassInfo("", "eval");
				Object.freeze($c);
				return $c;
			}())')
		.test('class : Base { }',
			'(function () {
				let $super = "Base";
				let $c = Object.create(su.global("Base"));
				$c.$setClassInfo("", "eval");
				Object.freeze($c);
				return $c;
			}())')
		.test('class : Base { M: 123; F(.a, ._M? = 4) { }  G(@args) { } }',
			`(function () {
				let $super = "Base";
				let $c = Object.create(su.global("Base"));
				$c.$setClassInfo("", "eval");
				$c.put("M", 123);
				$c.put("F", (function () {
					let $callNamed = function (named, a, m_Q) {
						su.maxargs(3, arguments.length);
						({ a = a } = named);
						if (named["m?"] !== undefined)
							m_Q = named["m?"];
						return $f.call(this, a, m_Q);
					};
					let $f = function (a = su.mandatory(), m_Q = su.dynparam("_m?", 4)) {
						su.put(this, "eval_a", a);
						su.put(this, "M?", m_Q);
						su.maxargs(2, arguments.length);
					};
					$f.$callableType = "METHOD";
					$f.$callableName = "eval#F";
					$f.$call = $f;
					$f.$callNamed = $callNamed;
					$f.$callAt = function (args) {
						return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
					};
					$f.$params = 'a, m?=4';
					return $f;
					})());
				$c.put("G", (function () {
					let $f = function (args = su.mandatory()) {
						su.maxargs(1, arguments.length);
					};
					$f.$callableType = "METHOD";
					$f.$callableName = "eval#G";
					$f.$callAt = $f;
					$f.$call = function (...args) {
						return $f.call(this, su.mkObject2(args));
					};
					$f.$callNamed = function (named, ...args) {
						return $f.call(this, su.mkObject2(args, su.obToMap(named)));
					};
					$f.$params = '@args';
					return $f;
					})());
				Object.freeze($c);
				return $c;
			}())`)
		.test('class B { New() { super() } }',
			`(function () {
				let $super = "B";
				let $c = Object.create(su.global("B"));
				$c.$setClassInfo("", "eval");
				$c.put("New", (function () {
					let $callNamed = function (named) {
						su.maxargs(1, arguments.length);
						return $f.call(this);
					};
					let $f = function () {
						su.maxargs(0, arguments.length);
						return su.invokeBySuper($super, "New", this);
					};
					$f.$callableType = "METHOD";
					$f.$callableName = "eval#New";
					$f.$call = $f;
					$f.$callNamed = $callNamed;
					$f.$callAt = function (args) {
						return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
					};
					$f.$params = '';
					return $f;
					})());
				Object.freeze($c);
				return $c;
			}())`)
		.testCatch('class B { Fn() { super() } }', 'super call only allowed in New')
		.testCatch('class B { New() { 1; super() } }', 'call to super must come first')
		}
	test(src, expected)
		{
		Assert(JsTranslate(src) like: expected)
		}
	testCatch(src, expected)
		{
		Assert({ JsTranslate(src) } throws: expected)
		}
	}
