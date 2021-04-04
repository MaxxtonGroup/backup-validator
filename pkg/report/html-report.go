package report

var HTML_REPORT = `
<!DOCTYPE html>
<html>
  <title>Backup Validation Test Report - {{ .Time }}</title>
  <style>
    a:hover {
      text-decoration: underline !important;
    }
  </style>
</html>
<body  style="font-family: 'Arial', sans-serif; background-color: #f6f6f7; font-size: 16px">
  <div style="max-width: 600px; margin-left: auto; margin-right: auto; padding: 20px; padding-bottom: 30px; background-color: white;">
    <h1 style="font-size: 20px; margin-bottom: 40px;">Backup Validation Test Report - {{ .Time }}</h1>
    <table style="width: 100%; border-collapse: collapse;">
      <thead>
        <tr style="background-color: #f6f6f7">
          <th style="text-align: left; padding: 10px; border: 1px solid #f6f6f7;">Test</th>
          <th style="text-align: left; padding: 10px; border: 1px solid #f6f6f7;">Result</th>
          <th style="text-align: left; padding: 10px; border: 1px solid #f6f6f7;">Duration</th>
        </tr>
      </thead>
      <tbody>
        {{- range .TestResults }}
        <tr>
          <td style="text-align: left; padding: 10px; border: 1px solid #f6f6f7; min-width: 150px">{{ .Name }}</td>
          {{- if .Error }}
          <td style="text-align: left; padding: 10px; border: 1px solid #f6f6f7; background-color: #f44336; color: white" class="passed">
            Error
            <ul style="font-size: 12px; margin: 0; padding-left: 20px;">
              <li>{{ .Error }}</li>
            </ul>
          </td>
          {{- else if and .FailedAsserts (gt (len .FailedAsserts) 0) }}
          <td style="text-align: left; padding: 10px; border: 1px solid #f6f6f7; background-color: #f44336; color: white" class="passed">
            Failed
            <ul style="font-size: 12px; margin: 0; padding-left: 20px;">
              {{- range .FailedAsserts }}
              <li>{{ . }}</li>
              {{- end }}
            </ul>
          </td>
          {{- else }}
          <td style="text-align: left; padding: 10px; border: 1px solid #f6f6f7; background-color: #4caf50; color: white" class="passed">Passed</td>
          {{- end }}
          <td style="text-align: left; padding: 10px; border: 1px solid #f6f6f7; font-size: 12px; min-width: 150px;">total: {{ .TotalDuration }}<br>(restore: {{ .RestoreDuration }}, import: {{ .ImportDuration }})</td>
        </tr>
        {{- end }}
      </tbody>
    </table>
  </div>
  <div style="text-align: center; margin-top: 20px;">
    <a href="https://github.com/MaxxtonGroup/backup-validator">Backup Validator</a>
  </div>
</body>
`
