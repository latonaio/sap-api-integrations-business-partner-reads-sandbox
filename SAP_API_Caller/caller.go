package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-business-partner-reads-rmq-kube/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"
	"golang.org/x/xerrors"
)

type RMQOutputter interface {
	Send(sendQueue string, payload map[string]interface{}) error
}

type SAPAPICaller struct {
	baseURL      string
	apiKey       string
	outputQueues []string
	outputter    RMQOutputter
	log          *logger.Logger
}

func NewSAPAPICaller(baseUrl string, outputQueueTo []string, outputter RMQOutputter, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:      baseUrl,
		apiKey:       GetApiKey(),
		outputQueues: outputQueueTo,
		outputter:    outputter,
		log:          l,
	}
}

func (c *SAPAPICaller) AsyncGetBP(businessPartner, businessPartnerRole, addressID, bankCountryKey, bankNumber, bPName string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "General":
			func() {
				c.General(businessPartner)
				wg.Done()
			}()
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
		case "Bank":
			func() {
				c.Bank(businessPartner, bankCountryKey, bankNumber)
				wg.Done()
			}()
		case "BPName":
			func() {
				c.BPName(bPName)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) General(businessPartner string) {
	generalData, err := c.callBPSrvAPIRequirementGeneral("A_BusinessPartner", businessPartner)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": generalData, "function": "BusinessPartnerGeneral"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(generalData)

	roleData, err := c.callToRole(generalData[0].ToRole)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": roleData, "function": "BusinessPartnerRole"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(roleData)

	addressData, err := c.callToAddress(generalData[0].ToAddress)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": addressData, "function": "BusinessPartnerAddress"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(addressData)

	bankData, err := c.callToBank(generalData[0].ToBank)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": bankData, "function": "BusinessPartnerBank"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(bankData)

}

func (c *SAPAPICaller) callBPSrvAPIRequirementGeneral(api, businessPartner string) ([]sap_api_output_formatter.General, error) {
	url := strings.Join([]string{c.baseURL, "API_BUSINESS_PARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithGeneral(req, businessPartner)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToGeneral(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToRole(url string) ([]sap_api_output_formatter.ToRole, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToRole(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToAddress(url string) ([]sap_api_output_formatter.ToAddress, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToAddress(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToBank(url string) ([]sap_api_output_formatter.ToBank, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToBank(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) Role(businessPartner, businessPartnerRole string) {
	roleData, err := c.callBPSrvAPIRequirementRole("A_BusinessPartnerRole", businessPartner, businessPartnerRole)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": roleData, "function": "BusinessPartnerRole"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(roleData)

}

func (c *SAPAPICaller) callBPSrvAPIRequirementRole(api, businessPartner, businessPartnerRole string) ([]sap_api_output_formatter.Role, error) {
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
	addressData, err := c.callBPSrvAPIRequirementAddress("A_BusinessPartnerAddress", businessPartner, addressID)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": addressData, "function": "BusinessPartnerAddress"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(addressData)

}

func (c *SAPAPICaller) callBPSrvAPIRequirementAddress(api, businessPartner, addressID string) ([]sap_api_output_formatter.Address, error) {
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

func (c *SAPAPICaller) Bank(businessPartner, bankCountryKey, bankNumber string) {
	bankData, err := c.callBPSrvAPIRequirementBank("A_BusinessPartnerBank", businessPartner, bankCountryKey, bankNumber)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": bankData, "function": "BusinessPartnerBank"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(bankData)

}

func (c *SAPAPICaller) callBPSrvAPIRequirementBank(api, businessPartner, bankCountryKey, bankNumber string) ([]sap_api_output_formatter.Bank, error) {
	url := strings.Join([]string{c.baseURL, "API_BUSINESS_PARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithBank(req, businessPartner, bankCountryKey, bankNumber)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToBank(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) BPName(bPName string) {
	bPNameData, err := c.callBPSrvAPIRequirementBPName("A_BusinessPartner", bPName)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": bPNameData, "function": "BusinessPartnerGeneral"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(bPNameData)

}

func (c *SAPAPICaller) callBPSrvAPIRequirementBPName(api, bPName string) ([]sap_api_output_formatter.General, error) {
	url := strings.Join([]string{c.baseURL, "API_BUSINESS_PARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithBPName(req, bPName)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToGeneral(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithGeneral(req *http.Request, businessPartner string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s'", businessPartner))
	req.URL.RawQuery = params.Encode()
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

func (c *SAPAPICaller) getQueryWithBank(req *http.Request, businessPartner, bankCountryKey, bankNumber string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s' and BankCountryKey eq '%s' and BankNumber eq '%s", businessPartner, bankCountryKey, bankNumber))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithBPName(req *http.Request, bPName string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("substringof('%s', BusinessPartnerName)", bPName))
	req.URL.RawQuery = params.Encode()
}
