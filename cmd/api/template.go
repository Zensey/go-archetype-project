package main

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
		<style>

        table.darkTable {
          font-family: "Arial Black", Gadget, sans-serif;
          border: 2px solid #000000;
          background-color: #4A4A4A;
          width: 600px;
          
          text-align: center;
          border-collapse: collapse;
        }
        table.darkTable td, table.darkTable th {
          border: 1px solid #4A4A4A;
          padding: 3px 2px;
        }
        table.darkTable tbody td {
          font-size: 13px;
          color: #E6E6E6;
        }
        table.darkTable tr:nth-child(even) {
          background: #888888;
        }
        table.darkTable thead {
          background: #000000;
          border-bottom: 3px solid #000000;
        }
        table.darkTable thead th {
          font-size: 15px;
          font-weight: bold;
          color: #E6E6E6;
          text-align: center;
          border-left: 2px solid #4A4A4A;
        }
        table.darkTable thead th:first-child {
          border-left: none;
        }

        table.darkTable tfoot {
          font-size: 12px;
          font-weight: bold;
          color: #E6E6E6;
          background: #000000;
          background: -moz-linear-gradient(top, #404040 0%, #191919 66%, #000000 100%);
          background: -webkit-linear-gradient(top, #404040 0%, #191919 66%, #000000 100%);
          background: linear-gradient(to bottom, #404040 0%, #191919 66%, #000000 100%);
          border-top: 1px solid #4A4A4A;
        }
        table.darkTable tfoot td {
          font-size: 12px;
        }
	</style>
	</head>
	<body>

  <table class="darkTable">
      <thead>
      <tr>
      <th>Approved</th>
      <th>Total</th>
      </tr>
      </thead>
      <tbody>
      {{range .Count}}
          <tr>
          <td>{{if .Status}} {{.Status}} {{else}} ? {{end}} </td>
          <td>{{ .Count }}</td>
          </tr>
      {{end}}
      </tbody>
  </table>
  </br>

  <table class="darkTable">
      <thead>
      <tr>
      <th>Name</th>
      <th>ProdID</th>
      <th>Approved</th>
      <th>Created</th>
      </tr>
      </thead>
      <tbody>
      {{range .Items}}
          <tr>
          <td>{{ .Name }}</td>
      	<td>{{ .ProductID }}</td>
          <td>{{if .Status}} {{.Status}} {{else}} ? {{end}} </td>
      	<td>{{ .Created }}</td>
          </tr>
      {{end}}
      </tbody>
  </table>

</br>
</br>


	</body>
</html>`
