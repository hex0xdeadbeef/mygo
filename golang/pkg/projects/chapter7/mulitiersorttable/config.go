package mulitiersorttable

import (
	htmltempl "html/template"
)

const (
	fileName = "table.html"
)

var tableTemplate = htmltempl.Must(
	htmltempl.
		New("table").Parse(`
<!DOCTYPE>
<html>

<head>
	<title>Table Widget</title>
	<style>
		h1 {
			text-align: center;
		}
		table {
			border-collapse: collapse;
			width: 100%;
		}
		th, td {
			border: 1px solid black;
			padding: 8px;
			text-align: left;
		}
		th {
			background-color: #dcfff2;
			cursor: pointer;
		}
	</style>
	<script>	
		function getColumnIndex(columnIndex) {
			window.location.href = "/sort?columnIndex=" + columnIndex;
		}
	</script>
</head>

<body>
	<h1>Info</h1>
	<table>
			<thead>
				<tr>
					{{range $index, $header := .TableHeaders}}
						<th onclick="getColumnIndex({{$index}})">{{$header.Title}}</th>
					{{end}}
				</tr>
			</thead>
			
			<tbody>
				{{range .TableRows}}
					<tr>
						{{range .}}
							<td>{{.}}</td>
						{{end}}
					</tr>
				{{end}}
			</tbody>		
	</table>
</body>

</html>
`))

func getColumnIndex(columnIndex int) int {
	return columnIndex
}
