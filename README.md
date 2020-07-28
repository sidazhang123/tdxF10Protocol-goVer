# tdxF10Protocol-goVer

This is the golang version of pytdx subject to limited implementations for personal usage, including 
* GetCodeName(ipPool) from get_security_list(), which returns a code-name mapping of stock codes that start with 000,001,002,300,600,601,603 in corresponding sz and sh markets.
* GetCompanyInfoCategory(codeSlice, ipPool) from get_company_info_category() with a retry policy; it processes the codes collectively
* GetCompanyInfoContent(code, filename, start, length, ipPool) from get_comapany_info_content(); it processes a code individually

**import "github.com/sidazhang123/tdxF10Protocol-goVer"**

```
	// this mod subjects to codes that start with 000, 001, 002, 300 in sz and, 600, 601, 603 in sh
	// retry with each ip in the ipPool recursively; MaxRetry=len(Addrs) if 0/undefined
	api := tdxF10Protocol_goVer.Socket{
		MaxRetry: 0,
	}
    // api.Init(addrs,timeout)
	// Addrs affects GetCompanyInfoCategory & GetCompanyInfoContent; nil for default
	// timeout=5s by default(<=0)
	api.Init(nil,0)
	defer api.Close()
```

```
	// get code-name mappings
	// it uses a separate ip pool; nil for default
	err, codeNameMap := api.GetCodeNameMap(nil)
	if err != nil {
		println(err.Error())
	}
```
{

  "000001": "平安银行",
  
  "000002": "万科Ａ",
  
  "000004": "国农科技",
  
  "000005": "世纪星源",
  
  "000006": "深振业Ａ",
  
  "000007": "全新好",
  
  "000008": "神州高铁",
  
  "000009": "中国宝安",
  
  ...
  
  **codename len = 3775**

  **GetCodeNameMap took 7.0307578s**
  
```
    // get fields to fetch the f10 content of given codes
    // this method applies the retry policy
    err, category := api.GetCompanyInfoCategory(codes)
    if err != nil {
        fmt.Println(err.Error())
    }
    _ = PrettyPrint(category) 
```

 "600000": {
 
"业内点评": [

"600000.txt", //filename

"275723", //start_pos

"4539"  //length

],

"主力追踪": [

"600000.txt",

"225153",

"10241"

],

"公司大事": [

"600000.txt",

"115728",

"79201"

],

"公司报导": [

"600000.txt",

"74669",

"22882"

...

**category len = 3775**

**GetCompanyInfoCategory took 2m6.4742127s**
    
```
	// get the specific content
    // avoid being coupled with the return type of GetCompanyInfoCategory, you can impl a retry mech likewise by yourself.
	code := "600000"
    x := category[code]["经营分析"]
	filename := x[0]
	start := x[1]
	length := x[2]
	err, s := api.GetCompanyInfoContent(code, filename, start, length, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(s)

}
```

☆经营分析☆ ◇600000 浦发银行 更新日期：2020-07-16◇ 港澳资讯 灵通V7.0
★本栏包括【1.主营业务】【2.主营构成分析】【3.经营投资】【4.关联企业经营状况】★

【1.主营业务】
吸收公众存款；发放短期、中期和长期贷款；办理结算；办理票据贴现；发行金融债券；代理发行、代理兑付、承销政府债券；买卖政府债券；同业拆借；提供信用证服务及担保；代理收付款项及代理保险业务；提供保险箱业务；外汇存款；外汇贷款；外汇汇款；外币兑换；国际结算；同业外汇拆借；外汇票据的承兑和贴现；外汇借款；外汇担保；结汇、售汇；买卖和代理买卖股票以外的外币有价证券；自营外汇买卖；代客外汇买卖；资信调查、咨询、见证业务；离岸银行业务；证券投资基金托管业务；全国社会保障基金托管业务；经中国人民银行和中国银行业监督管理委员会批准经营的其他业务。


【2.主营构成分析】
【截止日期】2019-12-31
项目名                        营业收入    营业利润   毛利率(%)  占主营业务

...

**GetCompanyInfoContent took 174.5871ms**
