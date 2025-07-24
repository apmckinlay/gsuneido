// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_wrapText()
		{
		fn = KnowledgeBase.KnowledgeBase_wrapText

		Assert(fn('', 5) is: [''])
		Assert(fn('hello', 10) is: ['hello'])
		Assert(fn('abcdefghijklmnopqrstuvwxyz', 5)
			is: ["abcde", "fghij", "klmno", "pqrst", "uvwxy", "z"])
		Assert(
			fn("This is a \n test sentence. \r\n It should be wrapped \r properly.", 10)
			is: ["This is a", "test", "sentence.", "It should", "be wrapped", "properly."]
			)
		Assert(
			fn("The quick \t\t brown fox jumps over the lazy dog.", 8)
			is: ["The", "quick", "brown", "fox", "jumps", "over the", "lazy", "dog."])
		Assert(
			fn("This sentence \n \n contains spaces and \r\n line breaks.", 10)
			is: ["This", "sentence", "contains", "spaces and", "line", "breaks."])
		}

	fakeOpenAI: class
		{
		callId: 0
		New(.log) { }
		Embeddings(inputs)
			{
			res = Object()
			for input in inputs
				{
				if input is 'throw'
					throw 'test throw'
				if input.Tr('\r\n', ' ').Split(' ').Size() > 3
					{
					.log.Add([input, 'throw'])
					throw 'OpenAI Embeddings: Please reduce your prompt'
					}
				.log.Add([input, .callId])
				res.Add([.callId++])
				}
			return res
			}
		}
	fakeVdb: class
		{
		New(.log) { }
		Upsert(vectors)
			{
			.log.Append(vectors)
			}
		}
	Test_update()
		{
		spy = .SpyOn(SuneidoLog).Return(true)
		log1 = Object()
		log2 = Object()
		openAI = (.fakeOpenAI)(log1)
		vdb = (.fakeVdb)(log2)

		fn = KnowledgeBase.KnowledgeBase_update

		Assert(fn([], openAI, vdb) is: 0)

		outputs = [[text: 'a b c', id: 'id 1', metadata: #(source: 'test')]]
		Assert(fn(outputs, openAI, vdb) is: 1)
		Assert(log1 is: [['a b c', 0]])
		Assert(log2 is: [[metadata: #(source: "test"), id: "id 1", values: #(0)]])

		log1.Delete(all:)
		log2.Delete(all:)
		outputs = [[text: 'a b c d', id: 'id 2', metadata: #(source: 'test')]]
		Assert(fn(outputs, openAI, vdb) is: 2)
		Assert(log1 is: [['a b c d', 'throw'], ['a b c', 1], ['d', 2]])
		Assert(log2 is: [
			[metadata: #(source: "test"), id: "id 2.0", values: #(1)],
			[metadata: #(source: "test"), id: "id 2.1", values: #(2)]])

		log1.Delete(all:)
		log2.Delete(all:)
		outputs = [
			[text: 'a b c d', id: 'id 3', metadata: #(source: 'test')],
			[text: 'a b c\nd e f\ng h i', id: 'id 4', metadata: #(source: 'test')]]
		Assert(fn(outputs, openAI, vdb) is: 5)
		Assert(log1 is: [
			['a b c d', 'throw'],
			['a b c\nd e f', 'throw'],
			['a b c', 3],
			['d e f', 4],
			['g h i', 5],
			['a b c d', 'throw'],
			['a b c', 6],
			['d', 7]])
		Assert(log2 is: [
			[metadata: #(source: "test"), id: "id 4.0.0", values: #(3)],
			[metadata: #(source: "test"), id: "id 4.0.1", values: #(4)],
			[metadata: #(source: "test"), id: "id 4.1", values: #(5)],
			[metadata: #(source: "test"), id: "id 3.0", values: #(6)],
			[metadata: #(source: "test"), id: "id 3.1", values: #(7)]])

		log1.Delete(all:)
		log2.Delete(all:)
		outputs = [
			[text: 'a b c d', id: 'id 5', metadata: #(source: 'test')],
			[text: 'throw', id: 'id 6', metadata: #(source: 'test')]]
		Assert(fn(outputs, openAI, vdb) is: 2)
		callLogs = spy.CallLogs()
		Assert(callLogs isSize: 1)
		Assert(callLogs[0].params is: #('id 6'))
		Assert(log1 is: [['a b c d', 'throw'], ['a b c', 8], ['d', 9]])
		Assert(log2 is: [
			[metadata: #(source: "test"), id: "id 5.0", values: #(8)],
			[metadata: #(source: "test"), id: "id 5.1", values: #(9)]])
		}
	}