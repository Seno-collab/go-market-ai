# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported |
| ------- | --------- |
| `main`  | ✅        |
| Latest release (vX.Y.Z) | ✅ |
| Older releases | ❌ |

> If you are using an unsupported version, please upgrade to a supported release.

---

## Reporting a Vulnerability

Please report security issues **privately**. Do **not** open a public GitHub Issue for vulnerabilities.

**Preferred reporting channels:**
1. GitHub Security Advisories (Private reporting):  
   Go to **Security → Advisories → Report a vulnerability**.
2. Email: `security@YOUR_DOMAIN` (replace with your email)

When reporting, please include:
- A clear description of the issue and potential impact
- Steps to reproduce (PoC is helpful, but keep it minimal)
- Affected versions / commit hash
- Any relevant logs, configs, or environment details
- Your suggested fix (if you have one)

We will acknowledge receipt within **72 hours**.

---

## Disclosure Policy

We follow a coordinated disclosure process:
- We will confirm the vulnerability and assess severity.
- We will work on a fix and prepare a release.
- We will coordinate a public disclosure date with the reporter when possible.

We aim to provide a fix or mitigation within:
- **7 days** for critical issues (where feasible)
- **30 days** for high/medium issues (depending on complexity)

Timelines may vary based on impact, exploitability, and available information.

---

## Security Updates

Security fixes are released as:
- Patch releases (e.g., `vX.Y.(Z+1)`) whenever possible
- Backports only for supported versions

Release notes will mention security fixes in a non-exploitable way until most users have had time to update.

---

## Dependency & Go Module Security

We take dependency security seriously:
- Dependencies are tracked via `go.mod` / `go.sum`
- We regularly scan for known vulnerabilities

Recommended checks for maintainers:
- Run Go vulnerability scanning:
  ```bash
  govulncheck ./...
