# Safe Git demo

These commands use a temporary directory and do not contact a remote. They do not
change global Git configuration.

```powershell
$gitDemo = Join-Path $env:TEMP "devduck-git-demo"
Remove-Item -Recurse -Force $gitDemo -ErrorAction SilentlyContinue
New-Item -ItemType Directory $gitDemo | Out-Null
Push-Location $gitDemo
duck git status
git init
duck git checkout definitely-not-existing-branch
Pop-Location
```
