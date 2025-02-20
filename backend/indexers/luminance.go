package indexers

import (
	"backend/database"
	"backend/helpers"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func LoginAndGetCookiesLuminance(username string, password string, cookiesProwlarr string, loginURL string, indexerInfo gjson.Result) string {

	if cookiesProwlarr != "" {
		// uses 2FA, cookies directly provided by prowlarr, no need to login
		return strings.Replace(cookiesProwlarr, " ", "", -1)
	}

	// body := indexerInfo.Get("login.body").String()
	fields := indexerInfo.Get("login.fields").Map()
	formData := url.Values{}
	formData.Add(fields["username"].String(), username)
	formData.Add(fields["password"].String(), password)

	tokens := getHiddenTokensLuminance(loginURL, indexerInfo.Get("domain").Str)

	for key, value := range tokens["inputs"] {
		formData.Add(key, value)
	}

	extraFields := fields["extra"].Map()
	for key, val := range extraFields {
		formData.Add(key, val.String())
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevents redirect
		},
	}

	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	u, _ := url.Parse(loginURL)
	var cookieSlice []*http.Cookie
	for name, value := range tokens["cookies"] {
		cookieSlice = append(cookieSlice, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}
	jar.SetCookies(u, cookieSlice)

	resp, err := client.PostForm(loginURL, formData)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	cookiesStr := ""
	for _, cookie := range cookies {
		cookiesStr += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
	}
	if !strings.Contains(cookiesStr, "sid") {
		// login failed
		return ""
	}
	cookiesStr = cookiesStr[:len(cookiesStr)-1]
	return cookiesStr
	// return ""
}

func getHiddenTokensLuminance(url string, domain string) map[string]map[string]string {

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", domain)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("DNT", "1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("TE", "trailers")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	tokens := map[string]map[string]string{"inputs": {}, "cookies": {}}
	// fmt.Println(doc.Find("body").Html())

	token := doc.Find("#cinfo")
	tokenName, _ := token.Attr("name")
	// tokenValue, _ := token.Attr("value")
	// this is just some info about the client browser, it is computed by a js script
	tokens["inputs"][tokenName] = "1920|935|24|0"

	token = doc.Find("input[name=\"token\"]")
	tokenName, _ = token.Attr("name")
	tokenValue, _ := token.Attr("value")
	tokens["inputs"][tokenName] = tokenValue

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		tokens["cookies"][cookie.Name] = cookie.Value
	}

	return tokens
}

func ConstructRequestLuminance(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str

	cookieStr := database.GetIndexerCookies(indexerId)
	// indexer is not setup yet, credentials are in prowlarr but session cookie needs to be stored in dasharr's db
	// TODO: move this execution to container initialization
	if cookieStr == "" {
		LoginAndSaveCookies(indexerName, "", "", "", "", indexerId)
		cookieStr = database.GetIndexerCookies(indexerId)
	}
	userPath := getUserPathLuminance(baseUrl, cookieStr)
	req, _ := http.NewRequest("GET", strings.TrimSuffix(baseUrl, "/")+userPath, nil)
	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req)

	return req
}

func getUserPathLuminance(baseUrl string, cookieStr string) string {
	// req, _ := http.NewRequest("", "", nil)
	req, _ := http.NewRequest("GET", baseUrl, nil)

	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req.Cookies())

	req.Header.Add("Cookie", cookieStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if resp.Status == "200 OK" {
		// indexerInfo := helpers.GetIndexerInfo(indexerName)
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		href, found := doc.Find(`a[class="username"]`).Attr("href")
		if !found {
			fmt.Printf("Couldn't find userid for indexer %s", baseUrl)
			return ""
		}
		return href
	} else {
		return ""
	}
}

func ProcessIndexerResponseLuminance(bodyString string, indexerInfo gjson.Result) map[string]interface{} {
	//todo: handle cookie refresh
	results := map[string]interface{}{}
	re := regexp.MustCompile(`([\d\.]+)[ \x{00a0}]?\s?(GiB|GB|MiB|MB|TiB|TB|KiB|KB|B)`)
	extractNumberRegex := regexp.MustCompile(`\d+`)
	extractNumberAfterColumnRegex := regexp.MustCompile(`: (\d+)`)
	// indexerName := indexerInfo.Get("indexer_name").Str

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	// _ = os.WriteFile("test.html", []byte(bodyString), 0644)
	uploadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_amount").Str).Text())
	if len(uploadRegexResult) == 0 {
		fmt.Printf("An error occured while parsing Luminance's response")
		return results
	}
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	results["uploaded_amount"] = helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2])

	var bonusPoints string
	bonusPoints = doc.Find(indexerInfo.Get("scraping.xpaths.bonus_points").Str).Text()
	results["bonus_points"] = strings.ReplaceAll(bonusPoints, ",", "")

	downloadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.downloaded_amount").Str).Text())
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	results["downloaded_amount"] = helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2])

	// requiredRatioRegexResult := extractNumberAfterColumnRegex.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.required_ratio").Str).Text())
	// cleanRequiredRatio, _ := strconv.ParseFloat(requiredRatioRegexResult[1], 64)
	// results["required_ratio"] = helpers.AnyUnitToBytes(cleanRequiredRatio, requiredRatioRegexResult[2])

	invitedRegexResult := extractNumberAfterColumnRegex.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.invited").Str).Text())
	results["invited"], _ = strconv.ParseFloat(invitedRegexResult[1], 64)

	seedingSizeXpath := indexerInfo.Get("scraping.xpaths.seeding_size")
	seedingSizeRegexResult := re.FindStringSubmatch(doc.Find(seedingSizeXpath.Str).Text())
	cleanSeedingSize, _ := strconv.ParseFloat(seedingSizeRegexResult[1], 64)
	results["seeding_size"] = helpers.AnyUnitToBytes(cleanSeedingSize, seedingSizeRegexResult[2])

	uploaded_torrents := extractNumberRegex.FindString(doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_torrents").Str).Text())
	results["uploaded_torrents"] = uploaded_torrents

	snatchedRegex := extractNumberRegex.FindString(doc.Find(indexerInfo.Get("scraping.xpaths.snatched").Str).Text())
	results["snatched"], _ = strconv.Atoi(snatchedRegex)

	seedingRegex := extractNumberAfterColumnRegex.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.seeding").Str).Text())
	results["seeding"], _ = strconv.Atoi(seedingRegex[1])

	leechingRegex := extractNumberAfterColumnRegex.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.leeching").Str).Text())
	results["leeching"], _ = strconv.Atoi(leechingRegex[1])

	// ratioXpath := indexerInfo.Get("scraping.xpaths.ratio")
	// ratio := doc.Find(indexerInfo.Get(ratioXpath.Str).Str).Text()
	// results["ratio"] = ratio

	torrent_comments := extractNumberRegex.FindString(doc.Find(indexerInfo.Get("scraping.xpaths.torrent_comments").Str).Text())
	results["torrent_comments"] = torrent_comments

	forum_posts := extractNumberRegex.FindString(doc.Find(indexerInfo.Get("scraping.xpaths.forum_posts").Str).Text())
	results["forum_posts"] = forum_posts

	freeleechTokensXpath := indexerInfo.Get("scraping.xpaths.freeleech_tokens")
	if freeleechTokensXpath.Exists() {
		freeleech_tokens := extractNumberRegex.FindString(doc.Find(freeleechTokensXpath.Str).Text())
		results["freeleech_tokens"] = freeleech_tokens
	}

	// warned, _ := strconv.Atoi(doc.Find(indexerInfo.Get("scraping.xpaths.warned").Str).Text())
	// results["warned"] = warned > 0

	fmt.Println(results)

	return results
}
