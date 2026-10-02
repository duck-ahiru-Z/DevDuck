# DevDuck hackathon demo

This sequence is designed for a short presentation. It uses only local files and
a temporary Git repository; it does not contact a network or change global Git
configuration.

## 0. Build DevDuck

Run this from the repository root. The binary stays in a temporary directory.

```powershell
$duckDemoBin = Join-Path $env:TEMP "devduck-demo"
New-Item -ItemType Directory -Force $duckDemoBin | Out-Null
go build -o (Join-Path $duckDemoBin "duck.exe") ./cmd/duck
$env:Path = "$duckDemoBin;$env:Path"
function prompt { "PS demo> " }
```

Use `cls` before each screenshot to keep the frame clean.

## 1. Ordinary Python error

Use whichever Python command is installed (`python`, `python3`, or `py`):

```powershell
python demo/python/zero_division.py
```

Capture this as the ordinary error. The program makes the cause visible: there
are 120 points but zero students, so an average cannot be calculated.

## 2. DevDuck beginner explanation

```powershell
duck --level beginner python demo/python/zero_division.py
```

The frame should show the original traceback, the local explanation, and Hint 1.

## 3. Progressive Hint 2

At the prompt, type `h` and press Enter. Capture the frame with Hint 2. Type `q`
when you are done.

## 4. Fixed version

```powershell
duck --level beginner python demo/python/fixed.py
```

This succeeds and shows only the normal result (`30.0`), with no DevDuck error
message.

## 5. NameError

```powershell
duck --level beginner python demo/python/name_error.py
```

This is a compact second example: the assigned name is `project_name`, but the
program tries to print `project`.

## 6. Git Adapter

Run the safe local scenario from the repository root:

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

The first command demonstrates “not a Git repository”; the second DevDuck
command demonstrates a missing branch. Both are offline and isolated in `$env:TEMP`.

## Suggested screenshots

1. The ordinary Python traceback.
2. DevDuck beginner output with Explanation and Hint 1.
3. The next screen after entering `h`, showing Hint 2.
4. The successful `fixed.py` output (`30.0`).

## Suggested 30-second video order

Clear the screen, run the ordinary error, run the same file through DevDuck,
press `h`, then run `fixed.py`. If there is time, finish with `duck git status`
in the temporary Git directory to show that the same teaching flow is not
Python-only.
