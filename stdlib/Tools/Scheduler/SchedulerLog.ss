// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
class
	{
	Import(folder, company)
		{
		.ensureLogTable()

		for file in Dir(folder $ 'scheduler*.*')
			{
			for line in GetFile(folder $ file).Lines()
				{
				.readline(line, company)
				}
			}
		}

	ensureLogTable()
		{
		Database(`ensure scheduler_log (schedlog_num, schedlog_company,
			schedlog_date, schedlog_type, schedlog_task)
			index (schedlog_task, schedlog_date, schedlog_company)
			key (schedlog_num)`)
		}

	readline(line, company)
		{
		date = Date(line.BeforeFirst('\t'))
		taskLine = line.AfterFirst('SchedulerLastProcess')
		type = taskLine.Prefix?('Started ') ? 'start' : 'end'
		task = taskLine.AfterFirst(' ').BeforeFirst('<#booklogsSeparator')[..-1]
		QueryOutput('scheduler_log', [schedlog_num: Timestamp(),
			schedlog_company: company, schedlog_type: type,
			schedlog_date: date, schedlog_task: task])
		}


	Summarize(task)
		{
		overlaps = Object()
		curCo = Object()
		QueryApply('scheduler_log
			where schedlog_task is ' $ Display(task) $ '
			sort schedlog_date, schedlog_company')
			{ |x|
			if x.schedlog_type is 'start'
				curCo.Add([co: x.schedlog_company, start: x.schedlog_date])

			if x.schedlog_type is 'end'
				{
				i = curCo.FindIf({ |y| x.schedlog_company is y.co })
				if false isnt first = curCo.FindOne({ |y| x.schedlog_company isnt y.co })
					{
					if i isnt false
						{
						s = curCo[i]
						overlapTime = Max(s.start, first.start)
						overlaps.Add(Object(
							time: x.schedlog_date.MinusSeconds(overlapTime)
							start: overlapTime
							end: x.schedlog_date
							co: s.co
							))
						}
					}
				curCo.Delete(i)
				}
			}
		return overlaps
		}

	GetTaskRunTime(file, task = 'all')
		{
		DeleteFile(file)
		result = Object()
		tasks = task is 'all' ? QueryList('scheduler_log', 'schedlog_task') : Object(task)
		for task in tasks.Sort!()
			{
			AddFile(file, '=='.Repeat(45/* = times*/) $ '\n')
			AddFile(file, task $ '\n')
			AddFile(file, '=='.Repeat(45/* = times*/) $ '\n')
			QueryApply('scheduler_log where schedlog_task is ' $ Display(task) $
				' sort schedlog_company, schedlog_date')
				{
				if it.schedlog_type is 'start'
					{
					if not result.Member?(it.schedlog_company)
						result[it.schedlog_company] = Object()
					if not result[it.schedlog_company].Member?(task)
						{
						AddFile(file, it.schedlog_company $ ': ' $ Display(task) $ '\n')
						result[it.schedlog_company][task] = Object()
						}
					result[it.schedlog_company][task].Add(Object(start: it.schedlog_date))
					}
				else
					{
					if not result.Member?(it.schedlog_company) or
						not result[it.schedlog_company].Member?(task)
						continue

					rec = result[it.schedlog_company][task].Last()
					rec.end = it.schedlog_date
					rec.duration = (rec.end).MinusSeconds(rec.start)
					AddFile(file, Display(rec) $ '\n')
					}
				}
			}
		}
	}
