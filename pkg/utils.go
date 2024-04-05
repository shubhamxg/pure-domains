package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func Request(url string) *http.Response {
	var (
		jar     = tls_client.NewCookieJar()
		options = []tls_client.HttpClientOption{
			tls_client.WithTimeoutSeconds(30),
			tls_client.WithClientProfile(profiles.Chrome_116_PSK),
			tls_client.WithNotFollowRedirects(),
			tls_client.WithCookieJar(jar),
		}
		client, client_err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	)

	if client_err != nil {
		log.Println(client_err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header = http.Header{
		"accept":          {"*/*"},
		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"user-agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	return resp
}

//	func Parse(data string, i int, left string, right string) [][]string {
//		r, err := regexp.Compile(fmt.Sprintf("%s(.*?)%s", left, right))
//		if err != nil {
//			log.Println("Failed to compile the regex")
//		}
//		matches := r.FindAllStringSubmatch(data, i)
//		return matches
//	}
func Parse(data string, left string, right string) string {
	r, err := regexp.Compile(fmt.Sprintf("%s(.*?)%s", left, right))
	if err != nil {
		log.Println("Failed to compile the regex")
	}
	matches := r.FindStringSubmatch(data)
	return matches[1]
}

// var dedupe_filtered_links = make(map[string]bool)

type Filtered struct {
	Unfiltered []string
	Filtered   []string
	Domains    []string
}

// func Filter(unfiltered_links_store []string) Filtered {

// 	var deduped_domains = make([]string, 0)
// 	deduped_filtered_links := make([]string, 0)
// 	deduped_unfiltered_links := make([]string, 0)
// 	for _, unfiltered_link := range unfiltered_links_store {
// 		domain := domains(unfiltered_link)
// 		if !dedupe_filtered_links[domain] {
// 			dedupe_filtered_links[domain] = true
// 			decoded_url_1 := strings.ReplaceAll(unfiltered_link, `\\u003d`, "=")
// 			decoded_url := strings.ReplaceAll(decoded_url_1, `\\u0026`, "&")

// 			deduped_unfiltered_links = append(deduped_unfiltered_links, decoded_url)
// 			deduped_domains = append(deduped_domains, domain)
// 			if strings.Contains(decoded_url, "=") {
// 				deduped_filtered_links = append(deduped_filtered_links, decoded_url)
// 			}
// 		}
// 	}

// 	return Filtered{
// 		Unfiltered: deduped_unfiltered_links,
// 		Filtered:   deduped_filtered_links,
// 		Domains:    deduped_domains,
// 	}
// }

// func domains(data string) string {
// 	r, err := regexp.Compile(`^(?:http:\/\/|www\.|https:\/\/)([^\/]+)`)
// 	if err != nil {
// 		log.Println("Failed to compile the regex")
// 	}
// 	matches := r.FindStringSubmatch(data)
// 	return matches[1]
// }

type SaveFile struct {
	mu         sync.Mutex
	folderName string
}

func Save() *SaveFile {
	current_time := fmt.Sprint(time.Now().Format("02-Jan-06 15:04:05"))
	current_time = strings.ReplaceAll(current_time, ":", "-")
	current_time = strings.ReplaceAll(current_time, " ", "_")
	foldername := fmt.Sprintf("Results/%s", current_time)
	err := os.MkdirAll(foldername, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return &SaveFile{
		folderName: foldername,
	}
}

func (s *SaveFile) Results(results Filtered, folder string) {
	// s.File(results.Filtered, "filtered_links.txt")
	// s.File(results.Unfiltered, "unfiltered_links.txt")
	s.File(results.Domains, folder)
	fmt.Printf("\n[~] Site > %s | Subdomains > %s", folder, len(results.Domains))

	// fmt.Printf("\nFiltered: %d\nUnfiltered: %d\nDomains: %d\n", len(results.Filtered), len(results.Unfiltered), len(results.Domains))
}

func (s *SaveFile) File(links []string, filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filename = fmt.Sprintf("%s/%s", s.folderName, filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range links {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return err
		}

		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (s *SaveFile) Test(links string, filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filename = fmt.Sprintf("%s/%s", s.folderName, filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// for _, line := range links {
	_, err = writer.WriteString(links + "\n")
	if err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}
	// }
	return nil
}
