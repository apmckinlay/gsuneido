// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// this is tests for methods defined in Objects
// ObjectTest is for built-in methods
// SuJsWebTest
Test
	{
	Test_ListToNamed()
		{
		Assert(#().ListToNamed() is: #())
		Assert(#(1, 2, 3).ListToNamed() is: #())
		Assert(#(1, 2, 3).ListToNamed(#a, #b, #c) is: #(a: 1, b: 2, c: 3))
		Assert(#(1, 2, 3).ListToNamed(#a, #b) is: #(a: 1, b: 2))
		}
	Test_ListToMembers()
		{
		Assert(#().ListToMembers() is: #())
		Assert(#(a, b, c).ListToMembers() is: #(a:, b:, c:))
		Assert(#(a, b, c, d:, e:).ListToMembers() is: #(a:, b:, c:))
		}
	Test_FlatMap()
		{
		twice = function (x) { return [x, x] }
		Assert([].FlatMap(twice) is: [])
		Assert([1, 2].FlatMap(twice) is: [1, 1, 2, 2])
		}
	Test_Each()
		{
		s = ""
		[1, 2, 3].Each()
			{|x| s $= x }
		Assert('123' is: s)

		s = ""
		[1, 2, 3, 4, 5].Each()
			{|x|
			if x is 2
				continue
			s $= x
			if x is 4
				break
			}
		Assert('134' is: s)
		}
	Test_Each2()
		{
		test = function (@ob)
			{
			result = Object()
			ob.Each2 {|k,v| result[k] = v }
			Assert(result is: ob)
			}
		test()
		test(1)
		test(a: 1)
		test(1, 2, a: 3, b: 4)
		}
	Test_Extract()
		{
		.extract([], 				#a, false)
		.extract([a: 123],			#b, false)
		.extract([a: 123],			#a,	123)
		.extract([a: 12, b: 34],	#b,	34)
		Assert([a: 123].Extract(#a) is: 123)
		Assert({ [].Extract(#a) } throws: "missing: a")
		}
	extract(ob, mem, expect)
		{
		x = ob.Extract(mem, false)
		Assert(expect is: x)
		Assert(ob hasntMember: mem)
		}
	Test_FindAllIf()
		{
		gt5 = function (x) { x > 5 }
		Assert([].FindAllIf(gt5) is: [])
		Assert([1, 2, a: 3, b: 4].FindAllIf(gt5) is: [])
		Assert([1, 6, a: 3, b: 8].FindAllIf(gt5) is: #(1, b))
		}
	Test_RemoveIf()
		{
		gt5 = function (x) { x > 5 }
		Assert([].RemoveIf(gt5) is: [])
		x = [1, 2, a: 3, b: 4]
		Assert(x.RemoveIf(gt5) is: x)
		Assert([1, 6, a: 3, b: 8].RemoveIf(gt5) is: [1, a: 3])
		}
	Test_FindAll()
		{
		Assert([1, 2, a: 2, b: 3].FindAll(5) is: [])
		Assert([1, 2, a: 2, b: 3].FindAll(2) is: #(1, a))
		}
	Test_Remove()
		{
		x = [1, 2, a: 2, b: 3]
		Assert(x.Remove() is: x)
		Assert(x.Remove(5) is: x)
		Assert(x.Remove(5, 6) is: x)
		Assert(x.Remove(2) is: [1, b: 3])
		Assert(x.Remove(7, 2, 8) is: [1, b: 3])
		Assert(x.Remove(0, 1, 2, 3) is: [])

		Assert([1, '', 2, '', '', 3, 4].Remove('') is: [1, 2, 3, 4])
		}
	Test_WithoutFields()
		{
		x = []
		Assert(x.WithoutFields() is: x)
		Assert(x.WithoutFields(#a, #b) is: x)
		x = [a: 12, b: 34, c: 45]
		Assert(x.WithoutFields() is: x)
		Assert(x.WithoutFields(#d, #e) is: x)
		Assert(x.WithoutFields(#b) is: [a: 12, c: 45])
		Assert(x.WithoutFields(#a, #b, #c) is: [])
		}
	Test_Sorted?()
		{
		Assert(#().Sorted?())
		Assert(#(12).Sorted?())
		Assert(#(12, 12, 12).Sorted?())
		Assert(#(12, 34, 56).Sorted?())
		Assert(#(12, 34, 56).Sorted?({|x,y| x < y }))
		Assert(#(65, 43, 21).Sorted?({|x,y| x > y }))
		Assert(not #(65, 43, 21).Sorted?())
		Assert(not #(12, 34, 0).Sorted?())
		Assert(#((name: "Fred", age: 35), (name: Sue, age: 25)).Sorted?(By(#name)))
		Assert(not #((name: "Fred", age: 35), (name: Sue, age: 25)).Sorted?(By(#age)))
		}
	Test_Trim()
		{
		Assert(#().Trim!() is: #())
		Assert(#().Trim!('') is: #())
		Assert(#(named: '').Trim!('') is: #(named: ''))
		Assert(#(123, 456, 789).Trim!('') is: #(123, 456, 789))
		Assert(['', 123, 456, 789, '', named:].Trim!('') is: #(123, 456, 789, named:))
		Assert(['', ''].Trim!('') is: #())
		Assert(#(123, '', 789).Trim!('') is: #(123, '', 789))
		}
	Test_Union()
		{
		n = #()
		d = #(a, b, a, b)
		ab = #(a, b)
		abc = #(a,b,c)
		Assert(n.Union(n) is: n)
		Assert(d.Union(n) is: ab)
		Assert(n.Union(d) is: ab)
		Assert(d.Union(d) is: ab)
		Assert(d.Union(abc) is: abc)
		Assert(abc.Union(d) is: abc)
		}
	Test_HasNonEmptyMember?()
		{
		fields = #()
		Assert(#().HasNonEmptyMember?(fields) is: false)
		fields = #(a, b)
		Assert(#().HasNonEmptyMember?(fields) is: false)
		Assert(#(c: 1).HasNonEmptyMember?(fields) is: false)
		Assert(#(b: '').HasNonEmptyMember?(fields) is: false)
		Assert([a: '', b: 1].HasNonEmptyMember?(fields) is: true)
		Assert(#(a: '', b: 1).HasNonEmptyMember?(fields) is: true)
		}
	Test_Project()
		{
		Assert(#().Project(#()) is: #())
		Assert(#(a:1, b:2, c:3, d:4).Project(#(a, c)) is: #(a:1, c:3))
		Assert(#(a:1, b:2, c:3, d:4).Project(#a, #c) is: #(a:1, c:3))
		Assert(#(a:1, b:2, c:3, d:4).Project(#(a, b, c, d)) is: #(a:1, b:2, c:3, d:4))

		Assert(Type(#(a: 1).Project(#(a))) is: 'Object')
		Assert(#(a: 1, b: 2, c: 3).Project(#(a, c)) is: #(a: 1, c: 3))
		Assert({ #(a: 1, b: 2, c: 3).Project(#(a, c, d)) } throws: 'member not found')

		}
	Test_ProjectValues()
		{
		Assert(#().ProjectValues(#()) is: #())
		Assert(#(a:1, b:2, c:3, d:4).ProjectValues(#(a, c)) is: #(1, 3))
		Assert(#(a:1, b:2, c:3, d:4).ProjectValues(#a, #c) is: #(1, 3))
		Assert(#(a:1, b:2, c:3, d:4).ProjectValues(#(a, b, c, d)) is: #(1, 2, 3, 4))
		}
	Test_FindOne()
		{
		Assert(#((a: 1, b: 2), (a: 6, b: 7)).FindOne({ it.b.Odd?() })
			is: #(a: 6, b: 7))
		Assert(#((a: 1, b: 2), (a: 6, b: 7)).FindOne({ it.b > 9 })
			is: false)
		}
	Test_Intersect()
		{
		Assert(#(1, 2, 3, 4, 5).Intersect(#(4, 5, 6, 7, 8)) is: #(4, 5))
		Assert(#(1, 2, 3, 4, 5).Intersect(#(1, 2, 3, 4, 5)) is: #(1, 2, 3, 4, 5))
		Assert(#(1, 2, 3, 4, 5).Intersect(#(6, 7, 8, 9, 0)) is: #())
		}
	Test_Difference()
		{
		Assert(#(1, 2, 3, 4, 5).Difference(#(4, 5, 6, 7, 8)) is: #(1, 2, 3))
		Assert(#(1, 2, 3, 4, 5).Difference(#(1, 2, 3, 4, 5)) is: #())
		Assert(#(1, 2, 3, 4, 5).Difference(#(6, 7, 8, 9, 0)) is: #(1, 2, 3, 4, 5))
		}
	Test_DuplicateValues()
		{
		Assert(#().DuplicateValues() is: #())
		Assert(#(1, 2, 3, 4).DuplicateValues() is: #())
		Assert(#(1, 2, 3, 1).DuplicateValues() is: #(1))
		Assert(#(1, 2, 3, 2).DuplicateValues() is: #(2))
		Assert(#(1, 1, 1, 1).DuplicateValues() is: #(1))
		Assert(#(1, 1, 2, 2).DuplicateValues() is: #(1, 2))
		Assert(#(1, 1, 2, 2, 3, 4, 5, 2, 3, 6).DuplicateValues() is: #(1, 2, 3))
		Assert(#(a: 1, b: 1, c: 2, d: 3).DuplicateValues() is: #(1))
		}
	Test_JoinCSV()
		{
		record = #('one', 'two', 3, 'f,o,u,r', '"five"')
		// Join with implicit fields
		line = record.JoinCSV()
		Assert(line  is: '"one","two",3,"f,o,u,r","""five"""')
		// Join with explicit fields
		line = record.JoinCSV(#(4, 3, 2, 1, 0))
		Assert(line  is: '"""five""","f,o,u,r",3,"two","one"')
		}
	Test_Fold()
		{
		Assert(#(1, 2, 3).Fold(0, {|sum,x| sum + x }) is: 6)
		Assert(#(a, b, c).Fold('', {|s,t| s $= t }) is: "abc")
		}
	Test_Reduce()
		{
		Assert(#(123).Reduce(Min) is: 123)
		Assert(#(3, 5, 4, 6, 2).Reduce(Min) is: 2)
		Assert(#(3, 5, 4, 6, 2).Reduce(Max) is: 6)
		Assert(#(3, 5, 4, 6, 2).Reduce({|x, y| x + y }) is: 20)
		}
	Test_Filter()
		{
		t = {|unused| true }
		f = {|unused| false }
		cond = {|x| x < 5 }
		Assert(#().Filter(t) is: #())
		Assert(#().Filter(f) is: #())
		Assert(#().Filter(cond) is: #())
		Assert(Seq(10).Filter(t) is: Seq(10))
		Assert(Seq(10).Filter(f) is: #())
		Assert(Seq(10).Filter(cond) is: Seq(5))
		}
	Test_UniqueValues()
		{
		Assert(#().UniqueValues() is: #())
		Assert(#(1, 2, 3).UniqueValues() is: #(1, 2, 3))
		Assert(#(1, 2, 2, 3, 1).UniqueValues() is: #(1, 2, 3))
		}
	Test_Zip()
		{
		Assert(#().Zip() is: #())
		Assert(#(#(), #()).Zip() is: #())
		Assert(#(#(1, 2, 3), #(11, 22, 33, 44), #(111, 222, 333)).Zip()
			is: #((1, 11, 111) (2, 22, 222) (3, 33, 333)))
		}
	Test_Min()
		{
		Assert(#(1).Min() is: 1)
		Assert(#(1 50 23).Min() is: 1)
		Assert(#(50 23 1).Min() is: 1)
		}
	Test_Max()
		{
		Assert(#(1).Max() is: 1)
		Assert(#(1 50 23).Max() is: 50)
		Assert(#(50 23 1).Max() is: 50)
		}
	Test_MaxWith()
		{
		Assert(#(one, two, three, four, five).MaxWith(#Size) is: "three")
		}
	Test_Append()
		{
		Assert(Object(1, 2, a: 3, b: 4).Append(#(2.3, 2.6, d: 5, e: 6))
			is: #(1, 2, 2.3, 2.6, a: 3, b: 4, d: 5, e: 6))
		}
	Test_Every?()
		{
		Assert(#().Every?(#Odd?))
		Assert(not #(1, 2, 3).Every?(#Odd?))
		Assert(#(1, 11, 3).Every?(#Odd?))
		}
	Test_Any?()
		{
		Assert(not #().Any?(#Even?))
		Assert(#(1, 2, 3).Any?(#Even?))
		Assert(not #(1, 11, 3).Any?(#Even?))
		}
	Test_BinarySearch?()
		{
		ob = #(1, 2, 3)
		Assert(ob.BinarySearch?(2) 		is: true)
		Assert(ob.BinarySearch?(3) 		is: true)
		Assert(ob.BinarySearch?(1.5) 	is: false)
		Assert(ob.BinarySearch?(4) 		is: false)
		}
	Test_GetDifferences()
		{
		Assert(#().GetDifferences(#()) is: #())
		Assert(#(a: 1, b: 2).GetDifferences(#(b: 2, a: 1)) is: #())
		diffs = #(a:1 b:2 c: 10).GetDifferences(#(a:3 b:4, z: 99))
		Assert(diffs is: #(c: "does not exist in ob", b: "2 isnt 4",
			z: "does not exist in this", a: "1 isnt 3"))
		}
	Test_PopFirst()
		{
		ob = Object()
		Assert(Same?(ob.PopFirst(), ob))
		test = function (@ob)
			{
			first = ob[0]
			after = ob.Copy().Delete(0)
			Assert(ob.PopFirst() is: first)
			Assert(ob is: after)
			}
		test(1)
		test(1, 2, 3)
		test(1, 2, a: 3, b: 4)
		}
	Test_PopLast()
		{
		ob = Object()
		Assert(Same?(ob.PopLast(), ob))
		test = function (@ob)
			{
			last = ob.Last()
			after = ob.Copy().Remove(last)
			Assert(ob.PopLast() is: last)
			Assert(ob is: after)
			}
		test(1)
		test(1, 2, 3)
		test(1, 2, a: 3, b: 4)
		}
	Test_GetInit()
		{
		ob = Object(a: 123)
		Assert(ob.GetInit(#a, 456) is: 123)
		Assert(ob is: #(a: 123))
		Assert(ob.GetInit(#b, 456) is: 456)
		Assert(ob is: #(a: 123, b: 456))
		n = 0
		Assert(ob.GetInit(#a, { ++n }) is: 123)
		Assert(ob.GetInit(#d, { ++n; 789 }) is: 789)
		Assert(n is: 1)
		}
	Test_DeepCopy()
		{
		Assert(#().DeepCopy() is: #())
		x = #()
		Assert({ x.y = 1 } throws: 'readonly')
		Assert(x.DeepCopy().Add(1) is: #(1))
		x = #(1, #(2, b: 3), c: 4, d: #(5, e: 6))
		y = x.DeepCopy()
		Assert(y is: x)
		Assert(y.Add(8) is: #(1, #(2, b: 3), 8, c: 4, d: #(5, e: 6)))
		y.d.f = 9
		Assert(y is: #(1, #(2, b: 3), 8, c: 4, d: #(5, e: 6, f: 9)))
		x = #((((((((((((1))))))))))))
		Assert({ x.DeepCopy() } throws: "nesting")
		}
	Test_DeepReplaceReferences()
		{
		ob = #()
		ob.DeepReplaceReferences!('a', 'b')
		Assert(ob is: #())

		ob = Object(a: 1, b: 2, c: 3)
		ob.DeepReplaceReferences!('a', 'd')
		Assert(ob is: #(d: 1, b: 2, c: 3))

		ob = Object(a: 1, b: 2, c: 3, d: Object(x: 9, y: 8, z: 7) z: 99)
		ob.DeepReplaceReferences!('z', 'zz')
		Assert(ob is: Object(a: 1, b: 2, c: 3, d: Object(x: 9, y: 8, zz: 7) zz: 99))

		ob = Object(a: 1, b: 'z', c: 3, d: Object(x: 'z', y: 8, z: 7) z: 99)
		ob.DeepReplaceReferences!('z', 'zz')
		Assert(ob is: Object(a: 1, b: 'zz', c: 3, d: Object(x: 'zz', y: 8, zz: 7) zz: 99))
		x = #((((((((((((1))))))))))))
		Assert({ x.DeepReplaceReferences!('z', 'zz') } throws: "nesting")
		}

	Test_Val_or_func()
		{
		ob = Object(a: 'a', b: function(){ return 'b' }, c: { 'c' })
		Assert(ob.Val_or_func(#a) is: 'a')
		Assert(ob.Val_or_func(#b) is: 'b')
		Assert(ob.Val_or_func(#c) is: 'c')

		cl = class
			{
			D: 'D'
			E()
				{
				return 'E'
				}
			}
		c = cl()
		Assert(c.Val_or_func(#D) is: 'D')
		Assert(c.Val_or_func(#E) is: 'E')
		}
	Test_Sum()
		{
		Assert(#().Sum() is: 0)
		Assert(#(1, 2, a: 3, b: 4).Sum() is: 10)
		}
	Test_SumWith()
		{
		Assert(#().SumWith({ it }) is: 0)
		Assert(#(1, 2, 3).SumWith({ it * it }) is: 14)
		}
	Test_EqualSet?()
		{
		Assert(#().EqualSet?(#()))
		Assert(not #(1).EqualSet?(#()))
		Assert(not #().EqualSet?(#(1)))
		Assert(not #(1, 2, 2).EqualSet?(#(2, 1, 1)))
		Assert(#(1, 2, 2, 3).EqualSet?(#(2, 1, 3, 2)))
		}
	Test_Count()
		{
		x = #(1, 2, a: 3, b: 2)
		Assert(x.Count() is: 4)
		Assert(x.Count(9) is: 0)
		Assert(x.Count(1) is: 1)
		Assert(x.Count(2) is: 2)
		}
	Test_Map2()
		{
		// named members, sorted
		Assert(#(a: 1, b: 2, c: 3).Map2({|m, v| m $ "->" $ v })
			is: #("a->1", "b->2", "c->3"))
		}
	Test_Concat()
		{
		Assert(#(a, b, c).Concat() is: #(a, b, c))

		seq = Seq(3)
		Assert(#(a, b, c).Concat(#()) is: #(a, b, c))
		Assert(#(a, b, c).Concat(seq) is: #(a, b, c, 0, 1, 2))
		Assert(#(a, b, c).Concat(seq, #(a, b, c)) is: #(a, b, c, 0, 1, 2, a, b, c))

		Assert(seq.Concat() is: #(0, 1, 2))
		Assert(seq.Concat(#(a, b, c)) is: #(0, 1, 2, a, b, c))
		Assert(seq.Concat(#(a, b, c), seq) is: #(0, 1, 2, a, b, c, 0, 1, 2))
		}
	}
