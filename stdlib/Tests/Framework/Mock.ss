// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// Mockito style mocks/fakes
class
	{
	UseDeepEquals: true
	New(.cls = false)
		{
		.When = new .stubber(this)
		.Verify = new .verifier(this)
		.calls = Object()
		.returns = Object()
		.throws = Object()
		.dos = Object()
		.mockId = Timestamp()
		if cls isnt false
			.loadMembers(cls)
		}
	loadMembers(cls)
		{
		for m in cls.Members(all:)
			if Type(cls[m]) isnt 'function'
				this[m] = cls[m]
		}
	StubReturn(call, result)
		{
		.returns[call] = result
		}
	StubCallThrough(call)
		{
		.dos[call] = .callThrough
		}
	callThrough(call)
		{
		Assert(.cls isnt: false, msg: 'class is required for Mock CallThrough')
		newCall = call.Copy()
		newCall[0] = .cls[call[0]]
		_callThrough? = .mockId
		return .Eval(@newCall)
		}
	StubThrow(call, exception)
		{
		.throws[call] = exception
		}
	StubDo(call, block)
		{
		.dos[call] = block
		}
	Default(@call)
		{
		call = .HandleIfPrivateMethod(call)
		.calls.Add(call)
		if false isnt pattern = .find(.returns, call)
			{
			returns = .returns[pattern]
			result = returns[0]
			if returns.Size() > 1
				returns.Delete(0)
			return Type(result) is 'Block' ? result() : result
			}
		if false isnt pattern = .find(.throws, call)
			throw .throws[pattern]
		if false isnt pattern = .find(.dos, call)
			return (.dos[pattern])(:call)

		callThrough? = false
		try callThrough? = _callThrough?
		if callThrough? is .mockId
			return .callThrough(call)
		}
	find(patterns, call)
		{
		anyArgsPattern = false
		for pattern in patterns.Members()
			{
			result = .match(call, pattern)
			if result is true
				return pattern
			if result is 'anyArgs'
				anyArgsPattern = pattern
			}
		return anyArgsPattern
		}
	CountForVerify(call_pattern)
		{
		call_pattern = .HandleIfPrivateMethod(call_pattern)
		return .calls.CountIf({|call| .match(call, call_pattern) isnt false })
		}
	match(call, call_pattern)
		{
		if call[0] isnt call_pattern[0] // method name
			return false
		if .anyArgs(call_pattern)
			return 'anyArgs'
		if call.Size() isnt call_pattern.Size()
			return false
		for m in call.Members()
			if not call_pattern.Member?(m) or not .match1(call[m], call_pattern[m])
				return false
		return true
		}
	anyArgs(call_pattern)
		{
		return call_pattern.Size() is 2 and call_pattern[1] is #(anyArgs:)
		}
	match1(arg, arg_pattern)
		{
		if arg_pattern is #(any:)
			return true
		if .matcher?(arg_pattern)
			return .matcherMatch(arg, arg_pattern)
		else
			return arg is arg_pattern // default to exact match
		}
	matcher?(arg_pattern) // like [is: 0] or [startsWith: "x"]
		{
		if Instance?(arg_pattern)
			return false
		return Object?(arg_pattern) and
			arg_pattern.Size() is 1 and
			arg_pattern.Size(list:) is 0 and
			.matcherName?(arg_pattern.Members()[0])
		}
	matcherMatch(arg, arg_pattern)
		{
		name = arg_pattern.Members()[0]
		matcher = Global('Matcher_' $ name)
		return matcher.Match(arg, arg_pattern[name])
		}
	matcherName?(name)
		{
		try
			{
			Global('Matcher_' $ name)
			return true
			}
		catch
			return false
		}

	HandleIfPrivateMethod(call)
		{
		if not call[0].Capitalized?()
			{
			Assert(.cls isnt: false, msg: 'class is required for private method')
			call = call.Copy()
			clName = String(.cls).BeforeFirst(" ")
			call[0] = clName $ "_"  $ call[0]
			}
		return call
		}

	verifier: class
		{
		New(.mock, .minTimes = 1, .maxTimes = 1)
			{
			}
		Never()
			{
			return this.Base()(.mock, 0, 0)
			}
		Times(times)
			{
			return this.Base()(.mock, times, times)
			}
		AtLeast(times)
			{
			return this.Base()(.mock, times, 999999) /*= max number of times*/
			}
		AtMost(times)
			{
			return this.Base()(.mock, 0, times)
			}
		Default(@call)
			{
			n = .mock.CountForVerify(call)
			if .minTimes is 0 and .maxTimes is 0 and n isnt 0
				throw "should not have been invoked: " $ .display(call)
			if .minTimes is 1 and .maxTimes is 1 and n is 0
				throw "wanted but not invoked: " $ .display(call)
			if .minTimes is .maxTimes and n isnt .minTimes
				throw "wanted " $ .minTimes $ " calls, but got " $ n
			if n < .minTimes
				throw "wanted at least " $ .minTimes $ " calls, but got " $ n
			if n > .maxTimes
				throw "wanted at most " $ .maxTimes $ " calls, but got " $ n
			}
		display(call)
			{
			return call[0] $ Display(call[1..])[1..].Replace(': true]', ':]')
			}
		}

	stubber: class
		{
		New(.mock)
			{
			}
		Default(@call)
			{
			call = .mock.HandleIfPrivateMethod(call)
			return new .aCall(.mock, call)
			}
		aCall: class
			{
			New(.mock, .call)
				{
				}
			Return(@results)
				{
				.mock.StubReturn(.call, results)
				.mock[.call[0]] = {|@args| .mock[.call[0]](@args) }
				}
			CallThrough()
				{
				.mock.StubCallThrough(.call)
				}
			Do(block)
				{
				.mock.StubDo(.call, block)
				.mock[.call[0]] = {|@args| .mock[.call[0]](@args) }
				}
			Throw(exception)
				{
				.mock.StubThrow(.call, exception)
				}
			}
		}
	}