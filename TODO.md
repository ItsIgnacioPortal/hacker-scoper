- Add **Verbose** (`-v, --verbose`) flag that will show what scopes were detected in firebounty/hackerone, etc. This option is incompatible with chain-mode.
- Add **UNSURE**: If a url is not within scope, but also not outside of the scope, report it as _UNSURE_. These urls will not be included when using chain-mode, unless `-iu / --include-unsure`
- Add **Select a company name when there's multiple matches**
- Add **HackerOne, BugCrowd and Intigriti API private program scope query**: If a company wasn't found on firebounty, ask the user for an API key. 
    "Do you want to add a HackerOne API Key to get private bug-bounty program scopes? ([Yes]/No/Later): "
    If "No" is selected, allow user to register an API key later, using `--hackerone API_KEY`, `--bugcrowd API_KEY`, `--intigriti API_KEY`
- Add **Domains-only output** 
- Add **Output ignore user/pass**: Remove duplicates ignoring the URL username/password
- Add **Output ignore protocol**:  Remove duplicates ignoring the URL protocol
- Add **Output ignore segment**:   Remove duplicates ignoring the URL segment
- Add **Output ignore path**:      Remove duplicates ignoring the URL path
- Add **Reverse matching**: `-r, --reverse` only show out-of-scope items

REFERENCE: https://github.com/root4loot/rescope#features
- Add **Combine private and public scopes**
- Add **Resolves conflicting includes/excludes**

Support Bug-Bounty Services (BBaaS)
 - bugcrowd.com
 - hackerone.com
 - hackenproof.com
 - intigriti.com
 - openbugbounty.com
 - yeswehack.com
 - bugbounty.jp
 - federacy.com

- Put the options in a table on the README.
- Add **Define multiple scopes**
- Add **Define private scopes by copy/pasting target definitions from pretty much anywhere**
- Add Unit Tests
- Add fully automated chocolatey releases
- If the company name didn't match firebounty, nor any BBaaS platform scope, attempt to get the scope from the ASN