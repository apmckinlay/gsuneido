// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Type(true) is: 'Boolean')
		Assert(Type(false) is: 'Boolean')
		Assert(Type(123) is: 'Number')
		Assert(Type(-123) is: 'Number')
		Assert(Type(12.3) is: 'Number')
		Assert(Type('') is: 'String')
		Assert(Type("hello") is: 'String')
		Assert(Type(#(a:).Members()[0]) is: 'String')
		Assert(Type(#20030830) is: 'Date')
		Assert(Type(Date()) is: 'Date')
		Assert(Type(#()) is: 'Object')
		Assert(Type(Object()) is: 'Object')
		Assert(Type(Record()) is: 'Record')
		Assert(Type(#{}) is: 'Record')
		Assert(Type([]) is: 'Record')
		Assert(Type(Type) is: 'BuiltinFunction')
		Assert(Type(function () { }) is: 'Function')
		Assert(Type({ }) is: 'Block')
		Assert(Type(Stack) is: 'Class')
		Assert(Type(new Stack) is: 'Instance')
		Assert(Type(Stack.Push) is: 'Method')
		Assert(Transaction(read:) { |t| Type(t) } is: 'Transaction')
		Assert(Cursor('tables') { |c| Type(c) } is: 'Cursor')
		Assert(WithQuery('tables') { |q| Type(q) } is: 'Query')
		}
	}