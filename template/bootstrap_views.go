package template

var BootstrapLayoutTmpl = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <link href="/static/css/site.css" rel="stylesheet"/>
    {{partial_r ".head"}}
</head>
<body>
    <div class="header">
        <ul>
            <li><a href="/">Index</a></li>
            <li><a href="/about">About</a></li>
        </ul>
    </div>
    <hr>
    {{yield}}
    <hr>
    <div>
        Copyright @teamlint
    </div>
</body>
</html>`

var BootstrapErrorTmpl = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
<body>
    错误
    <h2>Title:{{.Error.Title}}</h2>
    <h3>StatusCode:{{.Error.StatusCode}}</h3>
    <h3>Code:{{.Error.Code}}</h3>
    <p>Message:{{.Error.Message}}</p>
</body>
</html>
`
var BootstrapIndexTmpl = `<h2>{{.Title}}</h2><p>{{.Message}}</p>`
var BootstrapIndexHeadTmpl = `<style>
.header{
    background-color: #e9e9e9;
}
</style>`
var BootstrapAboutTmpl = `<h2>{{.Title}}</h2><p>{{.Message}}</p><p>Session["foo"]:{{.SessFoo}}</p>`
