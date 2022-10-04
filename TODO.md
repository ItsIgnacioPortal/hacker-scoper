- [x] Add **Verbose** (`-v, --verbose`) flag that will show what scopes were detected in firebounty/hackerone, etc. This option is incompatible with chain-mode.
- [x] Add **UNSURE**: If a URL is not within scope, but also not outside of the scope, report it as _UNSURE_. Set `-iu / --include-unsure` to enable.
- [x] Add **Allow user to select a company when there's multiple matches for the same company name**
- [ ] Add **BBaaS API private program scope query**: If a company wasn't found on firebounty, ask the user for an API key. 

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

REFERENCE: https://github.com/root4loot/rescope#features
- [ ] Add **Combine private and public scopes**
- [ ] Add **Resolves conflicting includes/excludes**
- [ ] Put the usage options in a fancy table on the README.
- [ ] Add **Define multiple inscopes sources and combine them**
- [ ] Add **Define private scopes by copy/pasting target definitions from pretty much anywhere**
- [ ] Add Unit Tests
- [ ] Add fully automated chocolatey releases
- [ ] If the company name didn't match firebounty, nor any BBaaS platform scope, attempt to get the scope using an ASN
- [ ] Improve the _Attempt to scrape security.txt files from your targets_ feature
