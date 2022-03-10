package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"context"
)

func InstaStart(in string){
	inS:=strings.Split(in, " ")
	login:="instagram.login"
	password:="instagram.password"
	InstagramCheckLogin(login,password)
	InstagramPostArray(InstaMakePostArray(inS))
}

func InstagramPostArray(posts []InstaPost) {
	for _,post:=range posts{
		InstagramPostGood(post)
		time.Sleep(Cooldown(600000))
	}
}

func InstagramPostGood(post InstaPost){
	file,err:=os.Create("tmpimage.jpg")
	if err!=nil{
		fmt.Println("При создании временного файла изображения возникла ошибка: ",err)
		return
	}
	imgpath,err:=filepath.Abs("tmpimage.jpg")
	if err!=nil{
		fmt.Println("При формировании абсолютного пути к файлу возникла ошибка: ",err)
		return
	}
	resp,err:=http.Get(post.imgurl)
	if err!=nil{
		fmt.Println("При скачивании изображения возникла ошибка: ",err)
		return
	}
	_,err=io.Copy(file,resp.Body)
	if err!=nil{
		fmt.Println("При копировании изображения во временный файл возникла ошибка: ",err)
		return
	}
	file.Close()
	err= chromedp.Run(ctx,
		chromedp.Navigate("https://instagram.com"),
		chromedp.Sleep(Cooldown(10000)),
		//наж. кнопки для доб. поста
		chromedp.Click(`nav div:nth-child(3)>div>button`),
		chromedp.Sleep(Cooldown(5000)),
		//доб. файла в форму
		chromedp.SetUploadFiles(`form input[multiple][type='file']`, []string{imgpath}, chromedp.ByQuery),
		chromedp.Sleep(Cooldown(5000)),
		//наж. кнопки далее
		chromedp.Click(`div[role='dialog'] div:nth-child(3)>div>button`),
		chromedp.Sleep(Cooldown(5000)),
		//наж. кнопки далее
		chromedp.Click(`div[role='dialog'] div:nth-child(3)>div>button`),
		chromedp.Sleep(Cooldown(5000)),
	)
	descStrings:=strings.Split(post.imgdesc,"{br}")
	for _,s:=range descStrings{
		chromedp.Run(ctx,
			chromedp.SendKeys(`textarea[aria-label*='подпись']`, s+kb.Enter, chromedp.ByQuery),
			chromedp.Sleep(Cooldown(300)),
		)
	}
	chromedp.Run(ctx,
		chromedp.Sleep(Cooldown(700)),
		chromedp.Click(`div[role='dialog'] div:nth-child(3)>div>button`),
	)
	os.Remove(imgpath)
}

func InstagramCheckLogin(login string, password string){
	page:=""
	err:=chromedp.Run(ctx,
		chromedp.Navigate("https://www.instagram.com/"),
		chromedp.Sleep(Cooldown(5000)),
		chromedp.OuterHTML(`document.querySelector("body")`,&page,chromedp.ByJSPath),
	)
	if err!=nil{
		fmt.Println("Ошибка при входе на сайт: ",err)
		return
	}
	if err!=nil{
		fmt.Println("Ошибка при получении html сайта: ",err)
		return
	}
	gq, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err!=nil{
		fmt.Println("Ошибка при cчитывании документа из ридера: ",err)
		return
	}
	logged:=true
	gq.Find("input").Each(func(_ int, inp *goquery.Selection) {
		attr, _:=inp.Attr("name")
		if attr=="password"{
			logged=false
			fmt.Println("Нужно залогиниться")
		}
	})
	if  logged==false {
		err = chromedp.Run(ctx,
			chromedp.SendKeys("#loginForm > div > div:nth-child(1) > div > label > input", login),
			chromedp.SendKeys("#loginForm > div > div:nth-child(2) > div > label > input", password),
			chromedp.Click("#loginForm > div > div:nth-child(3) > button"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				time.Sleep(Cooldown(5000))
				fmt.Println("Залогинены успешно")
				time.Sleep(Cooldown(5000))
				return nil
			}),
		)
	}else{
		fmt.Println("Уже залогинены")
	}
}

func InstagramPreparePost(goodid string)InstaPost{
	res:=InstaPost{}
	good:=DBGetGood(goodid)
	if good.Good.Currentprice==""{
		fmt.Printf("У товара с ID=%v нет цены, печалька\n",good.Good.Goodid)
		return res
	}
	if len(good.Good.Gooddescription)<10{
		fmt.Printf("У товара с ID=%v нет описания, печалька\n",good.Good.Goodid)
		return res
	}
	if len(good.Good.Goodname)<5{
		fmt.Printf("У товара с ID=%v нет названия, печалька\n",good.Good.Goodid)
		return res
	}
	reg:=regexp.MustCompile("&#..;")
	res.imgurl=AvonGetImageUrl(good.Good.Goodurl)
	tmpdesc:=strings.Split(good.Good.Gooddescription,"{br}")[0]
	res.imgdesc=good.Good.Goodname+"{br}{br}"+tmpdesc+"{br}{br}Код товара: "+goodid+"{br}Цeна в каталоге:"+good.Good.Currentprice+"р.{br}страница в каталоге:"+good.Good.Currentcatpage+"{br}{br}Информация о заказе и доставке по ссылке в шапке страницы"
	res.imgdesc=reg.ReplaceAllString(res.imgdesc,"")
	return res
}

func InstaMakePostArray(postarray []string)[]InstaPost  {
	res:=[]InstaPost{}
	for _,s:=range postarray{
		post:=InstagramPreparePost(s)
		if len(post.imgdesc)<15{
			continue
		}
		res=append(res,post)
	}
	return res
}
