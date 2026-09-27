# DevDuck

Understand the error. Keep the fun.

DevDuck is a developer CLI that helps you understand errors without immediately giving away the solution. It explains what happened and gives progressive hints so you can solve the problem yourself.

## Example

```bash
duck python app.py
duck gcc main.c -o main
duck javac Main.java
duck git push
duck docker compose up
```

## How it works

```text
Command fails
    ↓
Local explanation
    ↓
Progressive hints
    ↓
Cache
    ↓
AI fallback only when needed
```

Local rules run first. AI is used only when an adapter cannot explain the diagnostic locally.

## Supported tools

- Python
- GCC / C
- javac
- Git
- Docker

## Teaching levels

- `beginner`: more terminology and up to three hints
- `intermediate`: normal explanation and up to two hints
- `advanced`: concise output and one hint

## Installation

Release binaries will be provided when the first release is published. Until then, install or build from source:

```bash
go install github.com/duck-ahiru-Z/DevDuck/cmd/duck@latest
go build -o duck ./cmd/duck
```

## Usage

```bash
duck --level beginner python app.py
duck config show
duck config set level beginner
duck auth set gemini
duck auth status
duck auth delete gemini
duck version
duck doctor
```

`duck auth set gemini` stores the key in the operating system credential store. `GEMINI_API_KEY` remains available as a development and CI fallback.

## AI

DevDuck uses local explanations first. Gemini is an opt-in fallback for diagnostics that have no local rule. Configure it with the OS credential store or `GEMINI_API_KEY`.

## Security and privacy

- API keys are stored in OS credential storage.
- Secrets are redacted before AI requests.
- API keys are not stored in config files or the explanation cache.
- AI is called only when local explanations cannot handle the error.

## Development

```bash
go test ./...
go vet ./...
```

The core runs commands through a common runner and selects an adapter from the registry. Adapters parse diagnostics and provide local teaching rules. The Teaching Engine controls hint depth; Cache, Redactor, Credential Store, and AI Provider remain separate services.

## License

This repository does not currently include a `LICENSE` file. The project license is still to be chosen by the maintainers.
