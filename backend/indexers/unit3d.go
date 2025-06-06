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

func LoginAndGetCookiesUnit3d(username string, password string, twoFaCode string, loginURL string, domain string) string {
	formData := url.Values{}
	formData.Add("username", username)
	formData.Add("password", password)
	formData.Add("_username", "")
	formData.Add("remember", "on")

	tokens := getHiddenTokensUnit3d(loginURL, domain)

	for key, value := range tokens["inputs"] {
		formData.Add(key, value)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if twoFaCode != "" {
				// redirect needs to be done for 2fa
				return nil
			} else {
				// Prevents redirect
				return http.ErrUseLastResponse
			}
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

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Host", domain)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Referer", loginURL)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", fmt.Sprintf("https://%s", domain))
	req.Header.Add("DNT", "1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("TE", "trailers")

	// fmt.Println("Cookies in jar:", jar.Cookies(u))
	// fmt.Println(formData)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	if twoFaCode != "" {
		resp = twoFaHandlerUnit3d(resp, twoFaCode, cookies, domain)
		cookies = resp.Cookies()
	}

	// fmt.Println(resp.StatusCode)
	cookiesStr := ""
	for _, cookie := range cookies {
		// fmt.Println(cookie)
		cookiesStr += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
	}
	// this condition doesn't work on all trackers, find a better solution to see if the login failed
	if !strings.Contains(cookiesStr, "laravel_session") {
		// login failed
		return ""
	}
	cookiesStr = cookiesStr[:len(cookiesStr)-1]
	return cookiesStr
}

func twoFaHandlerUnit3d(loginResponse *http.Response, twoFaCode string, loginCookies []*http.Cookie, domain string) *http.Response {
	doc, err := goquery.NewDocumentFromReader(loginResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	tokens := map[string]string{}

	formData := url.Values{}

	// todo : move the xpaths from hard-coded to config.json

	// _token
	token := doc.Find("main section form input:nth-of-type(1)")
	fmt.Println(token)
	tokenName, _ := token.Attr("name")
	tokenValue, _ := token.Attr("value")
	tokens[tokenName] = tokenValue

	// _captcha
	token = doc.Find("main section form input:nth-of-type(2)")
	tokenName, _ = token.Attr("name")
	tokenValue, _ = token.Attr("value")
	tokens[tokenName] = tokenValue

	// random string
	token = doc.Find("main section form input:nth-of-type(3)")
	tokenName, _ = token.Attr("name")
	tokenValue, _ = token.Attr("value")
	tokens[tokenName] = tokenValue

	tokens["recovery_code"] = ""
	tokens["_username"] = ""

	for key, value := range tokens {
		formData.Add(key, value)
	}
	formData.Add("code", twoFaCode)

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s/two-factor-challenge", domain), strings.NewReader(formData.Encode()))

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
	req.Header.Add("Origin", "https://"+domain)
	req.Header.Add("Referer", "https://"+domain+"/two-factor-challenge")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range loginCookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevents redirect
			return http.ErrUseLastResponse
			// return nil
		},
	}
	twoFaResponse, _ := client.Do(req)

	return twoFaResponse
}

func getHiddenTokensUnit3d(url string, domain string) map[string]map[string]string {

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
	token := doc.Find("html > body > main > section > form > input:nth-of-type(2)")
	tokenName, _ := token.Attr("name")
	tokenValue, _ := token.Attr("value")
	tokens["inputs"][tokenName] = tokenValue

	doc.Find("html > body > main > section > form > input:nth-of-type(1)").Each(func(i int, s *goquery.Selection) {
		tokenName, nameExists := s.Attr("name")
		tokenValue, valueExists := s.Attr("value")
		if nameExists && valueExists {
			tokens["inputs"][tokenName] = tokenValue
		} else {
			fmt.Println("Token 2 attributes missing")
		}
	})

	doc.Find("html > body > main > section > form > input:nth-of-type(3)").Each(func(i int, s *goquery.Selection) {
		tokenName, nameExists := s.Attr("name")
		tokenValue, valueExists := s.Attr("value")
		if nameExists && valueExists {
			tokens["inputs"][tokenName] = tokenValue
		} else {
			fmt.Println("Token 2 attributes missing")
		}
	})

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		tokens["cookies"][cookie.Name] = cookie.Value
	}

	return tokens
}

func ConstructRequestUnit3d(indexerName string, indexerId int64) *http.Request {
	indexerInfo := helpers.GetIndexerInfo(indexerName)
	username := database.GetIndexerUsername(indexerId)
	baseUrl := indexerInfo.Get("base_url").Str + "users/" + username
	// fmt.Println(baseUrl)

	// indexerInfo := helpers.GetIndexerInfo(indexerName)

	req, _ := http.NewRequest("GET", baseUrl, nil)

	cookieStr := database.GetIndexerCookies(indexerId)
	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req)

	return req
}

func ProcessIndexerResponseUnit3d(bodyString string, indexerInfo gjson.Result) map[string]interface{} {
	//todo: handle cookie refresh
	results := map[string]interface{}{}
	re := regexp.MustCompile(`([\d\.]+)[ \x{00a0}]?\s?(GiB|MiB|TiB|KiB|B)`)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	uploadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_amount").Str).Text())
	if len(uploadRegexResult) == 0 {
		fmt.Printf("An error occured while parsing %s's response", indexerInfo.Get("indexer_name").Str)
		return results
	}
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	results["uploaded_amount"] = helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2])

	downloadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.downloaded_amount").Str).Text())
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	results["downloaded_amount"] = helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2])

	bufferRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.buffer").Str).Text())
	cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
	results["buffer"] = helpers.AnyUnitToBytes(cleanBuffer, downloadRegexResult[2])

	seedingSizeRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.seeding_size").Str).Text())
	cleanSeedingSize, _ := strconv.ParseFloat(seedingSizeRegexResult[1], 64)
	results["seeding_size"] = helpers.AnyUnitToBytes(cleanSeedingSize, seedingSizeRegexResult[2])

	bonusPoints := doc.Find(indexerInfo.Get("scraping.xpaths.bonus_points").Str).Text()
	results["bonus_points"] = strings.ReplaceAll(bonusPoints, " ", "")

	uploaded_torrents := doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_torrents").Str).Text()
	results["uploaded_torrents"] = uploaded_torrents

	snatched := doc.Find(indexerInfo.Get("scraping.xpaths.snatched").Str).Text()
	results["snatched"] = snatched

	seeding := doc.Find(indexerInfo.Get("scraping.xpaths.seeding").Str).Text()
	results["seeding"] = seeding

	leeching := doc.Find(indexerInfo.Get("scraping.xpaths.leeching").Str).Text()
	results["leeching"] = leeching

	ratio := doc.Find(indexerInfo.Get("scraping.xpaths.ratio").Str).Text()
	results["ratio"] = ratio

	torrent_comments := doc.Find(indexerInfo.Get("scraping.xpaths.torrent_comments").Str).Text()
	results["torrent_comments"] = torrent_comments

	forum_posts := doc.Find(indexerInfo.Get("scraping.xpaths.forum_posts").Str).Text()
	results["forum_posts"] = forum_posts

	freeleech_tokens := doc.Find(indexerInfo.Get("scraping.xpaths.freeleech_tokens").Str).Text()
	results["freeleech_tokens"] = freeleech_tokens

	warned, _ := strconv.Atoi(doc.Find(indexerInfo.Get("scraping.xpaths.warned").Str).Text())
	results["warned"] = warned > 0

	// fmt.Println(results)

	return results
}
