package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const firebountyAPIURL = "https://firebounty.com/api/v1/scope/all/url_only/"
const firebountyJSONFilename = "firebounty-scope-url_only.json"

//https://tutorialedge.net/golang/parsing-json-with-golang/
type Scope struct {
	Scope      string //either a domain, or a wildcard domain
	Scope_type string //we only care about "web_application"
}

type Program struct {
	Firebounty_url string //url.URL not allowed appearently
	Scopes         struct {
		In_scopes     []Scope
		Out_of_scopes []Scope
	}
	Slug string
	Tag  string
	Url  string //url.URL not allowed appearently
	Name string
}

type WhiteLists struct {
	Regex        string //can't be "*regexp.Regexp" because they're actually domain wildcards
	Program_slug string
}

type Firebounty struct {
	White_listed []WhiteLists
	Pgms         []Program
}

var chainMode bool
var targetsListFilepath string

const colorReset = "\033[0m"
const colorYellow = "\033[33m"
const colorRed = "\033[31m"

func main() {

	var company string
	var stxt bool
	var reuseList string  //should only be "Y", "N" or ""
	var explicitLevel int //should only be [0], 1, or 2
	var scopesListFilepath string
	var outofScopesListFilepath string
	var usedstdin bool
	usedstdin = false

	const usage = `Usage: ./hacker-scoper --file /path/to/targets [--company company | --custom-inscopes-file /path/to/inscopes [--custom-outofcopes-file] /path/to/outofscopes] [--explicit-level INT] [--reuse Y/N] [--chain-mode]
Example: ./hacker-scoper --file /home/kali/Downloads/recon-targets.txt --company google --explicit-level 2
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
	
  NOTE: Targets won't be matched if they don't have a valid scheme:
    ✅ http://target.com
    ✅ mongodb://127.0.0.1
    ❌ target.com
    ❌ 127.0.0.1
`

	flag.StringVar(&company, "c", "", "Specify the company name to lookup.")
	flag.StringVar(&company, "company", "", "Specify the company name to lookup.")
	flag.StringVar(&targetsListFilepath, "f", "", "Path to your file containing URLs")
	flag.StringVar(&targetsListFilepath, "file", "", "Path to your file containing URLs")
	flag.StringVar(&scopesListFilepath, "ins", "", "Path to a custom plaintext file containing scopes")
	flag.StringVar(&scopesListFilepath, "inscope-file", "", "Path to a custom plaintext file containing scopes")
	flag.StringVar(&outofScopesListFilepath, "oos", "", "Path to a custom plaintext file containing scopes exclusions")
	flag.StringVar(&outofScopesListFilepath, "outofcope-file", "", "Path to a custom plaintext file containing scopes exclusions")
	flag.IntVar(&explicitLevel, "e", 1, "Level of explicity expected. ([1]/2/3)")
	flag.IntVar(&explicitLevel, "explicit-level", 1, "Level of explicity expected. ([1]/2/3)")
	flag.BoolVar(&stxt, "cstxt", false, "Whether or not we will try to scrape security.txt from all domains and subdomains")
	flag.BoolVar(&stxt, "check-security-txt", false, "Whether or not we will try to scrape security.txt from all domains and subdomains")
	flag.StringVar(&reuseList, "r", "", "Reuse previously generated lists? (Y/N)")
	flag.StringVar(&reuseList, "reuse", "", "Reuse previously generated lists? (Y/N)")
	flag.BoolVar(&chainMode, "ch", false, "In \"chain-mode\" we only output the important information. No decorations.")
	flag.BoolVar(&chainMode, "chain-mode", false, "In \"chain-mode\" we only output the important information. No decorations.")
	//https://www.antoniojgutierrez.com/posts/2021-05-14-short-and-long-options-in-go-flags-pkg/
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	banner := `
'||                      '||                      '                                                 
 || ..    ....     ....   ||  ..    ....  ... ..     ....    ....    ...   ... ...    ....  ... ..  
 ||' ||  '' .||  .|   ''  || .'   .|...||  ||' ''   ||. '  .|   '' .|  '|.  ||'  || .|...||  ||' '' 
 ||  ||  .|' ||  ||       ||'|.   ||       ||       . '|.. ||      ||   ||  ||    | ||       ||     
.||. ||. '|..'|'  '|...' .||. ||.  '|...' .||.      |'..|'  '|...'  '|..|'  ||...'   '|...' .||.    
                                                                            ||                      
                                                                           ''''                     
`

	if !chainMode {
		fmt.Println(banner)
	}

	//validate arguments
	if (explicitLevel != 1) && (explicitLevel != 2) && explicitLevel != 3 {
		var err error
		crash("Invalid explicit-level selected", err)
	}

	if !chainMode {
		if explicitLevel != 3 {
			var howMany string
			if explicitLevel == 2 {
				howMany = "Some"
			} else {
				//explicitLevel = 1
				howMany = "A lot of"
			}
			warning("(--explicit-level=" + strconv.Itoa(explicitLevel) + ") " + howMany + " scopes might appear as duplicates if they are explicitly in the scope, and also covered by a wildcard. Consider running uniq on the output file.")
		}
	}

	//overwrite whathever was feeded to targetsListFilepath with the stdin input
	//https://stackoverflow.com/a/26567513/11490425
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		//data is being piped to stdin

		var stdinInput string

		tmpFile := createFile("hacker-scoper_stdin_scopes_tmp-file.txt", os.TempDir())

		//read stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdinInput += "\n" + scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			crash("bufio couldn't read stdin correctly.", err)
		}

		//write to disk
		err := os.WriteFile(os.TempDir()+"/hacker-scoper_stdin_scopes_tmp-file.txt", []byte(stdinInput), 0666)
		if err != nil {
			crash("Couldn't save write to tmp file.", err)
		}

		popLine(tmpFile)
		tmpFile.Close()

		targetsListFilepath = os.TempDir() + "/hacker-scoper_stdin_scopes_tmp-file.txt"
		usedstdin = true

	} //else { //stdin is from a terminal }

	//clean targetsListFilepath path for +speed
	targetsListFilepath = filepath.Clean(targetsListFilepath)

	//Verify existance of the targetsListFilepath file
	//https://stackoverflow.com/a/12518877/11490425
	if _, err := os.Stat(targetsListFilepath); err == nil {
		// path/to/whatever exists

		//check if security.txt exists
		//based on https://github.com/yeswehack/yeswehack_vdp_finder
		if stxt {

			var outputFileName string = "security-txt_URLs.txt"

			//attempt to create the file to later write the result URLs
			//https://stackoverflow.com/a/12518877/11490425
			if _, err := os.Stat(outputFileName); err == nil {
				//security-txt_URLs.txt exists
				//reuse?
				if reuseList == "" {
					fmt.Println("Previous " + outputFileName + " file found. Do you want to reuse it? ([Y]/N): ")
					fmt.Scanln(&reuseList)
				}
				if reuseList == "N" {
					//delete the old file
					err := os.Remove(outputFileName) // remove a single file
					if err != nil {
						fmt.Println(err)
					}

					createFile(outputFileName, "")

					//open the file
					//https://stackoverflow.com/a/16615559/11490425
					file, err := os.Open(targetsListFilepath)
					if err != nil {
						crash("Could not open targets URL-List file", err)
					}

					//scan the file using bufio
					scanner := bufio.NewScanner(file)

					//for each line in the file..
					//Scanner will error with lines longer than 65536 characters. If you know your line length is greater than 64K
					for scanner.Scan() {
						//https://gobyexample.com/url-parsing
						URL, err := url.Parse(scanner.Text())
						if err != nil {
							crash("Could not read a line on the input file. Lines longer than 65536 characters are not allowed. If this is an issue for you, open an issue.", err)
						}

						//get only domains & subdomains from page which start with HTTP/S
						if URL.Scheme == "http" || URL.Scheme == "https" {
							//remove query parameters from the URL
							//https://stackoverflow.com/a/55299809/11490425
							URL.RawQuery = ""

							//add the security.txt path
							//TODO: despite security.txt also being valid at the root directory, for now we will only look for it on the .well-known directory
							URL.Path = URL.Path + "/.well-known/security.txt"

							//open the output file for writing
							f, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_WRONLY, 0600)
							if err != nil {
								crash("Coulnd't open file "+outputFileName+" for writing security.txt URLs.", err)
							}

							//append the URL to the file
							if _, err = f.WriteString("\n" + URL.String()); err != nil {
								crash("Couldn't append a line to the security.txt-check output file.", err)
							}

							f.Close()

						}

						if err := scanner.Err(); err != nil {
							crash("Could not read URL List file successfully", err)
						}
					}

					file.Close()

					//pop the first line of the list, because it contains an unnecesary linejump
					//the line popper has it's own error handling.
					outputFile, _ := os.OpenFile(outputFileName, os.O_RDWR, 0666)
					popLine(outputFile)
					outputFile.Close()

				} //else { //user wants to reuse the list }

			} else if errors.Is(err, os.ErrNotExist) {
				//security-txt_URLs.txt does NOT exist
				//create it
				createFile(outputFileName, "")

			} else {
				// Schrodinger: file may or may not exist. See err for details.
				panic(err)

			}

			//open the file
			//https://stackoverflow.com/a/16615559/11490425
			file, err := os.Open(outputFileName)
			if err != nil {
				crash("Could not open the security.txt output file", err)
			}

			//Read the output file line per line
			//scan the file using bufio
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				const titleRegex = `(<title>).*(</title>)`
				allHTTPErrors := []int{300, 301, 302, 303, 304, 305, 306, 307, 308, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 421, 422, 423, 491, 424, 491, 425, 426, 427, 428, 429, 430, 431, 451, 500, 501, 502, 503, 504, 505, 506, 507, 508, 584, 509, 510, 511}

				//TODO: customizeable timeout
				//https://stackoverflow.com/a/25344458/11490425
				client := http.Client{
					Timeout: 5 * time.Second,
				}
				resp, err := client.Get(scanner.Text())
				if err != nil {
					//do not panic if a request fails
					if !chainMode {
						fmt.Println("[HTTP Fail]: Request failed for " + scanner.Text())
					}

				} else {
					if resp.StatusCode == 200 {
						body, err := ioutil.ReadAll(resp.Body)
						resp.Body.Close()
						if err != nil {
							fmt.Println(err)
						}

						regex, _ := regexp.Compile(titleRegex)
						result := regex.FindAllString(string(body), 2)
						var flag bool

					html:
						for titleCounter := 0; titleCounter < len(result); titleCounter++ {
							for i := 0; i < len(allHTTPErrors); i++ {
								if strings.Contains(result[titleCounter], strconv.Itoa(allHTTPErrors[i])) {
									if !chainMode {
										fmt.Println("ERROR - STATUS CODE " + strconv.Itoa(allHTTPErrors[i]))
									}
									flag = true
									break html
								}
							}

						}

						if !flag {
							//security.txt found!
							fmt.Println("[+] security.txt found at: " + scanner.Text())
							fmt.Println(string(body))
						}

					}
				}

			}
			if err := scanner.Err(); err != nil {
				crash("Could not read URL List file successfully", err)
			}
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		err = nil
		crash("The provided URL list file does not exist!", err)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		crash("Could not verify existance of provided URL List file!", err)

	}

	if company == "" && scopesListFilepath == "" {
		var err error
		crash("A company name is required to smartly weed-out out-of-scope URLs", err)
	} else {

		//default value. user will use the integrated scope list
		if scopesListFilepath == "" {
			if firebountyJSONFileStats, err := os.Stat(firebountyJSONFilename); err == nil {
				// path/to/whatever exists
				//check age. if age > 24hs
				yesterday := time.Now().Add(-24 * time.Hour)
				if firebountyJSONFileStats.ModTime().Before(yesterday) {
					if !chainMode {
						fmt.Println("[INFO]: 24hs have passed since the last update to the scopes file. Updating...")
					}
					updateFireBountyJSON()
				}

			} else if errors.Is(err, os.ErrNotExist) {
				//path/to/whatever does not exist
				if !chainMode {
					fmt.Println("[INFO]: Downloading scopes file...")
				}

				updateFireBountyJSON()

			} else {
				// Schrodinger: file may or may not exist. See err for details.
				panic(err)
			}

			//open json
			jsonFile, err := os.Open(firebountyJSONFilename)
			if err != nil {
				crash("Couldn't open firebounty JSON. Maybe run \"chmod 777 "+firebountyJSONFilename+"\"? ", err)
			}

			//read the json file as bytes
			byteValue, _ := ioutil.ReadAll(jsonFile)
			jsonFile.Close()

			var firebountyJSON Firebounty
			err = json.Unmarshal(byteValue, &firebountyJSON)
			if err != nil {
				crash("Couldn't parse firebountyJSON into pre-defined struct.", err)
			}

			//for every company...
			for companyCounter := 0; companyCounter < len(firebountyJSON.Pgms); companyCounter++ {
				fcompany := strings.ToLower(firebountyJSON.Pgms[companyCounter].Name)
				if strings.Contains(fcompany, company) {
					//match found!
					if !chainMode {
						fmt.Println("Company: " + fcompany)
					}
					//for every scope in the program
					for scopeCounter := 0; scopeCounter < len(firebountyJSON.Pgms[companyCounter].Scopes.In_scopes); scopeCounter++ {
						//if the scope type is "web_application" and it's not empty
						if firebountyJSON.Pgms[companyCounter].Scopes.In_scopes[scopeCounter].Scope_type == "web_application" && firebountyJSON.Pgms[companyCounter].Scopes.In_scopes[scopeCounter].Scope != "" {

							scope := firebountyJSON.Pgms[companyCounter].Scopes.In_scopes[scopeCounter].Scope

							if !chainMode {
								//alert the user about potentially mis-configured bug-bounty program
								if scope[0:4] == "com." || scope[0:4] == "org." {
									warning("Scope starting with \"com.\" or \"org. found. This may be a sign of a misconfigured bug bounty program. Consider editing the \"" + firebountyJSONFilename + " file and removing the faulty entries. Also, report the failure to the mainters of the bug bounty program.")
								}
							}

							parseScopesWrapper(scope, explicitLevel, targetsListFilepath, outofScopesListFilepath, firebountyJSON.Pgms[companyCounter].Scopes.Out_of_scopes, usedstdin)

						}
					}
				}
			}
		} else {
			//user chose to use their own scope list

			if _, err := os.Stat(scopesListFilepath); err == nil {
				// path/to/whatever exists

				//when using this custom scope, most likely there will be more targets than scopes, so we will nest scopes->targets for more efficiency

				//open the file
				//https://stackoverflow.com/a/16615559/11490425
				scopesFile, err := os.Open(scopesListFilepath)
				if err != nil {
					crash("Could not open "+scopesListFilepath, err)
				}

				//Read the file line per line using bufio
				scopesScanner := bufio.NewScanner(scopesFile)

				for scopesScanner.Scan() {
					parseScopesWrapper(scopesScanner.Text(), explicitLevel, targetsListFilepath, outofScopesListFilepath, nil, usedstdin)
				}
				scopesFile.Close()

			} else if errors.Is(err, os.ErrNotExist) {
				//path/to/whatever does not exist
				err = nil
				crash(scopesListFilepath+" does not exist.", err)

			} else {
				// Schrodinger: file may or may not exist. See err for details.
				panic(err)
			}
		}

	}

	if(usedstdin){
		//Developers using temporary files are expected to clean up after themselves.
		//https://superuser.com/a/296827
		os.Remove(targetsListFilepath)
	}
}

//path must not have the end bar (/)
func createFile(file string, path string) *os.File {
	outputFile, err := os.Create(path + "/" + file)
	if err != nil {
		panic(err)
	}

	return outputFile
}

//https://stackoverflow.com/a/30948278/11490425
func popLine(f *os.File) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	line, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return nil, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return line, nil
}

func updateFireBountyJSON() {
	// path/to/whatever does *not* exist
	//get the big JSON from the API
	jason, err := http.Get(firebountyAPIURL)
	if err != nil {
		crash("Could not download scopes from firebounty at: "+firebountyAPIURL, err)
	}

	//read the contents of the request
	body, err := ioutil.ReadAll(jason.Body)
	jason.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	//delete the previous file (if it even exists)
	os.Remove(firebountyJSONFilename)

	//write to disk
	err = os.WriteFile(firebountyJSONFilename, []byte(string(body)), 0666)
	if err != nil {
		crash("Couldn't save firebounty json to disk as"+firebountyJSONFilename, err)
	}

	if !chainMode {
		fmt.Println("[INFO]: Scopes file updated succesfully.")
	}

}

//we may recieve one like the following as scope:
// example.com
// *.example.com
// 192.168.0.1
// 192.168.0.1/24
// 192.168.0.1
// 192.168.0.1/24
func parseScopes(scope string, isWilcard bool, targetsListFilepath string, outofScopesListFilepath string, firebountyOutOfScopes []Scope, usedstdin bool) {
	schemedScope := "http://" + scope

	var CIDR *net.IPNet
	var parseAsIP bool
	var scopeURL *url.URL
	var err error

	//attempt to parse current scope as a CIDR range
	_, CIDR, _ = net.ParseCIDR(scope)
	scopeIP := net.ParseIP(scope)
	//if we can parse the scope as a CIDR range or as an IP address:
	if scopeIP.String() != "<nil>" || CIDR != nil {
		parseAsIP = true
	} else {
		parseAsIP = false
		scopeURL, err = url.Parse(schemedScope)
		if err != nil {
			if !chainMode {
				warning("Couldn't parse the scope " + scope + " as a valid URL.")
			}
			return
		}
	}

	//open the user-supplied URL list
	file, err := os.Open(targetsListFilepath)
	if err != nil {
		crash("Could not open your provided URL list file", err)
	}

	//Read the URLs file line per line
	//scan using bufio
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		//attempt to parse current target as an IP
		var currentTargetURL *url.URL
		currentTargetURL, err = url.Parse(scanner.Text())

		//If we couldn't parse it as is, attempt to add the "https://" prefix
		if (err != nil || currentTargetURL.Host == ""){
			currentTargetURL, err = url.Parse("https://" + scanner.Text())
		}

		portlessHostofCurrentTarget := removePortFromHost(currentTargetURL)
		targetIp := net.ParseIP(portlessHostofCurrentTarget)

		//if it fails...
		if (err != nil || currentTargetURL.Host == "") && !chainMode {
			if(usedstdin){
				warning("STDIN: Couldn't parse " + scanner.Text() + " as a valid URL.")
			} else {
				warning(targetsListFilepath + ": Couldn't parse " + scanner.Text() + " as a valid URL.")
			}

		} else {
			//we were able to parse the target as a URL
			//if we were able to parse the target as an IP, and the scope as an IP or CIDR range
			if targetIp.String() != "" && parseAsIP {
				//if the CIDR range is empty
				if CIDR == nil {
					//Couldn't parse scope as CIDR range, retrying as ip match")
					if targetIp.String() == scopeIP.String() {
						if !isOutOfScope(nil, outofScopesListFilepath, targetIp, firebountyOutOfScopes) {
							if !chainMode {
								fmt.Println("[+] IN-SCOPE: " + scanner.Text())
							} else {
								fmt.Println(scanner.Text())
							}
						}

					}
				} else {
					if CIDR.Contains(targetIp) {
						if !isOutOfScope(nil, outofScopesListFilepath, targetIp, firebountyOutOfScopes) {
							if !chainMode {
								fmt.Println("[+] IN-SCOPE: " + scanner.Text())
							} else {
								fmt.Println(scanner.Text())
							}
						}
					}
				}

			} else {
				//parse the scope & target as URLs

				if isWilcard {
					//parse the scope as a URL

					//if x is a subdomain of y
					//ex: wordpress.example.com with a scope of *.example.com will give a match
					//we DON'T do it by splitting on dots and matching, because that would cause errors with domains that have two top-level-domains (gov.br for example)
					if strings.HasSuffix(removePortFromHost(currentTargetURL), scopeURL.Host) {
						if !isOutOfScope(currentTargetURL, outofScopesListFilepath, nil, firebountyOutOfScopes) {
							if !chainMode {
								fmt.Println("[+] IN-SCOPE: " + scanner.Text())
							} else {
								fmt.Println(scanner.Text())
							}
						}

					}
				} else {
					if removePortFromHost(currentTargetURL) == scopeURL.Host {
						if !isOutOfScope(currentTargetURL, outofScopesListFilepath, nil, firebountyOutOfScopes) {
							if !chainMode {
								fmt.Println("[+] IN-SCOPE: " + scanner.Text())
							} else {
								fmt.Println(scanner.Text())
							}
						}

					}
				}

			}
		}

	}

	file.Close()

	if err := scanner.Err(); err != nil {
		crash("Could not read URL List file successfully", err)
	}
}

func parseScopesWrapper(scope string, explicitLevel int, targetsListFilepath string, outofScopesListFilepath string, firebountyOutOfScopes []Scope, usedstdin bool) {
	//if we have a wildcard domain
	if strings.Contains(scope, "*.") {
		//shorter way of saying if explicitLevel == 2 || explicitLevel ==1
		if explicitLevel != 3 {
			//remove wildcard ("*.")
			scope = strings.ReplaceAll(scope, "*.", "")
			parseScopes(scope, true, targetsListFilepath, outofScopesListFilepath, firebountyOutOfScopes, usedstdin)
		}
	} else if explicitLevel == 1 {
		//this is NOT a wildcard domain, but we'll treat it as such anyway
		parseScopes(scope, true, targetsListFilepath, outofScopesListFilepath, firebountyOutOfScopes, usedstdin)
	} else {
		//this is NOT a wildcard domain. we will parse it explicitly
		parseScopes(scope, false, targetsListFilepath, outofScopesListFilepath, firebountyOutOfScopes, usedstdin)
	}
}

func crash(message string, err error) {
	fmt.Print(string(colorRed) + "[ERROR]: " + message + string(colorReset) + "\n")
	panic(err)
}

func warning(message string) {
	fmt.Print(string(colorYellow) + "[WARNING]: " + message + string(colorReset) + "\n")
}

func removePortFromHost(url *url.URL) string {
	//code readability > efficiency
	portless := strings.Replace(string(url.Host), string(url.Port()), "", 1)
	//obligatory cleanup ("192.168.1.1:" -> "192.168.1.1")
	portless = strings.Replace(portless, ":", "", 1)
	return portless
}

//out-of-scopes are parsed as --explicit-level==2
func isOutOfScope(targetURL *url.URL, outofScopesListFilepath string, targetIP net.IP, firebountyOutOfScopes []Scope) bool {
	var err error

	if outofScopesListFilepath != "" {
		//user chose to use their own out-of-scopes file
		if _, err = os.Stat(outofScopesListFilepath); err == nil {
			// path/to/whatever exists
			//open the file
			//https://stackoverflow.com/a/16615559/11490425
			outOfScopesFile, err := os.Open(outofScopesListFilepath)
			if err != nil {
				crash("Could not open "+outofScopesListFilepath, err)
			}

			//Read the file line per line using bufio
			outofScopeScanner := bufio.NewScanner(outOfScopesFile)

			for outofScopeScanner.Scan() {

				if parseOutOfScopes(targetURL, outofScopeScanner.Text(), targetIP) {
					return true
				}
			}
			outOfScopesFile.Close()
			return false

		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			crash("OutOfScopes file supplied, but it does not exist!", err)

		} else {
			// Schrodinger: file may or may not exist. See err for details.
			crash("Couldn't verify existance of outofscopesFile", err)

		}
	} else {
		//check target agains firebounty out-of-scopes
		//for every outOfScope
		for outOfScopeCounter := 0; outOfScopeCounter < len(firebountyOutOfScopes); outOfScopeCounter++ {
			//if the scope_type is web_application and it's not empty
			if firebountyOutOfScopes[outOfScopeCounter].Scope_type == "web_application" && firebountyOutOfScopes[outOfScopeCounter].Scope != "" {
				outOfScope := firebountyOutOfScopes[outOfScopeCounter].Scope
				if !chainMode {
					//alert the user about potentially mis-configured bug-bounty program
					if outOfScope[0:4] == "com." || outOfScope[0:4] == "org." {
						warning("Scope starting with \"com.\" or \"org. found. This may be a sign of a misconfigured bug bounty program. Consider editing the \"" + firebountyJSONFilename + " file and removing the faulty entries. Also, report the failure to the maintainers of the bug bounty program.")
					}
				}
				if parseOutOfScopes(targetURL, outOfScope, targetIP) {
					return true
				}
			}

		}
	}

	//if we got no matches for any outOfScope
	return false
}

func parseOutOfScopes(targetURL *url.URL, outOfScope string, targetIP net.IP) bool {
	if targetURL != nil {
		//parse target as a URL
		isWildcard := strings.Contains(outOfScope, "*.")
		outOfScopeURL, err := url.Parse("https://" + outOfScope)
		if err != nil {
			if !chainMode {
				warning("Couldn't parse out-of-scope \"" + outOfScope + "\" as a URL.")
			}
			return false
		}

		if isWildcard {
			//if x is a subdomain of y
			//ex: wordpress.example.com with a scope of *.example.com will give a match
			//we DON'T do it by splitting on dots and matching, because that would cause errors with domains that have two top-level-domains (gov.br for example)
			if strings.HasSuffix(removePortFromHost(targetURL), outOfScopeURL.Host) {
				return true

			}
		} else {
			if removePortFromHost(targetURL) == outOfScopeURL.Host {
				return true

			}
		}
	} else {
		//IP mode
		//attempt to parse current outOfScope as an IP
		outOfScopeIp := net.ParseIP(outOfScope)
		//if we can parse the current outOfScope as an IP...
		if outOfScopeIp != nil {
			//try IP match
			if targetIP.String() == outOfScopeIp.String() {
				return true
			}
		}
	}

	//if nothing matched
	return false
}
