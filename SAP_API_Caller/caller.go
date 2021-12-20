package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-business-partner-reads/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"
	"golang.org/x/xerrors"
)

type SAPAPICaller struct {
	baseURL string
	apiKey  string
	log     *logger.Logger
}

func NewSAPAPICaller(baseUrl string, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL: baseUrl,
		apiKey:  GetApiKey(),
		log:     l,
	}
}

func (c *SAPAPICaller) AsyncGetBP(businessPartner, businessPartnerRole, addressID string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Role":
			func() {
				c.Role(businessPartner, businessPartnerRole)
				wg.Done()
			}()
		case "Address":
			func() {
				c.Address(businessPartner, addressID)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) Role(businessPartner, businessPartnerRole string) {
	data, err := c.callBPSrvAPIRequirementRole("A_BusinessPartnerRole", businessPartner, businessPartnerRole)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callBPSrvAPIRequirementRole(api, businessPartner, businessPartnerRole string) (*sap_api_output_formatter.Role, error) {
	url := strings.Join([]string{c.baseURL, "API_BUSINESS_PARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithRole(req, businessPartner, businessPartnerRole)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToRole(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) Address(businessPartner, addressID string) {
	data, err := c.callBPSrvAPIRequirementAddress("A_BusinessPartnerAddress", businessPartner, addressID)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callBPSrvAPIRequirementAddress(api, businessPartner, addressID string) (*sap_api_output_formatter.Address, error) {
	url := strings.Join([]string{c.baseURL, "API_BUSINESS_PARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithAddress(req, businessPartner, addressID)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToAddress(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithRole(req *http.Request, businessPartner, businessPartnerRole string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s' and BusinessPartnerRole eq '%s'", businessPartner, businessPartnerRole))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithAddress(req *http.Request, businessPartner, addressID string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s' and AddressID eq '%s'", businessPartner, addressID))
	req.URL.RawQuery = params.Encode()
}
