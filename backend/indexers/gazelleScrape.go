package indexers

import (
	"backend/database"
	"backend/helpers"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func LoginAndGetCookiesGazelleScrape(username string, password string, twoFaCode string, loginURL string, indexerInfo gjson.Result) string {

	// body := indexerInfo.Get("login.body").String()
	fields := indexerInfo.Get("login.fields").Map()
	formData := url.Values{}
	formData.Add(fields["username"].String(), username)
	formData.Add(fields["password"].String(), password)
	formData.Add(fields["twoFaCode"].String(), twoFaCode)

	extraFields := fields["extra"].Map()
	for key, val := range extraFields {
		formData.Add(key, val.String())
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevents redirect
		},
	}
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
	if !strings.Contains(cookiesStr, "session") {
		// login failed
		return ""
	}
	cookiesStr = cookiesStr[:len(cookiesStr)-1]
	return cookiesStr
}

func ConstructRequestGazelleScrape(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str

	cookieStr := database.GetIndexerCookies(indexerId)
	userPath := getUserPathGazelleScrape(baseUrl, cookieStr, indexerName)
	req, _ := http.NewRequest("GET", baseUrl+userPath, nil)
	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req)

	return req
}

func getUserPathGazelleScrape(baseUrl string, cookieStr string, indexerName string) string {
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
		indexerInfo := helpers.GetIndexerInfo(indexerName)
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		href, found := doc.Find(indexerInfo.Get("scraping.xpaths.user_path").Str).Attr("href")
		if !found {
			fmt.Println("Couldn't find userid")
			return ""
		}
		return href
	} else {
		return ""
	}
}

func ProcessIndexerResponseGazelleScrape(bodyString string, indexerInfo gjson.Result) map[string]interface{} {
	//todo: handle cookie refresh
	results := map[string]interface{}{}
	re := regexp.MustCompile(`([\d\.]+)[ \x{00a0}]?\s?(GiB|GB|MiB|MB|TiB|TB|KiB|KB|B)`)
	extractNumberRegex := regexp.MustCompile(`\d+`)
	extractNumberAfterColumnRegex := regexp.MustCompile(`: (\d+)`)
	indexerName := indexerInfo.Get("indexer_name").Str

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	uploadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_amount").Str).Text())
	if len(uploadRegexResult) == 0 {
		fmt.Printf("An error occured while parsing GazelleScrape's response")
		return results
	}
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	results["uploaded_amount"] = helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2])

	var bonusPoints string
	if indexerName == "AlphaRatio" {
		bonusPointsRegexResult := extractNumberAfterColumnRegex.FindStringSubmatch(strings.ReplaceAll(doc.Find(indexerInfo.Get("scraping.xpaths.bonus_points").Str).Text(), ",", ""))
		bonusPoints = bonusPointsRegexResult[1]
	} else {
		bonusPoints = doc.Find(indexerInfo.Get("scraping.xpaths.bonus_points").Str).Text()
	}
	results["bonus_points"] = strings.ReplaceAll(bonusPoints, ",", "")

	downloadRegexResult := re.FindStringSubmatch(doc.Find(indexerInfo.Get("scraping.xpaths.downloaded_amount").Str).Text())
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	results["downloaded_amount"] = helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2])

	seedingSizeXpath := indexerInfo.Get("scraping.xpaths.seeding_size")
	if seedingSizeXpath.Exists() {
		seedingSizeRegexResult := re.FindStringSubmatch(doc.Find(seedingSizeXpath.Str).Text())
		cleanSeedingSize, _ := strconv.ParseFloat(seedingSizeRegexResult[1], 64)
		results["seeding_size"] = helpers.AnyUnitToBytes(cleanSeedingSize, seedingSizeRegexResult[2])
	}

	uploaded_torrents := extractNumberRegex.FindString(doc.Find(indexerInfo.Get("scraping.xpaths.uploaded_torrents").Str).Text())
	results["uploaded_torrents"] = uploaded_torrents

	// snatched, seeding and leeching are retrieved via undocumented api : https://anthelion.me/ajax.php?action=community_stats&userid=
	// snatched := doc.Find(indexerInfo.Get("scraping.xpaths.snatched").Str).Text()
	// results["snatched"] = snatched

	// seeding := doc.Find(indexerInfo.Get("scraping.xpaths.seeding").Str).Text()
	// results["seeding"] = seeding

	// leeching := doc.Find(indexerInfo.Get("scraping.xpaths.leeching").Str).Text()
	// results["leeching"] = leeching

	ratioXpath := indexerInfo.Get("scraping.xpaths.seeding_size")
	if ratioXpath.Exists() {
		ratio := doc.Find(indexerInfo.Get(ratioXpath.Str).Str).Text()
		results["ratio"] = ratio
	}

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

	// fmt.Println(results)

	return results
}
