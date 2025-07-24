// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.test('function (@args) { return args }',
			`(function () {
				let $f = function (args = su.mandatory()) {
					su.maxargs(1, arguments.length);
					return args;
				};
				$f.$callableType = "FUNCTION";
				$f.$callableName = "eval";
				$f.$callAt = $f;
				$f.$call = function (...args) {
					return $f.call(this, su.mkObject2(args));
				};
				$f.$callNamed = function (named, ...args) {
					return $f.call(this, su.mkObject2(args, su.obToMap(named)));
				};
				$f.$params = '@args';
				return $f;
				})()`)
		.test('function (a?, b = "\n") { return #(123) }',
			`(function () {
				let $const = [
					su.mkObject(123),
				];
				let $callNamed = function (named, a_Q, b) {
					su.maxargs(3, arguments.length);
					({ b = b } = named);
					if (named["a?"] !== undefined)
						a_Q = named["a?"];
					return $f.call(this, a_Q, b);
				};
				let $f = function (a_Q = su.mandatory(), b = "\n") {
					su.maxargs(2, arguments.length);
					return $const[0];
				};
				$f.$callableType = "FUNCTION";
				$f.$callableName = "eval";
				$f.$call = $f;
				$f.$callNamed = $callNamed;
				$f.$callAt = function (args) {
					return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
				};
				$f.$params = 'a?, b="\n"';
				return $f;
				})()`)
		.test(`function (_a? = 1) { _b = 1; return a? + _b + 1.1 }`,
			`(function () {
				let $const = [
					su.mknum("1.1"),
				];
				let $callNamed = function (named, a_Q) {
					su.maxargs(2, arguments.length);
					if (named["a?"] !== undefined)
						a_Q = named["a?"];
					return $f.call(this, a_Q);
				};
				let $f = function (a_Q = su.dynparam("_a?", 1)) {
					try { su.dynpush();
					su.maxargs(1, arguments.length);
					su.dynset("_b", 1);
					return su.add(su.add(a_Q, su.dynget("_b")), $const[0]);
					} finally { su.dynpop(); }
				};
				$f.$callableType = "FUNCTION";
				$f.$callableName = "eval";
				$f.$call = $f;
				$f.$callNamed = $callNamed;
				$f.$callAt = function (args) {
					return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
				};
				$f.$params = 'a?=1';
				return $f;
				})()`)
		.test(`function () { b = { |a| c = a++ } }`,
			`(function () {
				let $callNamed = function (named) {
					su.maxargs(1, arguments.length);
					return $f.call(this);
				};
				let $f = function () {
					su.maxargs(0, arguments.length);
					var a, b, c;
					return b = (function () {
						let $callNamed = function (named, a) {
							su.maxargs(2, arguments.length);
							({ a = a } = named);
							return $f.call(this, a);
						};
						let $f = function (a = su.mandatory()) {
							su.maxargs(1, arguments.length);
							var $tmp0;
							return c = ((a = su.inc($tmp0 = a)), $tmp0);
						};
						$f.$callableType = "BLOCK";
						$f.$callableName = "eval$b";
						$f.$call = $f;
						$f.$callNamed = $callNamed;
						$f.$callAt = function (args) {
							return $callNamed.call(this, su.mapToOb(args.map), ` $
								`...args.vec);
						};
						$f.$params = 'a';
						return $f;
						})();
				};
				$f.$callableType = "FUNCTION";
				$f.$callableName = "eval";
				$f.$call = $f;
				$f.$callNamed = $callNamed;
				$f.$callAt = function (args) {
					return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
				};
				$f.$params = '';
				return $f;
				})()`)
		.test(`function () { b = { .a } }`,
			`(function () {
			let $callNamed = function (named) {
				su.maxargs(1, arguments.length);
				return $f.call(this);
			};
			let $f = function () {
				su.maxargs(0, arguments.length);
				var b;
				return b = (function () {
					let $callNamed = function (named) {
						su.maxargs(1, arguments.length);
						return $f.call(this);
					};
					let $f = function () {
						var blockThis = su.getBlockThis(this);
						su.maxargs(0, arguments.length);
						return su.get(blockThis, "a");
					};
					$f.$callableType = "BLOCK";
					$f.$callableName = "eval$b";
					$f.$blockThis = this;
					$f.$call = $f;
					$f.$callNamed = $callNamed;
					$f.$callAt = function (args) {
						return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
					};
					$f.$params = '';
					return $f;
					}).call(this);
			};
			$f.$callableType = "FUNCTION";
			$f.$callableName = "eval";
			$f.$call = $f;
			$f.$callNamed = $callNamed;
			$f.$callAt = function (args) {
				return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
			};
			$f.$params = '';
			return $f;
			})()`)
		.test(`function (a = #(1234), b = 1.1, c = #20240101)
				{
				b = 1.1
				c = #20240101
				d = 1.1 + 1.1
				return Same?(a, #(1234))
				}`,
			`(function () {
			let $const = [
				su.mkObject(1234),
				su.mknum("1.1"),
				su.mkdate('20240101'),
				su.mknum("2.2"),
			];
			let $callNamed = function (named, a, b, c) {
				su.maxargs(4, arguments.length);
				({ a = a, b = b, c = c } = named);
				return $f.call(this, a, b, c);
			};
			let $f = function (a = $const[0], b = $const[1], c = $const[2]) {
				su.maxargs(3, arguments.length);
				var d;
				b = $const[1];
				c = $const[2];
				d = $const[3];
				return su.call(su.global("Same?"), a, $const[0]);
			};
			$f.$callableType = "FUNCTION";
			$f.$callableName = "eval";
			$f.$call = $f;
			$f.$callNamed = $callNamed;
			$f.$callAt = function (args) {
				return $callNamed.call(this, su.mapToOb(args.map), ...args.vec);
			};
			$f.$params = 'a=, b=1.1, c=#20240101';
			return $f;
			})()`
			)
		}
	test(src, expected)
		{
		Assert(JsTranslate(src) like: expected)
		}
	}
