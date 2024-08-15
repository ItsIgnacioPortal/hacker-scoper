![Hacker-scoper icon](icon/icon.png)

<h1 align="center">Hacker Scoper</h1>

<p align="center">
  <a href="https://hits.seeyoufarm.com">
    <img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2FItsIgnacioPortal%2Fhacker-scoper&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false">
  </a>
  <a href="https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml">
    <img src="https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml/badge.svg">
  </a>
</p>

<p align="center">
  <img src="https://forthebadge.com/images/badges/made-with-go.svg">
  <img src="https://forthebadge.com/images/badges/built-with-love.svg">
  <img src="https://forthebadge.com/images/badges/cc-sa.svg">
</p>

---

This is a Go v1.17.2 application made for quickly filtering out URLs and IP addresses which are outside of our scope. Designed with bug-bounty programs in mind, the tool will match your given `targets` URLs/IPv4s/IPv6s with those from a locally stored copy of the full [firebounty](https://firebounty.com) json of scraped scopes, OR with your own list of scopes!

## Features

- Automagically match your targets from an automatically-updated local scopes collection.
- Use your own scopes file.
- Set "explicit-level" (Parse (all as) wildcards?)
- Parse advanced wildcards as regex (supports scope filters like `amzn*.example.com` and `dev.*.example.com`).
- Match IPv4s
- Match IPv6s
- Match any valid URL ([RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html) Compliant)
- Attempt to scrape security.txt files from your targets
- Easily chainable with other tools. Use `-ch`/`--chain-mode` to disable the fancy text decorations.
- TLD-Based detection of mis-configured bug-bounty programs: Sometimes, bug bounty programs set apk package names such as `com.my.businness.gatewayportal` as `web_application` resources instead of as `android_application` resources. hacker-scoper will detect that, and alert the user of the mis-configuration.
- If no company name and no custom files are specified, hacker-scoper will look for `.inscope` and `.noscope` files in the current or parent directories.
- Save output to a file.

## Installation

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

## Usage

### Demo with company lookup 
[![asciicast](https://asciinema.org/a/WMeGitIu0VEjaFQrbv45fjhJG.svg)](https://asciinema.org/a/WMeGitIu0VEjaFQrbv45fjhJG)
<br>
<br>
<br>
<br>

### Demo with custom scopes file
[![asciicast](https://asciinema.org/a/SWtH3kLbEOmyPzrGFQe9ic9BB.svg)](https://asciinema.org/a/SWtH3kLbEOmyPzrGFQe9ic9BB)

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
| -cstxt | --check-security-txt |  Whether or not we will try to scrape security.txt from all domains and subdomains (Warning: experimental feature.) |
| -r | --reuse |  Reuse previously generated security.txt lists? (Y/N). Only needed if using "-cstxt" |
| -f | --file |  Path to your file containing URLs |
| -ins | --inscope-file |  Path to a custom plaintext file containing scopes |
| -oos | --outofcope-file |  Path to a custom plaintext file containing scopes exclusions |
| -e | --explicit-level int |  How explicit we expect the scopes to be:    <br> 1 (default): Include subdomains in the scope even if there's not a wildcard in the scope    <br> 2: Include subdomains in the scope only if there's a wildcard in the scope    <br> 3: Include subdomains in the scope only if they are explicitly within the scope |
| -ch | --chain-mode |  In "chain-mode" we only output the important information. No decorations.. Default: false |
| --fire |  | Set this to specify a path for the FireBounty JSON. |
| -iu | --include-unsure |  Include "unsure" URLs in the output. An unsure URL is a URL that's not in scope, but is also not out of scope. Very probably unrelated to the bug bounty program. |
| -o | --output |  Save the inscope urls to a file |
| -ho | --hostnames-only |  Output only hostnames instead of the full URLs |
| --verbose |  | Show what scopes were detected for a given company name. |
| --version |  | Show the installed version |

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

## F.A.Q
- **Q: How does the "company" scope matching actually work?**     
A: It works by looking for company-name matches in a local copy of the [firebounty](https://firebounty.com/) database (`firebounty-scope-url_only.json`). The company name that you specify will be lowercase'd, and hacker-scoper will check if any company name in the JSON contains that string. Once it finds a name match, it'll filter your URLs according to the scopes that firebounty detected for that company. You can test how this would perform by just searching some name in [the firebounty website](https://firebounty.com/).

## Special thank you
This project was inspired by the [yeswehack_vdp_finder](https://github.com/yeswehack/yeswehack_vdp_finder)

## License
All of the code on this repository is licensed under the *Creative Commons Attribution-ShareAlike License*. A copy can be seen as `LICENSE` on this repository.
