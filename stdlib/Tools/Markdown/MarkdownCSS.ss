/* === Embedded CSS for Markdown HTML === */
`:root {
	--page-max-width: 820px;
	--page-padding: 28px;
	--body-bg: #ffffff;
	--body-fg: #1b1f23;
	--muted: #6b7280;
	--link: #0b6cff;
	--border: #e6e9ee;
	--code-bg: #f6f8fa;
	--pre-bg: #0f1720;
	--accent: #ff7b72;
	--quote-bg: #f8fafc;
	--radius: 10px;
	--mono: ui-monospace, SFMono-Regular, Menlo, Monaco, \"Roboto Mono\", \"Courier New\", monospace;
	--ui-font: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial;
	--base-font-size: 16px;
	--line-height: 1.7;
	--shadow: 0 6px 18px rgba(18, 22, 30, 0.06);
	--toc-width: 260px;
}

@media (prefers-color-scheme: dark) {
	:root {
		--body-bg: #0b0f13;
		--body-fg: #e6eef6;
		--muted: #94a3b8;
		--link: #7cc2ff;
		--border: #1f2937;
		--code-bg: #071022;
		--pre-bg: #071022;
		--accent: #ff9e94;
		--quote-bg: rgba(255, 255, 255, 0.03);
		--shadow: 0 8px 30px rgba(2, 6, 23, 0.6);
	}
}

html,
body {
	height: 100%;
	margin: 0;
	padding: 0;
	background: var(--body-bg);
	color: var(--body-fg);
	font-family: var(--ui-font);
	font-size: var(--base-font-size);
	line-height: var(--line-height);
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
}

.container {
	max-width: var(--page-max-width);
	margin: 36px auto;
	padding: var(--page-padding);
}

@media (max-width:900px) {
	.container {
		margin: 18px;
		padding: 18px;
	}
}

h1,
h2,
h3,
h4,
h5,
h6 {
	margin: 1.2em 0 0.35em;
	line-height: 1.15;
	font-weight: 600;
}

h1 {
	font-size: 2.0rem;
}

h2 {
	font-size: 1.45rem;
}

h3 {
	font-size: 1.15rem;
}

h4 {
	font-size: 1rem;
	color: var(--muted);
}

h5,
h6 {
	font-size: 0.95rem;
	color: var(--muted);
}

p {
	margin: 0 0 1em;
}

ol,
ul {
	margin: 0 0 1em 1.4em;
	padding: 0;
}

li {
	margin: 0.35em 0;
}

a {
	color: var(--link);
	text-decoration: none;
	border-bottom: 1px dotted rgba(0, 0, 0, 0.06);
}

a:hover {
	text-decoration: underline;
}

code {
	font-family: var(--mono);
	background: var(--code-bg);
	padding: 0.12em 0.38em;
	border-radius: 6px;
	font-size: 0.95em;
}

pre code {
	font-family: var(--mono);
	font-size: 0.92rem;
	margin: 1em 0;
	overflow: auto;
	border-radius: 10px;
	box-shadow: var(--shadow);
	padding: 14px;
}

pre code {
	display: block;
	white-space: pre;
	line-height: 1.5;
}

blockquote {
	margin: 1em 0;
	padding: 14px 18px;
	border-radius: 8px;
	background: var(--quote-bg);
	border-left: 3px solid var(--accent);
	color: var(--muted);
}

hr {
	height: 1px;
	border: 0;
	background: linear-gradient(90deg, transparent, var(--border), transparent);
	margin: 1.4em 0;
}

table {
	width: 100%;
	border-collapse: collapse;
	margin: 0.9em 0 1.2em;
	display: block;
	overflow: auto;
}

th,
td {
	padding: 10px 12px;
	border: 1px solid var(--border);
	text-align: left;
}

th {
	background: linear-gradient(180deg, rgba(0, 0, 0, 0.02), transparent);
	font-weight: 600;
}

img {
	max-width: 100%;
	height: auto;
	vertical-align: text-bottom;
	border-radius: 8px;
	box-shadow: 0 6px 18px rgba(2, 6, 23, 0.04);
}

figcaption {
	font-size: 0.9rem;
	color: var(--muted);
	margin-top: 6px;
}`
