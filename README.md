# sap-api-integrations-business-partner-reads
sap-api-integrations-business-partner-reads は、外部システム(特にエッジコンピューティング環境)をSAPと統合することを目的に、SAP API で BP(ビジネスパートナ) - 一般 データを取得するマイクロサービスです。    
sap-api-integrations-business-partner-reads には、サンプルのAPI Json フォーマットが含まれています。   
sap-api-integrations-business-partner-reads は、オンプレミス版である（＝クラウド版ではない）SAPS4HANA API の利用を前提としています。クラウド版APIを利用する場合は、ご注意ください。   
https://api.sap.com/api/OP_API_BUSINESS_PARTNER_SRV/overview   

## 動作環境  
sap-api-integrations-business-partner-reads は、主にエッジコンピューティング環境における動作にフォーカスしています。  
使用する際は、事前に下記の通り エッジコンピューティングの動作環境（推奨/必須）を用意してください。  
・ エッジ Kubernetes （推奨）    
・ AION のリソース （推奨)    
・ OS: LinuxOS （必須）    
・ CPU: ARM/AMD/Intel（いずれか必須）    

## クラウド環境での利用
sap-api-integrations-business-partner-reads は、外部システムがクラウド環境である場合にSAPと統合するときにおいても、利用可能なように設計されています。  

## 本レポジトリ が 対応する API サービス
sap-api-integrations-business-partner-reads が対応する APIサービス は、次のものです。

* APIサービス概要説明 URL: https://api.sap.com/api/OP_API_BUSINESS_PARTNER_SRV/overview    
* APIサービス名(=baseURL): API_BUSINESS_PARTNER

## 本レポジトリ に 含まれる API名
sap-api-integrations-business-partner-reads には、次の API をコールするためのリソースが含まれています。  

* A_BusinessPartner（ビジネスパートナ - 一般）※ビジネスパートナの詳細データを取得するために、ToRole、ToAddress、ToBankと合わせて利用されます。  
* ToRole（ビジネスパートナ - ロール）
* ToAddress（ビジネスパートナ - アドレス）
* ToBank（ビジネスパートナ - 銀行）
* A_BusinessPartnerRole（ビジネスパートナ - ロール）
* A_BusinessPartnerAddress（ビジネスパートナ - アドレス）
* A_BusinessPartnerBank（ビジネスパートナ - 銀行）

## API への 値入力条件 の 初期値
sap-api-integrations-business-partner-readsにおいて、API への値入力条件の初期値は、入力ファイルレイアウトの種別毎に、次の通りとなっています。  

### SDC レイアウト

* inoutSDC.BusinessPartner.BusinessPartner（ビジネスパートナ）
* inoutSDC.BusinessPartner.Role.BusinessPartnerRole（ビジネスパートナロール）
* inoutSDC.BusinessPartner.Address.AddressID（アドレスID）
* inoutSDC.BusinessPartner.Bank.BankCountryKey（銀行国コード）
* inoutSDC.BusinessPartner.Bank.BankNumber（銀行コード）
* inoutSDC.BusinessPartner.BusinessPartnerName（ビジネスパートナ名）

## SAP API Bussiness Hub の API の選択的コール

Latona および AION の SAP 関連リソースでは、Inputs フォルダ下の sample.json の accepter に取得したいデータの種別（＝APIの種別）を入力し、指定することができます。  
なお、同 accepter にAll(もしくは空白)の値を入力することで、全データ（＝全APIの種別）をまとめて取得することができます。  

* sample.jsonの記載例(1)  

accepter において 下記の例のように、データの種別（＝APIの種別）を指定します。  
ここでは、"Role" が指定されています。    
  
```
	"api_schema": "sap.s4.beh.businesspartner.v1.BusinessPartner.Created.v1",
	"accepter": ["Role"],
	"business_partner_code": "1000140",
	"deleted": false
```
  
* 全データを取得する際のsample.jsonの記載例(2)  

全データを取得する場合、sample.json は以下のように記載します。  

```
	"api_schema": "sap.s4.beh.businesspartner.v1.BusinessPartner.Created.v1",
	"accepter": ["All"],
	"business_partner_code": "1000140",
	"deleted": false
```

## 指定されたデータ種別のコール

accepter における データ種別 の指定に基づいて SAP_API_Caller 内の caller.go で API がコールされます。  
caller.go の func() 毎 の 以下の箇所が、指定された API をコールするソースコードです。  

```
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
```
## Output  
本マイクロサービスでは、[golang-logging-library-for-sap](https://github.com/latonaio/golang-logging-library-for-sap) により、以下のようなデータがJSON形式で出力されます。  
以下の sample.json の例は、ビジネスパートナ の ロール が取得された結果の JSON の例です。  
以下の項目のうち、"BusinessPartner" ～ "ValidTo" は、/SAP_API_Output_Formatter/type.go 内 の Type Role {} による出力結果です。"cursor" ～ "time"は、golang-logging-library-for-sap による 定型フォーマットの出力結果です。  

```
{
	"cursor": "/Users/latona2/bitbucket/sap-api-integrations-business-partner-reads/SAP_API_Caller/caller.go#L179",
	"function": "sap-api-integrations-business-partner-reads/SAP_API_Caller.(*SAPAPICaller).Role",
	"level": "INFO",
	"message": [
		{
			"BusinessPartner": "1000140",
			"BusinessPartnerRole": "FLCU01",
			"ValidFrom": "",
			"ValidTo": ""
		}
	],
	"time": "2022-01-27T17:47:10+09:00"
}
```
