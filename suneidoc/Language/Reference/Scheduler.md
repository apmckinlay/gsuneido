### Scheduler

``` suneido
(source = 'schedtasks', sourceName = 'scheduler')
```

Run functions at scheduled times.

source is either a table name (which will be passed to SchedTable, or a function that returns a list of tasks with: (sched_when:, sched_func:, sched_args:)

For example:

``` suneido
tasks = function ()
    {
    return #((sched_when: 'every 1 minute', 
        sched_func: 'Print', sched_args: "'hello'"))
    }
Scheduler(tasks)
```

The Scheduler sleeps for minute between checking if tasks are due.

The tasks function is re-run each time, so the tasks can be updated while the Scheduler is running.

Currently supports:

-	SchedEvery e.g. "every 5 minutes"
-	SchedAt e.g.  "at 19:00" or "at 8:00 skip weekends"
-	SchedOn e.g. "on Wed at 15:00"   
	Adds day of week to SchedAt
-	SchedMonthlyOn e.g. "on EndMonth at 23:00" or "on 5 at 7:00"   
	Adds day of month to SchedAt


To stop a Scheduler e.g. if it's running in a [Thread](<Thread.md>), call SchedExit()