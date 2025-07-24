// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
ContinuousTest_Base
	{
	CallClass(libs = false)
		{
		if libs is false
			{
			libs = GetContributions('ApplicationLibraries')
			libs.RemoveIf({ not QueryColumns(it).Has?('group') })
			}
		return .runQc(libs)
		}

	ExceptionFile: 'exceptions.txt'
	exceptionSplit: '\n---------\n'
	runQc(libs)
		{
		asOf = Date()
		allStatistics = Object()
		hasExceptions? = false
		PutFile(.ExceptionFile, '')
		for lib in libs
			{
			statistics = Object(recsInLib: 0, worst: Object(), libAvgRating: 0,
				totalLibWarnings: 0, :lib)
			for .. .worstSize
				statistics.worst.Add(#(name: '', rating: 5))
			QueryApply(lib, group: -1)
				{
				//QUESTION: Ignore tests?
				try
					{
					warnings = Qc_Main.CheckWithExtra(
						lib, it.name, it.lib_current_text, minimizeOutput?:)
					if false isnt rating = Qc_Main.CalcRatings(warnings)
						{
						rating /= 2 // convert to 5 star format
						if rating < 5 /*= Max rating*/
							.saveWorst(rating, statistics.worst, it.name)
						statistics.totalLibWarnings += .countWarnings(warnings)
						statistics.libAvgRating += rating
						statistics.recsInLib++
						}
					}
				catch(e)
					{
					hasExceptions? = true
					AddFile(.ExceptionFile, lib $ ':' $ it.name $ '\t' $ e $ '\n' $
						FormatCallStack(e.Callstack(), levels: 20) $ .exceptionSplit)
					}
				}
			statistics.libAvgRating /= statistics.recsInLib
			statistics.libAvgRating = statistics.libAvgRating.Round(2)
			statistics.worst.RemoveIf({ it.name is '' }).Reverse!()
			allStatistics.Add(statistics)
			}
		exceptmsg = ''
		if hasExceptions?
			exceptmsg = .exceptionsMessage()
		return .buildHtml(allStatistics, exceptmsg, asOf)
		}

	countWarnings(warnings)
		{
		return warnings.Values().Map({
			it.GetDefault('size', it.GetDefault('warnings', #()).Size())
			}).Sum()
		}

	worstSize: 10
	saveWorst(rating, worst, name) // bigger value in the front
		{
		if rating > worst.First().rating
			return

		if rating < worst.Last().rating
			{
			worst.Add([:name, :rating])
			worst.Delete(0)
			return
			}

		for (i = 0; i < worst.Size(); ++i)
			{
			if rating >= worst[i].rating
				{
				worst.Add([:name, :rating], at: i)
				if worst.Size() > .worstSize
					worst.Delete(0)
				break
				}
			}
		}

	tableTemplate:
		`<table id="myTable">
			<thead><tr>
				<th onclick="sortTable(0)" id='col0'>Library</th>
				<th onclick="sortTable(1)" id='col1'>Number of Records</th>
				<th onclick="sortTable(2)" id='col2'>Average Rating</th>
				<th onclick="sortTable(3)" id='col3'>Total Warnings</th>
				<th onclick="sortTable(4)" id='col4'>10 Worst in Lib</th>
			</tr></thead>
			<tbody>
			@for(lib in .allStatistics)
				{
				<tr>
					<td>@lib.lib</td>
					<td>@lib.recsInLib</td>
					<td>@lib.libAvgRating</td>
					<td>@lib.totalLibWarnings</td>
					<td>
						@for(worstRec in lib.worst)
							{
							<div>@worstRec.name: @worstRec.rating</div>
							}
					</td>
				</tr>
				}
			</tbody>
		</table>`
	buildHtml(allStatistics, exceptmsg, asOf)
		{
		return Razor(`<html>
			<head>
			<link rel="stylesheet" href="./styles.css" type="text/css" />
			<title>Quality Checking Result as of @.asOf.ShortDateTime()</title>` $
				.sortScript() $
			`</head>
			<body>
			<h1>Quality Checking Result</h1>
			<h2>as of @.asOf.ShortDateTime()</h2>` $
			(exceptmsg.Blank?() ? `` : exceptmsg) $
			.tableTemplate $
			`</body></html>`, [:allStatistics, :asOf])
		}

	sortScript()
		{
		return `
			<script>
				var dir;
				var sortCol;
				function setSortCol(n) {
					sortCol = n;
				}
				function getSortCol() {
					return sortCol;
				}
				function sortRowAsc(a, b) {
					var col = getSortCol();
					if ((col !== 0 && col !== 4))
						return compareNum(a, b, col);
					return compareStr(a, b, col);
				}
				function sortRowDes(a, b) {
					var col = getSortCol();
					if ((col !== 0 && col !== 4))
						return -compareNum(a, b, col);
					return -compareStr(a, b, col);
				}
				function compareNum(a, b, col) {
					return Number(getInnerHTML(a, col)) - Number(getInnerHTML(b, col));
				}
				function compareStr(a, b, col) {
					return (getInnerHTML(a, col).toLowerCase() >=
						getInnerHTML(b, col).toLowerCase())
						? 1
						: -1;
				}
				function getInnerHTML(row, col) {
					return row.getElementsByTagName("TD")[col].innerHTML;
				}` $
				.sortTable() $
			`</script>`
		}

	sortTable()
		{
		return `
				function sortTable(n) {
					const table = document.getElementById("myTable");
					const rows = table.getElementsByTagName("TR");
					const sortHdr = rows[0].getElementsByTagName("TH")[n];
					const arr = document.getElementById('sortArr');
					var items = [].slice.call(rows);
					var parent = items[1].parentNode;
					delete items[0];

					setSortCol(n);
					if (dir == 'des' || dir == undefined) {
						items.sort(sortRowAsc);
						dir = 'asc';
					}
					else {
						items.sort(sortRowDes);
						dir = 'des';
					}

					for (let i = 0; i < items.length - 1; i++) {
						let removedItem = parent.removeChild(items[i]);
						parent.appendChild(removedItem);
					}

					if (arr != null)
						arr.remove();
					const img = document.createElement('img');
					img.setAttribute('id', 'sortArr');
					img.setAttribute('src', (dir == 'asc') ? './up.png' : './down.png');
					sortHdr.insertAdjacentElement("beforeend", img);
				}
				`
		}

	exceptionsMessage()
		{
		str = `<h1>QC Stats FAILED</h1>`
		count = 0
		File(.ExceptionFile)
			{ |file|
			while false isnt line = file.Readline()
				{
				if line is .exceptionSplit.Trim()
					{
					++count
					str $= `<br>` $ .exceptionSplit $ `<br>`
					continue
					}
				if count >= 5 /*=max errors reported*/
					{
					str $= `<br>More errors in file...<br><br>`
					break
					}
				str $= line $ `<br>`
				}
			}
		return str
		}
	}
