# tdxF10Protocol-goVer

This is the golang version of pytdx subject to limited implementations for personal usage, including 
* GetCodeName() from get_security_list(), which returns a code-name mapping of stock codes that start with 000,001,002,300,600,601,603 in corresponding sz and sh markets.
* GetCompanyInfoCategory() from get_company_info_category()
* GetCompanyInfoContent() from get_comapany_info_content()

**import "github.com/sidazhang123/tdxF10Protocol-goVer"**

```
	// this mod subjects to codes that start with 000, 001, 002, 300 in sz and, 600, 601, 603 in sh
	**api := tdxF10Protocol_goVer.Socket{}**
	defer api.Close()
	// get a code-name mapping
  // using built-in ip pool (ipPool=nil)
	err, codeNameMap := **api.GetCodeNameMap(nil)**
	if err != nil {
		println(err.Error())
	}
	_ = PrettyPrint(codeNameMap)
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
  
```
	//get fields to fetch the f10 content of given codes
	code := "600001"
	err, codeF10InfoMap := api.GetCompanyInfoCategory([]string{code}, nil)
	if err != nil {
		println(err.Error())
	}
	_ = PrettyPrint(codeF10InfoMap)
```
 "600001": {
    "业内点评": [
      "600001.txt", //filename
      "275723", //start_pos
      "4539"  //length
    ],
    "主力追踪": [
      "600001.txt",
      "225153",
      "10241"
    ],
    "公司大事": [
      "600001.txt",
      "115728",
      "79201"
    ],
    "公司报导": [
      "600001.txt",
      "74669",
      "22882"
      ...
```
	// get the specific content
	x := codeF10InfoMap[code]["经营分析"]
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
☆经营投资☆ ◇600001 邯郸钢铁 更新日期：2009-08-12◇ 港澳资讯 灵通V5.0
★本栏包括【1.主营构成】【2.经营投资】【3.关联企业经营状况】
【1.主营构成】
【主营构成】
【2009年中期概况】
┌────────────┬─────┬─────┬───┬──────┐
｜项目名称                ｜营业收入( ｜营业利润( ｜毛利率｜占主营业务收｜
｜                        ｜万元)     ｜万元)     ｜(%)   ｜入比例(%)   ｜
├────────────┼─────┼─────┼───┼──────┤
｜冶金(行业)              ｜ 946168.32｜ 107153.60｜ 11.33｜       61.28｜
├────────────┼─────┼─────┼───┼──────┤
｜板材(产品)              ｜ 611673.44｜  52597.79｜  8.60｜       39.62｜
｜棒材(产品)              ｜ 225067.84｜  35478.23｜ 15.76｜       14.58｜
├────────────┼─────┼─────┼───┼──────┤
｜华北地区(地区)          ｜ 511827.95｜         -｜     -｜       33.15｜
｜华东地区(地区)          ｜ 242969.05｜         -｜     -｜       15.74｜
｜中南地区(地区)          ｜  96296.62｜         -｜     -｜        6.24｜
｜西北地区(地区)          ｜  48702.27｜         -｜     -｜        3.15｜
｜西南地区(地区)          ｜  35086.98｜         -｜     -｜        2.27｜
...
