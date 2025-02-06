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
	userPath := getUserPathGazelleScrape(baseUrl, cookieStr)
	req, _ := http.NewRequest("GET", baseUrl+userPath, nil)
	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req)

	return req
}

func getUserPathGazelleScrape(baseUrl string, cookieStr string) string {
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
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		href, found := doc.Find("body > div:nth-of-type(1) > div:nth-of-type(1) > div:nth-of-type(2) > ul > li:nth-of-type(10) > a").Attr("href")
		if !found {
			log.Fatal("Couldn't find userid")
		}
		return href
	} else {
		return ""
	}
}

func ProcessIndexerResponseGazelleScrape(bodyString string, indexerConfig gjson.Result) map[string]interface{} {
	//todo: handle cookie refresh
	results := map[string]interface{}{}
	re := regexp.MustCompile(`([\d\.]+)[ \x{00a0}]?\s?(GiB|MiB|TiB|KiB|B)`)
	extractNumberRegex := regexp.MustCompile(`\d+`)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	uploadRegexResult := re.FindStringSubmatch(doc.Find(indexerConfig.Get("scraping.xpaths.uploaded_amount").Str).Text())
	if len(uploadRegexResult) == 0 {
		fmt.Printf("An error occured while parsing GazelleScrape's response")
		return results
	}
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	results["uploaded_amount"] = helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2])

	downloadRegexResult := re.FindStringSubmatch(doc.Find(indexerConfig.Get("scraping.xpaths.downloaded_amount").Str).Text())
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	results["downloaded_amount"] = helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2])

	// bufferRegexResult := re.FindStringSubmatch(doc.Find(indexerConfig.Get("scraping.xpaths.buffer").Str).Text())
	// cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
	// results["buffer"] = helpers.AnyUnitToBytes(cleanBuffer, downloadRegexResult[2])

	seedingSizeRegexResult := re.FindStringSubmatch(doc.Find(indexerConfig.Get("scraping.xpaths.seeding_size").Str).Text())
	cleanSeedingSize, _ := strconv.ParseFloat(seedingSizeRegexResult[1], 64)
	results["seeding_size"] = helpers.AnyUnitToBytes(cleanSeedingSize, seedingSizeRegexResult[2])

	bonusPoints := doc.Find(indexerConfig.Get("scraping.xpaths.bonus_points").Str).Text()
	results["bonus_points"] = strings.ReplaceAll(bonusPoints, ",", "")

	uploaded_torrents := extractNumberRegex.FindString(doc.Find(indexerConfig.Get("scraping.xpaths.uploaded_torrents").Str).Text())
	results["uploaded_torrents"] = uploaded_torrents

	// snatched, seeding and leeching are retrieved via undocumented api : https://anthelion.me/ajax.php?action=community_stats&userid=
	// snatched := doc.Find(indexerConfig.Get("scraping.xpaths.snatched").Str).Text()
	// results["snatched"] = snatched

	// seeding := doc.Find(indexerConfig.Get("scraping.xpaths.seeding").Str).Text()
	// results["seeding"] = seeding

	// leeching := doc.Find(indexerConfig.Get("scraping.xpaths.leeching").Str).Text()
	// results["leeching"] = leeching

	ratio := doc.Find(indexerConfig.Get("scraping.xpaths.ratio").Str).Text()
	results["ratio"] = ratio

	torrent_comments := extractNumberRegex.FindString(doc.Find(indexerConfig.Get("scraping.xpaths.torrent_comments").Str).Text())
	results["torrent_comments"] = torrent_comments

	forum_posts := extractNumberRegex.FindString(doc.Find(indexerConfig.Get("scraping.xpaths.forum_posts").Str).Text())
	results["forum_posts"] = forum_posts

	freeleech_tokens := extractNumberRegex.FindString(doc.Find(indexerConfig.Get("scraping.xpaths.freeleech_tokens").Str).Text())
	results["freeleech_tokens"] = freeleech_tokens

	// warned, _ := strconv.Atoi(doc.Find(indexerConfig.Get("scraping.xpaths.warned").Str).Text())
	// results["warned"] = warned > 0

	// fmt.Println(results)

	return results
}
