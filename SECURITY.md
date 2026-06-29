# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |

## Reporting a Vulnerability

Dbasement is an offline tool that runs locally on your machine. It does not
communicate with any external services, send telemetry, or make network
requests. The risk of external exploitation is minimal.

However, if you discover a security vulnerability:

1. **Do not** open a public GitHub issue.
2. Email the maintainers directly or use GitHub's private vulnerability
   reporting feature.
3. Include a detailed description of the vulnerability and steps to reproduce.

You should receive a response within 48 hours. If the vulnerability is
confirmed, a patch will be released as soon as possible.

## Scope

The following are considered in scope for security reports:

- SQL injection in the memory database
- Path traversal in project path handling
- Remote code execution via malicious MCP messages
- Information disclosure through error messages
- Insecure default configurations

## Out of Scope

- Features explicitly marked as experimental
- Dependencies with known CVEs (please report those upstream)
- Social engineering attacks

## Safe Harbor

We consider security research conducted responsibly to be protected. We will
not take legal action against anyone who:

- Makes a good faith effort to avoid privacy violations
- Does not cause harm or data destruction
- Reports the issue privately first

## Recommendations

When using Dbasement in your projects:

1. Always run the latest version
2. Review the contents of `.dbasement/memory.db` before sharing your repository
   publicly, as it may contain project-specific information
3. Add `.dbasement/` to your `.gitignore` to prevent committing memory files
   (the project memory is designed to be local to each developer)
