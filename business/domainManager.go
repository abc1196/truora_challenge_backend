package business

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"../data"
	"../utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/likexian/whois-go"
)

// GetDomains - Get domains with their servers ordered by id
func GetDomains() []data.Domain {
	var domains []data.Domain
	utils.DB.Preload("Servers").Order("id").Find(&domains)
	return domains
}

// CreateDomain - Returns a domain given the id (hostname)
func CreateDomain(id string) (data.Domain, *utils.Exception) {
	// Make HTTP GET request
	response, err := http.Get("https://" + id)
	var persistedDomain data.Domain
	var dbQuery = utils.DB.Where("ID = ?", id).First(&persistedDomain)
	if dbQuery.RecordNotFound() {
		if err != nil {
			return persistedDomain, utils.ThrowException(400, "Host not Found")
		}
		logo, logoErr := getLogo(id, response)
		if logoErr != nil {
			return persistedDomain, logoErr
		}
		servers, serverErr := getServers(id)
		if serverErr != nil {
			return persistedDomain, serverErr
		}
		var title = getTitle(id, response)
		var sslGrade = getSslGrade(servers)

		persistedDomain.ID = id
		persistedDomain.SslGrade = sslGrade
		persistedDomain.Title = title
		persistedDomain.Servers = servers
		persistedDomain.Logo = logo
		persistedDomain.PreviousSslGrade = sslGrade
		persistedDomain.ServersChanged = false
		persistedDomain.IsDown = false
		persistedDomain.Created = time.Now().UTC()

		utils.DB.NewRecord(persistedDomain) // => returns `true` as primary key is blank

		utils.DB.Create(&persistedDomain)

		utils.DB.NewRecord(persistedDomain) // => return `false` after `user` created
	} else {
		if err != nil {
			persistedDomain.IsDown = true
			persistedDomain.Updated = time.Now().UTC()
			utils.DB.Save(&persistedDomain)
		} else {
			var persistedServers []data.Server
			utils.DB.Where("domain_refer = ?", persistedDomain.ID).Find(&persistedServers)
			servers, serverErr := getServers(id)
			if serverErr != nil {
				return persistedDomain, serverErr
			}
			var sslGrade = getSslGrade(servers)
			var now = time.Now().UTC()
			var diff = now.Sub(persistedDomain.Updated)
			if diff.Hours() >= 1 {
				var serversChanged = checkServers(persistedServers, servers)
				persistedDomain.ServersChanged = serversChanged
				persistedDomain.PreviousSslGrade = persistedDomain.SslGrade
			}
			persistedDomain.Servers = servers
			persistedDomain.SslGrade = sslGrade
			persistedDomain.Updated = time.Now().UTC()
			utils.DB.Save(&persistedDomain)
		}
	}
	return persistedDomain, nil
}

// getHost - Gets the api ssl labs host information
func getHost(id string) (data.Host, *utils.Exception) {
	var host data.Host
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + id)
	if err != nil {
		return host, utils.ThrowException(500, "Internal Server Error: APILabs")
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return host, utils.ThrowException(500, "Internal Server Error: APILabs")
		}
		jsonError := json.Unmarshal(bodyBytes, &host)
		if jsonError != nil {
			return host, utils.ThrowException(500, "Internal Server Error: APILabs")
		}
	}

	return host, nil
}

// getServers - Creates the servers given the id (hostname) and the api ssl labs endpoints
func getServers(id string) ([]data.Server, *utils.Exception) {
	host, err := getHost(id)
	if err != nil {
		return nil, err
	}
	var servers []data.Server
	for _, endpoint := range host.Endpoints {
		var server data.Server
		server.Address = endpoint.IPAddress
		server.SslGrade = endpoint.Grade
		endpointInfo, err := whois.Whois(endpoint.IPAddress)
		if err == nil {
			scanner := bufio.NewScanner(strings.NewReader(endpointInfo))
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "OrgName") {
					var orgNameArray = strings.Split(scanner.Text(), ":")
					var orgName = strings.TrimSpace(orgNameArray[1])
					server.Owner = orgName
				}
				if strings.Contains(scanner.Text(), "Country") {
					var countryArray = strings.Split(scanner.Text(), ":")
					var country = strings.TrimSpace(countryArray[1])
					server.Country = country
				}
			}
		}
		server.DomainRefer = id
		servers = append(servers, server)
	}
	return servers, nil
}

// getSslGrade - Returns the servers minimun SSL Grade
func getSslGrade(servers []data.Server) string {
	var min = 7
	for _, server := range servers {
		var gradeNumber = sslGradeToNumber(server.SslGrade)
		if gradeNumber < min {
			min = gradeNumber
		}
	}
	return numberToSsl(min)
}

// sslGradeToNumber - Maps the grade to a number
func sslGradeToNumber(grade string) int {
	var number = -1
	switch grade {
	case "A+":
		number = 6
	case "A":
		number = 5
	case "B":
		number = 4
	case "C":
		number = 3
	case "D":
		number = 2
	case "E":
		number = 1
	case "F":
		number = 0
	}
	return number
}

// numberToSsl - Maps the number to a SSL Grade
func numberToSsl(number int) string {
	var grade = ""
	switch number {
	case 6:
		grade = "A+"
	case 5:
		grade = "A"
	case 4:
		grade = "B"
	case 3:
		grade = "C"
	case 2:
		grade = "D"
	case 1:
		grade = "E"
	case 0:
		grade = "F"
	}
	return grade
}

// getTitle - Returns the id's (hostname) title in their HTML
func getTitle(id string, response *http.Response) string {
	var pageTitle = ""
	// Get the response body as a string
	dataInBytes, _ := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)
	// Find a substr
	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		titleStartIndex = strings.Index(pageContent, "<title id="+`"`+"pageTitle"+`"`+">")
		if titleStartIndex != -1 {
			// The start index of the title is the index of the first
			// character, the < symbol. We don't want to include
			// <title> as part of the final value, so let's offset
			// the index by the number of characers in <title>
			titleStartIndex += 22

		}
	} else {
		// The start index of the title is the index of the first
		// character, the < symbol. We don't want to include
		// <title> as part of the final value, so let's offset
		// the index by the number of characers in <title>
		titleStartIndex += 7
	}

	// Find the index of the closing tag
	titleEndIndex := strings.Index(pageContent, "</title>")

	// (Optional)
	// Copy the substring in to a separate variable so the
	// variables with the full document data can be garbage collected
	if titleEndIndex != -1 && titleStartIndex != -1 {
		pageTitle = string([]byte(pageContent[titleStartIndex:titleEndIndex]))
	} else {
		pageTitle = "No title element found"

	}

	return pageTitle
}

// getLogo - Returns the id's (hostname) logo in their HTML
func getLogo(id string, response *http.Response) (string, *utils.Exception) {
	var logo string

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return logo, utils.ThrowException(500, "Internal Server Error: Logo")
	}
	// Find and print image URLs
	document.Find("head").Each(func(index int, element *goquery.Selection) {
		element.Find("link").Each(func(i int, el *goquery.Selection) {
			value, exists := el.Attr("type")
			if exists {
				if strings.Contains(value, "image/x-icon") {
					l, e := el.Attr("href")
					if e {
						logo = l
					}
				}
			} else {
				value, exists := el.Attr("rel")
				if exists {
					if strings.Contains(value, "shortcut icon") || strings.Contains(value, "icon") {
						l, e := el.Attr("href")
						if e {
							logo = l
						}
					}
				}
			}
		})
	})
	return logo, nil
}

// checkServers - Compares the two given servers arrays
func checkServers(persistedServers []data.Server, servers []data.Server) bool {
	return cmp.Equal(persistedServers, servers)
}
