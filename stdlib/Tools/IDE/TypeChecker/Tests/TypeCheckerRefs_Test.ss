// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_constructed_via_new()
		{
		refs = TypeCheckerRefs("class { New() { .p = new Point(1, 2) } }")
		Assert(refs.constructed is: #(Point: 1))
		Assert(refs.called is: #())
		}
	Test_bare_call_kept_separate()
		{
		// X(...) is ambiguous (class construction or plain function), so it
		// goes in .called, not .constructed
		refs = TypeCheckerRefs("class { F() { x = new Point(1); y = Bar(2) } }")
		Assert(refs.constructed is: #(Point: 1))
		Assert(refs.called is: #(Bar: 1))
		}
	Test_counts_repeats()
		{
		refs = TypeCheckerRefs(
			"class { F() { new Point(1); new Point(2); new Rect() } }")
		Assert(refs.constructed is: #(Point: 2, Rect: 1))
		}
	Test_nested_construction_in_args()
		{
		refs = TypeCheckerRefs("class { F() { new Outer(new Inner()) } }")
		Assert(refs.constructed is: #(Outer: 1, Inner: 1))
		}
	Test_method_call_is_not_construction()
		{
		// value-method calls are ignored: x.Foo() has a lowercase (local)
		refs = TypeCheckerRefs("class { F(x) { x.Foo(); .Bar() } }")
		Assert(refs.constructed is: #())
		Assert(refs.called is: #())
		}
	Test_lowercase_and_allcaps_calls_skipped()
		{
		// helper() is a local/block, CONST() is a constant neither a class
		refs = TypeCheckerRefs("class { F() { helper(1); CONST(2); Real(3) } }")
		Assert(refs.called is: #(Real: 1))
		Assert(refs.constructed is: #())
		}
	Test_base_class_not_counted()
		{
		// inheritance is gathered by TypeCheckerLineage, not here
		refs = TypeCheckerRefs("BaseControl { New() { super.New() } }")
		Assert(refs.constructed is: #())
		Assert(refs.called is: #())
		}
	Test_parse_error_is_empty()
		{
		refs = TypeCheckerRefs("class { this is not valid (((")
		Assert(refs.constructed is: #())
		Assert(refs.called is: #())
		}
	}
