package indexers

import (
	"backend/database"
	"backend/helpers"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func LoginAndGetCookiesTL(username string, password string, twoFaCode string, loginURL string, indexerInfo gjson.Result) string {

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
	if !strings.Contains(cookiesStr, "tlpass") {
		// login failed
		return ""
	}
	cookiesStr = cookiesStr[:len(cookiesStr)-1]
	return cookiesStr
}

func ConstructRequestTL(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str

	cookieStr := database.GetIndexerCookies(indexerId)
	// indexer is not setup yet, credentials are in prowlarr but session cookie needs to be stored in dasharr's db
	// TODO: move this execution to container initialization
	if cookieStr == "" {
		LoginAndSaveCookies(indexerName, "", "", "", "", indexerId)
		cookieStr = database.GetIndexerCookies(indexerId)
	}
	username := prowlarrIndexerConfig.Get("extraFieldData.username").Str
	userPageUrl := fmt.Sprintf("%sprofile/%s", baseUrl, username)
	req, _ := http.NewRequest("GET", userPageUrl, nil)
	req = addCookiesToRequest(req, cookieStr)
	// fmt.Println(req)

	return req
}

func ProcessIndexerResponseTL(bodyString string, indexerInfo gjson.Result) map[string]interface{} {
	//todo: handle cookie refresh
	results := map[string]interface{}{}
	re := regexp.MustCompile(`([\d\.]+)[ \x{00a0}]?\s?(GB|MB|TB|KB|B)`)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	// fmt.Println(bodyString)
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

	bonusPoints := doc.Find(indexerInfo.Get("scraping.xpaths.bonus_points").Str).Text()
	results["bonus_points"] = strings.ReplaceAll(bonusPoints, "'", "")

	seeding_html := doc.Find(indexerInfo.Get("scraping.xpaths.seeding").Str).Text()
	results["seeding"] = regexp.MustCompile(`\((\d+)\)`).FindStringSubmatch(seeding_html)[1]

	leeching_html := doc.Find(indexerInfo.Get("scraping.xpaths.leeching").Str).Text()
	results["leeching"] = regexp.MustCompile(`\((\d+)\)`).FindStringSubmatch(leeching_html)[1]

	ratio_html := doc.Find(indexerInfo.Get("scraping.xpaths.ratio").Str).Text()
	results["ratio"] = regexp.MustCompile(`(\d+)`).FindStringSubmatch(ratio_html)[1]

	torrent_comments := doc.Find(indexerInfo.Get("scraping.xpaths.torrent_comments").Str).Text()
	results["torrent_comments"] = torrent_comments

	class := doc.Find(indexerInfo.Get("scraping.xpaths.class").Str).Text()
	results["class"] = class

	return results
}
