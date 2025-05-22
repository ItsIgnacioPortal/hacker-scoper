<div align="center"><img src="icon/icon_v2_universal.png" alt ="Hacker-scoper icon"></div>
<br>
<p align="center">
  <a href="https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml"><img src="https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml/badge.svg"></a>
  <a href="https://go.dev"><img alt="Golang icon" src="https://img.shields.io/badge/Built_with-GoLang-00acd7?logo=go"></a>
  <a href="https://github.com/ItsIgnacioPortal/Hacker-Scoper/releases"><img alt="Link to the latest version" src="https://img.shields.io/github/v/release/itsignacioportal/hacker-scoper"></a>
  <a href="LICENSE.md"><img alt="Badge depicting the proyect license, the aGPLv3" src="https://img.shields.io/badge/License-aGPLv3-663366?logo=GNU"></a>
</p>


---

Hacker-scoper is a Go (v1.17.2) tool designed to assist cybersecurity professionals in bug bounty programs. It identifies and excludes URLs and IP addresses that fall outside a program's scope by comparing input targets (URLs/IPs) against a locally cached [FireBounty](https://firebounty.com) database of scraped scope data. Users may also supply a custom scope list for validation.

## 🌟 Features

- Automagically match your targets from an automatically-updated local scopes collection.
- Use your own scopes file.
- Set "explicit-level" (Parse (all as) wildcards?)
- Parse advanced wildcards as regex (supports scope filters like `amzn*.example.com` and `dev.*.example.com`).
- Match IPv4s
- Match IPv6s
- Match any valid URL ([RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html) Compliant)
- Easily chainable with other tools. Use `-ch`/`--chain-mode` to disable the fancy text decorations.
- TLD-Based detection of mis-configured bug-bounty programs: Sometimes, bug bounty programs set apk package names such as `com.my.businness.gatewayportal` as `web_application` resources instead of as `android_application` resources. hacker-scoper will detect that, and alert the user of the mis-configuration.
- If no company name and no custom files are specified, hacker-scoper will look for `.inscope` and `.noscope` files in the current or parent directories.
- Save output to a file.

## 📦 Installation

**Using Chocolatey**

```
choco install hacker-scoper
```

**Using go install**

```
go install github.com/ItsIgnacioPortal/hacker-scoper
```

**From the releases page**

Download a pre-built binary from [the releases page](https://github.com/ItsIgnacioPortal/hacker-scoper/releases)
<br>
## 🎥 Demos

### Demo with company lookup
[![asciicast](https://asciinema.org/a/WMeGitIu0VEjaFQrbv45fjhJG.svg)](https://asciinema.org/a/WMeGitIu0VEjaFQrbv45fjhJG)
<br>
<br>
<br>
<br>

### Demo with custom scopes file
[![asciicast](https://asciinema.org/a/SWtH3kLbEOmyPzrGFQe9ic9BB.svg)](https://asciinema.org/a/SWtH3kLbEOmyPzrGFQe9ic9BB)

## 🏭 Company scope matching
- **Q: How does the "company" scope matching actually work?**
- A: It works by looking for company-name matches in a cached copy of the [firebounty](https://firebounty.com/) database. The company name that you specify will be lowercase'd, and then the tool will check if any company name in the database contains that string. Once it finds a name match, it will filter your supplied targets according to the scopes that firebounty detected for that company. You can test how this would perform by just searching some name in [the firebounty website](https://firebounty.com/).

## 🤔 Usage
Usage: hacker-scoper --file /path/to/targets [--company company | --custom-inscopes-file /path/to/inscopes [--custom-outofcopes-file /path/to/outofscopes] [--verbose]] [--explicit-level INT] [--reuse Y/N] [--chain-mode] [--fire /path/to/firebounty.json] [--include-unsure] [--output /path/to/outputfile] [--hostnames-only]

### Usage examples:
- Example: Cat a file, and lookup scopes on firebounty    
  `cat recon-targets.txt | hacker-scoper -c google`

- Example: Cat a file, and use the .inscope & .noscope files    
  `cat recon-targets.txt | hacker-scoper`

- Example: Manually pick a file, lookup scopes on firebounty, and set explicit-level    
  `hacker-scoper -f recon-targets.txt -c google -e 2`

- Example: Manually pick a file, use custom scopes and out-of-scope files, and set explicit-level    
  `hacker-scoper -f recon-targets.txt -ins inscope -oos noscope.txt -e 2`

**Usage notes:** If no company and no inscope file are specified, hacker-scoper will look for ".inscope" and ".noscope" files in the current or in parent directories.

### Table of all possible arguments:
| Short | Long | Description |
|-------|------|-------------|
| -c | --company |  Specify the company name to lookup. |
| -f | --file |  Path to your file containing URLs |
| -ins | --inscope-file |  Path to a custom plaintext file containing scopes |
| -oos | --outofcope-file |  Path to a custom plaintext file containing scopes exclusions |
| -e | --explicit-level int |  How explicit we expect the scopes to be:    <br> 1 (default): Include subdomains in the scope even if there's not a wildcard in the scope    <br> 2: Include subdomains in the scope only if there's a wildcard in the scope    <br> 3: Include subdomains in the scope only if they are explicitly within the scope |
| -ch | --chain-mode |  In "chain-mode" we only output the important information. No decorations.. Default: false |
| --database |  | Custom path to the cached firebounty database |
| -iu | --include-unsure |  Include "unsure" URLs in the output. An unsure URL is a URL that's not in scope, but is also not out of scope. Very probably unrelated to the bug bounty program. |
| -o | --output |  Save the inscope urls to a file |
| -ho | --hostnames-only |  Output only hostnames instead of the full URLs |
| --verbose |  | Show what scopes were detected for a given company name. |
| --version |  | Show the installed version |
|_______________|___________________| _____________________________________ |

list example:
```javascript
example.com
dev.example.com
1.dev.example.com
2.dev.example.com
ads.example.com
192.168.1.10
192.168.2.10
192.168.2.8
```

Custom .inscope file example:
```javascript
*.example.com
*.sub.domain.example.com
amzn*.domain.example.com
192.168.1.10
192.168.2.1/24
FE80:0000:0000:0000:0202:B3FF:FE1E:8329
FE80::0202:B3FF:FE1E:8329
FE80::0204:B3FF::/24
```

Custom .noscope file example:
```javascript
community.example.com
thirdparty.example.com
*.thirdparty.example.com
dev.*.example.com
192.168.2.254
FE80::0202:B3FF:FE1E:8330
```

## :heart: Special thank you
This project was inspired by the [yeswehack_vdp_finder](https://github.com/yeswehack/yeswehack_vdp_finder)

## 📄 License
All of the code on this repository is licensed under the *GNU Affero General Public License v3*. A copy can be seen as `LICENSE` on this repository.

The library `golang.org/x/net/publicsuffix`, used within this project is licensed with [BSD-3-Clause](https://pkg.go.dev/golang.org/x/net/publicsuffix?tab=licenses).
