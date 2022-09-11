# Hacker-Scoper
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2FItsIgnacioPortal%2Fhacker-scoper&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)
[![goreleaser](https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml/badge.svg)](https://github.com/ItsIgnacioPortal/hacker-scoper/actions/workflows/gorelease.yml)  
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) 
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/cc-sa.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/check-it-out.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/fo-real.svg)](https://forthebadge.com)

This is a go1.17.2 application made for quickly filtering out URLs and IP addresses which are outside of our scope. Designed with bug-bounty programs in mind, the tool will match your given `targets` URLs/IPv4s/IPv6s with those from a locally stored copy of the full [firebounty](https://firebounty.com) json of scraped scopes, OR with your own list of scopes!

## Features

- Automagically match your targets from an automatically-updated local scopes collection.
- Use your own scopes file
- Set "explicit-level" (Parse (all as) wildcards?)
- Match IPv4s
- Match IPv6s
- Match any valid URL ([RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html) Compliant)
- Attempt to scrape security.txt files from your targets
- 100% chainable with other tools: Just use `--chain-mode`, and begin piping targets into STDIN!
- Basic detection of mis-configured bug-bounty programs: Detect if an APK package name was set as a domain (`com.android.example` for example)
- If no company name and no custom files are specified, look for `.inscope` and `.noscope` files in the current or parent directories.

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
[![asciicast](https://asciinema.org/a/fKXTmmbtNPqKQmn3CrAaXNTB3.svg)](https://asciinema.org/a/fKXTmmbtNPqKQmn3CrAaXNTB3)
<br>
<br>
<br>
<br>

### Demo with custom scopes file
[![asciicast](https://asciinema.org/a/N8hohdAxM9cM0RkC6ptHFJygE.svg)](https://asciinema.org/a/N8hohdAxM9cM0RkC6ptHFJygE)

```
Usage: hacker-scoper --file /path/to/targets [--company company | --custom-inscopes-file /path/to/inscopes [--custom-outofcopes-file] /path/to/outofscopes] [--explicit-level INT] [--reuse Y/N] [--chain-mode] [--fire /path/to/firebounty.json]

Usage examples:
  Example: Cat a file, and lookup scopes on firebounty
  cat recon-targets.txt | hacker-scoper -c google

  Example: Cat a file, lookup scopes on firebounty, disable the fancy output, sort, and remove duplicats
  cat recon-targets.txt | hacker-scoper -c google -ch | sort -u

  Example: Manually pick a file, lookup scopes on firebounty, and set explicit-level
  hacker-scoper -f recon-targets.txt -c google -e 2

  Example: Manually pick a file, use custom scopes and out-of-scope files, and set explicit-level
  hacker-scoper -f recon-targets.txt -ins inscope -oos noscope.txt -e 2

List of all possible arguments:
  -c, --company string
      Specify the company name to lookup.

  -cstxt, --check-security-txt
      Whether or not we will try to scrape security.txt from all domains and subdomains

  -r, --reuse string
      Reuse previously generated security.txt lists? (Y/N)
          Only needed if using "-cstxt"

  -f, --file string
      Path to your file containing URLs

  -ins, --inscope-file string
      Path to a custom plaintext file containing scopes

  -oos, --outofcope-file string
      Path to a custom plaintext file containing scopes exclusions

  -e, --explicit-level int
      How explicit we expect the scopes to be:
       1 (default): Include subdomains in the scope even if there's not a wildcard in the scope
       2: Include subdomains in the scope only if there's a wildcard in the scope
       3: Include subdomains in the scope only if they are explicitly within the scope

  -ch, --chain-mode
      In "chain-mode" we only output the important information. No decorations.
            Default: false

  --fire string
      Set this to specify a path the FireBounty JSON.
```

The firebounty json is automatically updated every 24hs

list example:
```powershell
example.com
dev.example.com
1.dev.example.com
2.dev.example.com
ads.example.com
192.168.1.10
192.168.2.10
192.168.2.8
```

Custom scopes file example:
```powershell
*.example.com
*.sub.domain.example.com
192.168.1.10
192.168.2.1/24
FE80:0000:0000:0000:0202:B3FF:FE1E:8329
FE80::0202:B3FF:FE1E:8329
FE80::0204:B3FF::/24
```

Custom out-of-scope file example:
```javascript
community.example.com
thirdparty.example.com
*.thirdparty.example.com
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
