// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.scan = ScanCode('hello')
		Assert(.scan.Next() is: 'hello')
		}
	Test_context()
		{
		.assert_context('x', '')

		.assert_context('function () { x }', 'code')
		.assert_context('function (x) { }', 'code')

		.assert_context('class { x }', 'class')
		.assert_context('Controller { x }', 'class')
		.assert_context('#{ x }', 'constant')
		.assert_context('#{ { x } }', 'constant')
		.assert_context('#( x )', 'constant')
		.assert_context('#( ( ) x )', 'constant')
		.assert_context('class { m: ( X () ) }', 'constant')
		.assert_context('class : X { }', 'code')
		.assert_context('class X { }', 'code')
		.assert_context('X { }', 'code')
		.assert_context('_X { }', 'code')
		.assert_context('#{ F (x) }', 'constant')

		.assert_context('function () { } x', '')
		.assert_context('class { } x', '')
		.assert_context('#( ) x', '')
		.assert_context('#( ( ) ) x', '')

		.assert_context('function () { #(x) }', 'constant')
		.assert_context('function () { #() x }', 'code')
		.assert_context('#{ function () { x } }', 'code')
		.assert_context('#{ function () { } x }', 'constant')

		.assert_context('class { X(){ } }', 'class')
		.assert_context('class { F(x){ } }', 'code')
		.assert_context('class { F(q, y={}){ y e() } z:(x) }', 'constant' )
		.assert_context('class { F(q){ y({ |z| z }) x } }', 'code')
		.assert_context('class { F(q){ y({ |z| z }) y = #(x, 2) } }', 'constant')
		.assert_context('class { F(){ x } }', 'code')
		.assert_context('class { f(){ x } }', 'code')

		.assert_context('dll ( x )', 'dll')

		.assert_context('class { F() { if a is "PDF" { x() } } }', 'code')
		.assert_context('class { F() { if a is #PDF { x() } } }', 'code')
		}
	assert_context(code, context)
		{ // checks context of 'x'
		.scan = ScanCode(code)
		while 'x' isnt .scan.Next().Lower().Tr('_')
			{ }
		Assert(.scan.Context() is: context)
		}
	}