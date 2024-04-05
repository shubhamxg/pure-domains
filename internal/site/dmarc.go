package site

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/shubhamxg/pure-domains/pkg"
)

type DmarcResp struct {
	Total   string   `json:"total"`
	Hash    string   `json:"hash"`
	Domains []string `json:"domains"`
}

func Dmarc(domain string) []string {
	dmarc_req := pkg.Request(fmt.Sprintf("https://dmarc.live/info/%s", domain))
	defer dmarc_req.Body.Close()
	dmarc_bytes, _ := io.ReadAll(dmarc_req.Body)
	dmarc_hash := pkg.Parse(string(dmarc_bytes), `dmarc_hash:'`, `',`)

	dmarc_api_req := pkg.Request(fmt.Sprintf("https://dmarc.live/api/related/%s/1000000", dmarc_hash))
	defer dmarc_api_req.Body.Close()
	dmarc_api_bytes, _ := io.ReadAll(dmarc_api_req.Body)
	dmarc_resp := DmarcResp{}
	_ = json.Unmarshal(dmarc_api_bytes, &dmarc_resp)
	// fmt.Println(dmarc_resp.Domains)

	return dmarc_resp.Domains
}
