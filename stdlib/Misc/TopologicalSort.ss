// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// list = #(
//	(name: 'name1', deps: (...)),
//	(name: 'name2', deps: ('name1', ...)),
//	(name: 'name3', deps: ('name1', 'name2', ...)))
// return a new list with items sorted topologically
class
	{
	CallClass(list)
		{
		queue = Queue()
		result = Object()

		children = Object().Set_default(Object())
		counts = Object()
		for item in list
			{
			if counts.Member?(item.name)
				throw 'Find duplicate name: ' $ item.name
			counts[item.name] = item.deps.Size()
			if item.deps.Empty?()
				queue.Enqueue(item)
			else
				item.deps.Each({ children[it].Add(item) })
			}

		while queue.Count() > 0
			{
			current = queue.Dequeue()
			result.Add(current)
			for child in children[current.name]
				{
				childCount = --counts[child.name]
				Assert(childCount greaterThanOrEqualTo: 0)
				if childCount is 0
					queue.Enqueue(child)
				}
			}
		if result.Size() isnt list.Size()
			throw .analyze(counts, list)
		return result
		}

	analyze(counts, list)
		{
		analyzed = Object()
		_results = Object(circles: Object(), unknown: Object())
		for m in counts.Members().Sort!()
			if counts[m] isnt 0
				{
				.analyze2(m, list, counts, Object(), analyzed, )
				}
		return [
			Opt('Find circles:\r\n',
				_results.circles.Map({ it.Join(' -> ') }).Join('\r\n')),
			Opt('Find unknown deps:\r\n',
				_results.unknown.Map({ it.dep $ ' in ' $ it.name }).Join('\r\n'))].
			Join('\r\n').Trim()
		}

	analyze2(name, list, counts, currents, analyzed)
		{
		if false isnt pos = currents.Find(name)
			{
			_results.circles.Add(currents[pos..].Add(name))
			return
			}
		if analyzed.Member?(name)
			return
		analyzed[name] = true

		currents.Add(name)
		table = list.FindOne({ it.name is name })
		for dep in table.deps
			if not counts.Member?(dep)
				_results.unknown.Add([:name, :dep])
			else if counts[dep] isnt 0
				.analyze2(dep, list, counts, currents.Copy(), analyzed)
		}
	}