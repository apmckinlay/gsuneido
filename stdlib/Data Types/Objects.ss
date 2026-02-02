// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	SizeLimit: 64000
	Has?(value) // overridden in Sequences
		{
		return .Find(value) isnt false
		}
	Empty?() // overridden in Sequences
		{
		return .Size() is 0
		}
	NotEmpty?() // overridden in Sequences
		{
		return .Size() > 0
		}
	HasNamed?() // overridden in Sequences
		{
		return .Size(named:) > 0
		}
	Map!(block)
		{
		for i in this.SafeMembers()
			this[i] = block(this[i])
		return this
		}
	Map(block) // overridden in Sequences
		{
		return Map(this, block)
		}
	Map2(block)
		{
		return .Size(named:) is 0
			? Map2(this, block)
			: Sequence(new Objects.Map2Iterator(this, .Members().Sort!(), block))
		}
	Map2Iterator: class
		{
		i: 0
		New(.ob, .mems, .block)
			{
			}
		Next()
			{
			if .i >= .mems.Size()
				return this
			m = .mems[.i++]
			return (.block)(m, .ob[m])
			}
		Dup()
			{
			return new (.Base())(.ob, .mems, .block)
			}
		Infinite?()
			{
			return false
			}
		}
	MapMembers(block)
		{
		ob = Object()
		for i in this.SafeMembers()
			ob[block(i)] = this[i]
		return ob
		}
	FlatMap(block) //TODO change to Sequence
		{
		results = Object()
		for x in this
			results.Add(@block(x))
		return results
		}
	First() // overridden in Sequences
		{
		return this[0]
		}
	Last() // overridden in Sequences
		{
		return this[.Size(list:) - 1]
		}
	MinWith(block)
		{
		iter = .Iter()
		min = iter.Next()
		Assert(min isnt iter, "cannot use MinWith on empty object")
		while (iter isnt x = iter.Next())
			if Cmp(block(x), block(min)) < 0
				min = x
		return min
		}
	MaxWith(block)
		{
		iter = .Iter()
		max = iter.Next()
		Assert(max isnt iter, "cannot use MaxWith on empty object")
		while (iter isnt x = iter.Next())
			if Cmp(block(x), block(max)) > 0
				max = x
		return max
		}
	Project(@fields)
		{
		if fields.Size(list:) is 1 and Object?(fields[0])
			fields = fields[0]
		ob = Object()
		for f in fields
			ob[f] = this[f]
		return ob
		}
	ProjectValues(@fields)
		{
		if fields.Size(list:) is 1 and Object?(fields[0])
			fields = fields[0]
		ob = Object()
		for f in fields
			ob.Add(this[f])
		return ob
		}
	ListToNamed(@fields) // e.g. #(1, 2, 3).ListToNamed(#a, #b, #c) => #(a: 1, b: 2, c: 3)
		{
		ob = Object()
		for (i = 0; i < fields.Size(); ++i)
			ob[fields[i]] = this[i]
		return ob
		}
	ListToMembers() // e.g. #(a, b, c) => #(a:, b:, c:) so .Member? is faster than .Has?
		{
		ob = Object()
		for x in this.Values(list:)
			ob[x] = true
		return ob
		}
	Flatten()
		{
		return .flatten(Object(), this)
		}
	flatten(dst, src)
		{
		for x in src
			if (Object?(x))
				.flatten(dst, x)
			else
				dst.Add(x)
		return dst
		}
	Intersects?(other)
		{
		other.Instantiate() // other MUST be an object, strings are NOT valid
		return .Any?({ other.Has?(it) })
		}
	Disjoint?(other)
		{
		return not .Intersects?(other)
		}
	Intersect(other)
		// returns list of values from this that are also in other
		// maintains order of this
		{
		other.Instantiate() // other MUST be an object, strings are NOT valid
		return .Filter({ other.Has?(it) }).Instantiate()
		}
	Difference(other)
		{
		// returns list of values from this that are not in other
		// maintains order of this
		other.Instantiate() // other MUST be an object, strings are NOT valid
		return .Filter({ not other.Has?(it) }).Instantiate()
		}
	Union(other)
		{
		result = .UniqueValues()
		for x in other
			result.AddUnique(x)
		return result
		}
	Subset?(other) // TODO rename to SubsetOf?
		{
		other.Instantiate() // other MUST be an object, strings are NOT valid
		return .Every?({ other.Has?(it) })
		}
	MergeUnion(other)
		// TODO: Should this be renamed MergeUnion!(other) or Union!(other) b/c
		//       'this' is modified? I think Union!() would be better b/c the
		//       'Merge' series of functions operate on keys, whereas the
		//       'Union' series operates on values...
		{
		.Instantiate()
		for x in other
			if not .Has?(x)
				.Add(x)
		return this
		}
	DuplicateValues()
		{
		dups = Object()
		ob = Object()
		for val in this.Values()
			{
			if ob.Has?(val)
				dups.AddUnique(val)
			ob.AddUnique(val)
			}
		return dups
		}
	FindOne(block)
		{
		for x in this
			if block(x) is true
				return x
		return false
		}
	FindIf(block)
		{
		for i in this.SafeMembers()
			if block(this[i]) is true
				return i
		return false
		}
	FindLastIf(block)
		{
		for (i = .Size(list:) - 1; i >= 0; --i)
			if block(this[i]) is true
				return i
		return false
		}
	FindAllIf(block)
		{
		list = Object()
		for m in this.SafeMembers()
			if block(this[m]) is true
				list.Add(m)
		return list
		}
	FindAll(value)
		{
		list = Object()
		for m in this.SafeMembers()
			if this[m] is value
				list.Add(m)
		return list
		}
	MembersIf(block)
		{
		mems = Object()
		for m in this.Members()
			if block(m) is true
				mems.Add(m)
		return mems
		}
	Val_or_func(member)
		{
		x = this[member]
		if Function?(x)
			x = x()
		return x
		}
	// Merge will overwrite unnamed (as well as named) members with ob's members
	Merge(ob)
		{
		for member in ob.SafeMembers()
			this[member] = ob[member]
		return this
		}
	MergeNew(ob) // won't overwrite existing members
		{
		for member in ob.SafeMembers()
			if not .Member?(member)
				this[member] = ob[member]
		return this
		}
	SafeMembers() // overridden by Records
		{
		return .Members()
		}
	Replace(oldvalue, newvalue)
		{
		if oldvalue isnt newvalue
			while false isnt m = .Find(oldvalue)
				this[m] = newvalue
		return this
		}
	JoinCSV(fields = false)
		{
		line = ""
		for (field in (fields is false ? .SafeMembers() : fields))
			{
			value = this[field]
			if String?(value)
				line $= '"' $ value.Replace('"', '""') $ '"' $ ','
			else
				line $= Display(value) $ ','
			}
		return line[.. -1]
		}
	Sum()
		{
		sum = 0
		for x in this
			sum += x
		return sum
		}
	SumWith(block)
		{
		sum = 0
		for x in this
			sum += block(x)
		return sum
		}
	Fold(val, block)
		{
		for x in this
			val = block(val, x)
		return val
		}
	Reduce(block)
		{
		iter = .Iter()
		if iter is x = iter.Next()
			return Nothing() // no return value for empty sequence
		while iter isnt y = iter.Next()
			x = block(x, y)
		return x
		}
	Remove1(x) // remove first occurrence
		{
		// faster than Remove and won't get object-modified-during-iteration
		if false isnt m = .Find(x)
			.Delete(m)
		}
	Remove(@values)
		// Remove all occurrences of values
		{
		if 0 is n = values.Size()
			return this
		members = (n is 1 ? .FindAll(values[0]) : .FindAllIf({ values.Has?(it) }))
		for m in members.Reverse!()
			.Delete(m)
		return this
		}
	RemoveIf(block)
		{
		for m in .FindAllIf(block).Reverse!()
			.Delete(m)
		return this
		}
	DeleteIf(block)
		{
		for m in .MembersIf(block).Reverse!()
			.Delete(m)
		return this
		}
	Without(@values)
		{
		return .Copy().Remove(@values)
		}
	WithoutFields(@fields)
		{
		return .Copy().Delete(@fields)
		}
	Swap(idx1, idx2)
		{
		temp = this[idx1]
		this[idx1] = this[idx2]
		this[idx2] = temp
		return this
		}
	AddUnique(value)
		{
		if (not .Has?(value))
			.Add(value)
		return this
		}
	Filter(block)
		{
		return Filter(this, block)
		}
	Zip()
		{
		ob = Object()
		if (.Empty?())
			return ob
		n = .Map(#Size).Min()
		for (i = 0; i < n; ++i)
			ob.Add(.Map({|seq| seq[i] }).Add()) // Add forces instantiation
		return ob
		}
	UniqueValues()
		{
		ob = Object()
		for x in this
			ob.AddUnique(x)
		return ob
		}
	Append(ob)
		{
		if Seq?(ob)
			for x in ob // avoid instantiating
				.Add(x)
		else
			{
			i = 0
			for m in ob.SafeMembers()
				if m is i++
					.Add(ob[m])
				else
					this[m] = ob[m]
			}
		return this
		}
	Every?(block)
		{
		for x in this
			if not block(x)
				return false
		return true
		}
	Any?(block)
		{
		for x in this
			if block(x)
				return true
		return false
		}
	HasIf?(block)
		{
		for x in this
			if block(x)
				return true
		return false
		}
	Shuffle!()
		{
		n = this.Size(list:)
		for (i = 0; i < n - 1; ++i)
			this.Swap(i, i + Random(n - i))
		return this
		}
	AddMany!(value, n)
		{
		while n-- > 0
			.Add(value)
		return this
		}
	Each(block) // TODO: make builtin
		{
		for x in this
			try
				block(x)
			catch (e, "block:")
				if e is "block:break"
					break
				// else block:continue ... so continue
		return this
		}
	Each2(block) // passes member and value to block
		{
		for m in .SafeMembers()
			try
				block(m, this[m])
			catch (e, "block:")
				if e is "block:break"
					break
				// else block:continue ... so continue
		return this
		}
	SortWith!(block)
		{
		return .Sort!({|x, y| block(x) < block(y) })
		}
	Extract(member, x = #extract_no_default) //TODO rename to Extract!
		{
		if .Member?(member)
			{
			x = this[member]
			.Delete(member)
			}
		else if x is #extract_no_default
			throw "missing: " $ member
		return x
		}
	Count(value = #(0))
		{
		if Same?(value, #(0)) // depends on compiler sharing constants
			return .Size()
		n = 0
		for x in this
			if x is value
				++n
		return n
		}
	CountIf(block)
		{
		n = 0
		for x in this
			if block(x)
				++n
		return n
		}
	Sorted?(block = false)
		{
		n = .Size(list:)
		if block is false
			{
			for (i = 1; i < n; ++i)
				if Cmp(this[i], this[i-1]) < 0
					return false
			}
		else
			{
			for (i = 1; i < n; ++i)
				if block(this[i], this[i-1])
					return false
			}
		return true
		}
	Trim!(@values)
		{
		while .Member?(0) and values.Has?(this[0])
			.Delete(0)
		while .Member?(0) and values.Has?(.Last())
			.Delete(.Size(list:) - 1)
		return this
		}
	HasNonEmptyMember?(members)
		{
		return members.Any?({ this.Member?(it) and this[it] isnt "" })
		}
	BinarySearch?(value)
		{
		// this must be sorted
		i = .BinarySearch(value)
		return this.Member?(i) and this[i] is value
		}
	GetDifferences(ob)
		{
		ob = ob.Copy()
		diffs = Object()
		for m in this.SafeMembers()
			{
			if not ob.Member?(m)
				diffs[m] = "does not exist in ob"
			else if this[m] isnt ob[m]
				diffs[m] = Display(this[m]) $ " isnt " $ Display(ob[m])
			ob.Delete(m)
			}
		for m in ob.SafeMembers()
			diffs[m] = "does not exist in this"
		return diffs
		}
	GetInit(member, block)
		{
		// handle multiple threads initializing and deleting
		marker = #(1)
		for ..4 /*= unlikely to conflict many times */
			{
			if not Same?(marker, x = .GetDefault(member, marker))
				return x
			newValue = Function?(block) ? block() : block
			// multiple threads may reach here
			if .CompareAndSet(member, newValue)
				// only one thread will reach here
				return newValue
			}
		throw "GetInit: too many thread conflicts"
		}
	Take(n) // overridden in Sequences
		{
		return this[:: n]
		}
	Drop(n) // overridden in Sequences
		{
		return this[n ..]
		}
	Nth(n) // overridden in Sequences
		{
		return this[n]
		}
	Concat(@iterables)
		{
		iterables.Add(this at: 0)
		return Concat(@iterables)
		}
	AddTo(ob)
		{
		// avoid instantiating if this is a sequence
		for x in this
			ob.Add(x)
		}
	Instantiate() // overridden in Sequences
		{
		return this
		}
	Grep(regex, block = false)
		{
		return Grep(this, regex, block)
		}
	DeepCopy(nesting = 0)
		{
		nesting = .checkNesting('DeepCopy', nesting)
		return .Copy().Map!({ Object?(it) ? it.DeepCopy(:nesting) : it })
		}
	DeepReplaceReferences!(oldName, newName, nesting = 0)
		{
		nesting = .checkNesting('DeepReplaceReferences!', nesting)
		for m in this.Members().Copy()
			{
			if Object?(this[m])
				this[m].DeepReplaceReferences!(oldName, newName, :nesting)
			if m is oldName
				this[newName] = this.Extract(m)
			else if this[m] is oldName
				this[m] = newName
			}
		}
	checkNesting(methodName, nesting)
		{
		if ++nesting > 10 /*= max depth*/
			throw "object." $ methodName $ " nesting > 10"
		return nesting
		}
	RandVal()
		{
		return this[Random(this.Size(list:))]
		}
	EqualSet?(that)
		{
		if .Size() isnt that.Size()
			return false
		that = that.Copy()
		for x in this
			if false isnt i = that.Find(x)
				that.Delete(i)
			else
				return false
		return that.Empty?()
		}
	}
