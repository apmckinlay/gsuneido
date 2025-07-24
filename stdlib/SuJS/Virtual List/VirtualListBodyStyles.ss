// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
`.su-vlistbody tr {
	color: black;
	position: relative;
}
.su-vlistbody td {
	position: relative;
	border: none;
	padding: 0;
}
.su-vlistbody tr[data-type="data-row"].su-vlistbody-row-selected {
	background-color: #0076d7; /*=COLOR.HIGHLIGHT*/
	color: white;
}
.su-vlistbody tr[data-type="data-row"].su-vlistbody-row-highlighted:not(.su-vlistbody-row-selected) {
	background-color: var(--su-vlistbody-row-color);
}
.su-vlistbody-cell {
	overflow: hidden;
	white-space: nowrap;
	padding-left: 4px;
	padding-right: 4px;
	height: 100%;
}
.su-vlistbody-cell:hover {
	outline: 1px solid grey;
	outline-offset: -1px;
}
.su-vlistbody td[data-type="mark-cell"],
.su-listhead th[data-type="mark-cell"],
.su-vlistbody td[data-type="expand-mark-cell"] {
	background-color: var(--su-color-buttonface);
	font-family: suneido;
	font-style: normal;
	font-weight: normal;
	color: red;
	text-align: center;
	left: 0px;
	position: sticky;
	z-index: 2;
	width: 100%;
}
.su-vlistbody tr[data-type="expanded-row"] {
	background-color: var(--su-color-buttonface);
}
.su-vlistbody .su-vlist-edit-button {
	color: grey;
	position: absolute;
	top: 2px;
	left: 0px;
	right: 0px;
}
.su-vlistbody .su-vlist-edit-button:hover {
	color: black;
	border: deepskyblue 1px solid;
}
.su-vlistbody tr[data-editing] .su-vlist-edit-button {
	color: black;
	border: deepskyblue 1px solid;
	background-color: azure;
}
.su-vlistbody td[data-type="fill-cell"] {
	height: 2em;
}
.su-vshadinglist tr[data-type="data-row"]:nth-child(odd):not(.su-vlistbody-row-selected):not(.su-vlistbody-row-highlighted) {
	background-color: azure;
}
.su-vlist-expand-buttons {
	font-family: suneido;
	font-style: normal;
	font-weight: normal;
	background-color: var(--su-color-buttonface);
	position: absolute;
	left: 0px;
	top: 0px;
	z-index: 1;
}
.su-vlist-expand-button {
	color: grey;
	margin-left: 0.2em;
	margin-right: 0.2em;
}
.su-vlist-expand-button:hover {
	color: black;
	outline: deepskyblue solid 1px;
}
.su-vlist-dragging	{
	cursor: move;
}
`
