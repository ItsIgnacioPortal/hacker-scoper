- [x] Add **Verbose** (`-v, --verbose`) flag that will show what scopes were detected in firebounty/hackerone, etc. This option is incompatible with chain-mode.
- [x] Add **UNSURE**: If a URL is not within scope, but also not outside of the scope, report it as _UNSURE_. Set `-iu / --include-unsure` to enable.
- [x] Add **Allow user to select a company when there's multiple matches for the same company name**
- [ ] Comply with [OpenSSF Best Practices](https://www.bestpractices.dev)
- [ ] Add **BBaaS API private program scope query**: If a company wasn't found on firebounty, ask the user for an API key. REFERENCE: https://github.com/root4loot/rescope#features

	"Do you want to add a HackerOne API Key to get private bug-bounty program scopes? ([Yes]/No/Later): "
	If "No" is selected, allow user to register an API key later, using `--hackerone API_KEY`, `--bugcrowd API_KEY`, `--intigriti API_KEY`, etc.    
	List of Bug-Bounty as a Service platforms (BBaaS): 
	- [ ] bugcrowd.com
	- [ ] hackerone.com
	- [ ] hackenproof.com
	- [ ] intigriti.com
	- [ ] openbugbounty.com
	- [ ] yeswehack.com
	- [ ] bugbounty.jp
	- [ ] federacy.com
- [X] Add **Hostname-only output** 
- [X] Put the usage options in a fancy table on the README.
- [ ] Add fully automated chocolatey releases
- [ ] If the company name didn't match firebounty, nor any BBaaS platform scope, attempt to get the scope using an ASN
- [ ] Improve the _Attempt to scrape security.txt files from your targets_ feature
- [ ] Add **Combine private and public scopes**
- [ ] Add **Resolves conflicting includes/excludes**
- [ ] Add **Define multiple inscopes sources and combine them**
