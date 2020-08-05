package model

import (
	"fmt"
	"gopkg.in/check.v1"
	"jccAPITest/common"
)

type Goods struct {
	MerchName   string `json:"merch_name"`
	Title       string `json:"title"`
	SpuSkuCode  string `json:"spu_sku_code"`
	Supplier    string `json:"supplier"`
	ProductArea string `json:"product_area"`
	Brand       string `json:"brand"`
	CateOne     string `json:"cate_one"`
	CateTwo     string `json:"cate_two"`
	CateThree   string `json:"cate_three"`
	PageNo      string `json:"page_no"`
	PageSize    string `json:"page_size"`
	Order       string `json:"order"`
	OrderType   string `json:"order_type"`
	Status      string `json:"status"`
	FreightId   string `json:"freight_id"`
	common.BaseJccAPITesting
}

func init() {
	common.RegisterCheck(new(Goods))
}

func (g *Goods) UrlPath() string {
	return "/manager/platform_goods/goods_list"
}

func (g *Goods) SetUpSuite(c *check.C) {
	fmt.Println("-----------SetUpSuiteSetUpSuite----" + fmt.Sprint() + "---------------")
	l := new(Login)
	l.TestLoginSuccess(c)
	g.Token = l.Response.JSON().Object().Raw()["token"].(string)
	common.SetHeaderToken(g.Token)
}

func (g *Goods) TestSuccess(c *check.C) {
	g.Status = "2"
	res := common.HttpPost(c, g)
	ob := res.JSON().Object()
	fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("msg").Equal("操作成功")
}
