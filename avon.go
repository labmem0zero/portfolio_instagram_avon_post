package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"context"
)

func AvonStart(){
	avonlogin:="login"
	avonpassword:="password"

	AVONLogin(avonlogin,avonpassword)
	AVONNewOrder()
	AVONLogOut()
}

func AVONGetPageImage(limage string, curcat string) string{
	imgurl:="https://www.avon.ru/REPSuite/static/all_brochures/ebrochure/ebrochureC02_22_01/ru_RU/{limage}"
	imgurl=strings.ReplaceAll(imgurl,"{limage}",limage)
	dir:="C:\\avon\\goods"+curcat
	fname:=strings.Split(limage,"/")[3]
	err:=os.MkdirAll(dir, 0777)
	if err!=nil{
		fmt.Println("При создании папки ",dir," произошла ошибка: ",err)
		return "error"
	}
	file, err:=os.Create(dir+"\\"+fname)
	if err!=nil{
		fmt.Println("При создании файла ",dir+"\\"+fname," произошла ошибка: ",err)
		return "error"
	}
	defer file.Close()
	resp,err:=http.Get(imgurl)
	if err!=nil{
		fmt.Println("При скачивании файла",imgurl," произошла ошибка: ",err)
		return "error"
	}
	_,err=io.Copy(file,resp.Body)
	if err!=nil{
		fmt.Println("При копировании изображения в файл",fname," произошла ошибка: ",err)
		return "error"
	}
	return dir+"\\"+fname
}

func AVONGetCatImages(){
	schema:="https://www.avon.ru/REPSuite/static/all_brochures/ebrochure/ebrochureC{month}_{year}_01/ru_RU/pages.xml"
	time:=time.Now()
	month:=strconv.Itoa(int(time.Month()))
	if len(month)<2{
		month="0"+month
	}
	year:=strconv.Itoa(int(time.Year()))
	year=year[2:]
	fmt.Printf("Month: %v Year: %v\n",month,year)
	schema=strings.ReplaceAll(schema,"{month}",month)
	schema=strings.ReplaceAll(schema,"{year}",year)
	fmt.Println("XML: ",schema)
	xmlResp,err:=http.Get(schema)
	if err!=nil{
		fmt.Println("При скачивании XML с товарами возникла ошибка:",err)
	}
	xmlContent, err:=ioutil.ReadAll(xmlResp.Body)
	if err!=nil{
		fmt.Println("При считывании тела XML с товарами возникла ошибка:",err)
	}
	xmlGoods:=avon{}
	err=xml.Unmarshal(xmlContent, &xmlGoods)
	if err!=nil{
		fmt.Println("При анмаршале XML с товарами возникла ошибка:",err)
	}
	AVONGetPageImage(xmlGoods.Cover.Limage.Url,month+year)
	for _,page:=range xmlGoods.Pages{
		AVONGetPageImage(page.Limage.Url,month+year)
	}
}

func AVONGetGoods(){
	schema:="https://www.avon.ru/REPSuite/static/all_brochures/ebrochure/ebrochureC{month}_{year}_01/ru_RU/pages.xml"
	time:=time.Now()
	month:=strconv.Itoa(int(time.Month()))
	if len(month)<2{
		month="0"+month
	}
	year:=strconv.Itoa(int(time.Year()))
	year=year[2:]
	fmt.Printf("Month: %v Year: %v\n",month,year)
	schema=strings.ReplaceAll(schema,"{month}",month)
	schema=strings.ReplaceAll(schema,"{year}",year)
	fmt.Println("XML: ",schema)
	xmlResp,err:=http.Get(schema)
	if err!=nil{
		fmt.Println("При скачивании XML с товарами возникла ошибка:",err)
	}
	xmlContent, err:=ioutil.ReadAll(xmlResp.Body)
	if err!=nil{
		fmt.Println("При считывании тела XML с товарами возникла ошибка:",err)
	}
	xmlGoods:=avon{}
	err=xml.Unmarshal(xmlContent, &xmlGoods)
	if err!=nil{
		fmt.Println("При анмаршале XML с товарами возникла ошибка:",err)
	}
	for inx,page:=range xmlGoods.Pages{
		pag:=strconv.Itoa((inx+1)*2)+"-"+strconv.Itoa((inx+1)*2+1)
		for _,product:=range page.Products{
			tmpGood:=Good{}
			mined:=AVONMine(product.ProductId)
			tmpGood.Good.Goodid=product.ProductId
			tmpGood.Good.Goodcategory=CategoryWithName(mined.Data.Name,mined.Data.Category)
			tmpGood.Good.Goodname=mined.Data.Name
			tmpGood.Good.Gooddescription=SubsToString(AVONMineDecription(mined))
			tmpGood.Good.Currentcat=month+"/"+year
			tmpGood.Good.Currentcatpage=pag
			tmpGood.Good.Currentprice=product.ProductPrice
			tmpGood.Good.Goodurl="https://my.avon.ru/tovar/"+strconv.Itoa(mined.Data.Id)
			DBInsertGood(tmpGood)
			//goods=append(goods, tmpGood)
			for _, subprod:=range product.Subproducts{
				tmpGood:=Good{}
				mined:=AVONMine(subprod.ProductId)
				tmpGood.Good.Goodid=subprod.ProductId
				tmpGood.Good.Goodcategory=CategoryWithName(mined.Data.Name,mined.Data.Category)
				tmpGood.Good.Goodname=mined.Data.Name
				tmpGood.Good.Gooddescription=SubsToString(AVONMineDecription(mined))
				tmpGood.Good.Currentcat=month+"/"+year
				tmpGood.Good.Currentcatpage=pag
				tmpGood.Good.Currentprice=subprod.ProductPrice
				tmpGood.Good.Goodurl="https://my.avon.ru/tovar/"+strconv.Itoa(mined.Data.Id)
				DBInsertGood(tmpGood)
			}
		}
	}
	DBShowByCategoryCount()
}

func AvonGetImageUrl(goodurl string)string{
	html:=""
	err:=chromedp.Run(ctx,
		chromedp.Navigate(goodurl),
		chromedp.Sleep(Cooldown(2000)),
		chromedp.OuterHTML(`document.querySelector('body')`,&html,chromedp.ByJSPath),
	)
	if err!=nil{
		fmt.Println("Ошибка при посещении страницы ",goodurl,":", err)
		return ""
	}
	gq,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		fmt.Println("Ошибка при создании goquery документа", err)
		return ""
	}
	img:=gq.Find(`#ProductMediaCarousel div.Slides > div:nth-child(1) Img.GalleryImage[src*='product']`).First()
	url,_:=img.Attr("src")
	return url
}

func AVONLogin(login string, password string){
	html:=""
	chromedp.Run(ctx,
		chromedp.Navigate("https://www.avon.ru/REPSuite/loginMain.page"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(Cooldown(10000))
			return nil
		}),
		chromedp.SendKeys(`input[id='loginuserid']`,login),
		chromedp.SendKeys(`input[id='loginpassword']`,password),
		chromedp.Click(`button[id='submitBtn']`),
		chromedp.OuterHTML(`document.querySelector("body")`,&html,chromedp.ByJSPath),
	)
	time.Sleep(Cooldown(5000))
}

func AVONLogOut(){
	err:=os.RemoveAll(cookieDir)
	if err!=nil{
		fmt.Println("Не смог удалить куки по причине:",err)
	}
}

func AVONNewOrder()  {
	chromedp.Run(ctx,
		chromedp.Navigate("https://www.avon.ru/REPSuite/home.page"),
		chromedp.ActionFunc(func( context.Context) error {
			time.Sleep(Cooldown(5000))
			return nil
		}),
		chromedp.Click(`button[data-button='Start_New_Order']`),
		chromedp.ActionFunc(func( context.Context) error {
			time.Sleep(Cooldown(5000))
			return nil
		}),
		chromedp.Reload(),
	)
	time.Sleep(Cooldown(5000))
}

func AVONMine(goodid string)good_mined{
	var goodmined good_mined
	data := []byte(`{"lineNumber":"`+goodid+`"}`)
	r := bytes.NewReader(data)
	resp, err := http.Post("https://my.avon.ru/Api/SearchApi/ProductIdSubmit", "application/json", r)
	if err != nil {
		fmt.Println(err)
	}
	bytes,_:=ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &goodmined)
	fmt.Println("Для товара с кодом ",goodid," SKU=",goodmined.Data.Id)
	return goodmined
}

func AVONMineDecription(good good_mined)[]string{
	p:=bluemonday.StripTagsPolicy()
	descArray:=strings.Split(good.Data.Description, "<br />")
	for i,s:=range descArray{
		descArray[i]=p.Sanitize(s)
	}
	return descArray
}
