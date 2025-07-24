// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
`.su-listbody tr {
	color: black;
}
.su-listbody.su-list-dragging {
	cursor: move;
}
.su-listbody-hover tr:hover,
.su-listbody tr.su-listbody-hovered {
	color: #c0c0c0; /*CLR.silver*/
}
.su-listbody td {
	position: relative;
	border: none;
	padding: 0;
}
.su-listbody-cell {
	overflow: hidden;
	white-space: nowrap;
	padding-left: 4px;
	padding-right: 4px;
	height: 100%;
}
.su-listbody-cell-image {
	height: 1em;
}
.su-listbody-cell-rect {
	border-style: solid;
	box-sizing: border-box;
	height: calc(100% - 2px);
	width: calc(100% - 2px);
	margin: 1px;
}
.su-listbody-cell-circle {
	border-style: solid;
	border-radius: 50%;
	box-sizing: border-box;
	width: calc(var(--su-row-height) * 0.8);
	height: calc(var(--su-row-height) * 0.8);
	margin-left: calc(var(--su-row-height) * 0.5);
	margin-top: calc(var(--su-row-height) * 0.1);
}
.su-listbody-cell-part {
	pointer-events: none;
	overflow: hidden;
	padding: 0px;
	margin: 0px;
	border: 0px;
}
.su-listbody td[data-type="mark-cell"] {
	background-color: var(--su-color-buttonface);
	font-family: suneido;
	font-style: normal;
	font-weight: normal;
	color: red;
	text-align: center;
}
.su-listbody td[data-type="fill-cell"] {
	height: 2em;
}
.su-shadinglist tr[data-type="data-row"]:nth-child(odd) {
	background-color: azure;
}
.su-listbody-cell:hover {
	outline: 1px solid grey;
	outline-offset: -1px;
}
`
