package trackers

import (
	"backend/database"
	"backend/helpers"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func LoginAndGetCookiesUnit3d(username string, password string, loginURL string, domain string) string {
	formData := url.Values{}
	formData.Add("username", username)
	formData.Add("password", password)
	formData.Add("_username", "")

	tokens := getHiddenTokensUnit3d(loginURL, domain)

	for key, value := range tokens["inputs"] {
		formData.Add(key, value)
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
	// fmt.Println(resp.StatusCode)
	cookiesStr := ""
	for _, cookie := range cookies {
		fmt.Println(cookie)
		cookiesStr += cookie.String()
	}
	return cookiesStr
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

func ConstructRequestUnit3d(trackerConfig gjson.Result, trackerName string, indexerId int64) *http.Request {
	configFile, _ := os.ReadFile(fmt.Sprintf("config/trackers/%s.json", trackerName))
	configFileJson := gjson.Parse(string(configFile))
	baseUrl := configFileJson.Get("base_url").Str + "user"

	indexerInfo := helpers.GetIndexerInfo(trackerName)

	req, _ := http.NewRequest("GET", baseUrl, nil)

	cookie := &http.Cookie{
		Name:  indexerInfo.Get("login.cookie_name").Str,
		Value: database.GetIndexerCookie(indexerId),
		Path:  "/",
		// Expires: time.Now().Add(24 * time.Hour), // Optional: Set expiry time
	}

	req.AddCookie(cookie)
	return req
}

func ProcessTrackerResponseUnit3d(results gjson.Result, bodyString string) gjson.Result {
	re := regexp.MustCompile(`^([\d\.]+)\s?(GiB|MiB|TiB)$`)

	uploadRegexResult := re.FindStringSubmatch(results.Get("uploaded").Str)
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	edited_results, _ := sjson.Set(bodyString, "uploaded", helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2]))
	downloadRegexResult := re.FindStringSubmatch(results.Get("downloaded").Str)
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	edited_results, _ = sjson.Set(edited_results, "downloaded", helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2]))
	bufferRegexResult := re.FindStringSubmatch(results.Get("buffer").Str)
	cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
	edited_results, _ = sjson.Set(edited_results, "buffer", helpers.AnyUnitToBytes(cleanBuffer, downloadRegexResult[2]))

	return gjson.Parse(edited_results)
}
