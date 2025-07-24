class
	{
	CallClass()
		{
		return Join(' | ', .heap(), .res(),
			.cursors(), .transactions(), .threads())
		}
	res()
		{
		return 'Rsrc ' $ ResourceCounts().Sum()
		}
	heap()
		{
		return 'Heap ' $
			.format(Suneido.GoMetric("/memory/classes/heap/objects:bytes")) $
			" / " $ .format(MemoryArena()) $ " mb"
		}
	format(n)
		{
		return (n / 1_000_000).Round(0)
		}
	transactions()
		{
		n = ServerEval('WorkSpaceStatus.TranSize')
		return n is 0 ? "" : 'Transactions ' $ n
		}
	TranSize()
		{
		return Database.Transactions().Size()
		}
	cursors()
		{
		n = Database.Cursors()
		return n is 0 ? "" : 'Cursors ' $ n
		}
	threads()
		{
		return 'Threads ' $ Thread.Count()
		}
	ResourceDetails()
		{
		return 'Resources: ' $ ResourceCounts().Map2({|m, v| m $ ' ' $ v }).Join(' | ')
		}
	}