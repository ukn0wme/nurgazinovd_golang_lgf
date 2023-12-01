package templates

{{define "subject"}}Welcome to Salemmusic!{{end}}

{{define "plainBody"}}
Hi,

Thanks for signing up for a SalemMusic account. We're excited to have you on board!

For future reference, your user ID number is {{.ID}}.

Thanks,

The SalemMusic Team
{{end}}
{{define "htmlBody"}}
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width" />
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head>
<body>
<p>Hi,</p>
<p>Thanks for signing up for a Greenlight account. We're excited to have you on board!</p>
<p>For future reference, your user ID number is {{.ID}}.</p>
<p>Thanks,</p>
<p>The Greenlight Team</p>
</body>
</html>
{{end}}

