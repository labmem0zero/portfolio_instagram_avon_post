package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var ctx context.Context
var parentCtx context.Context
var	cookieDir string

type Order struct {
	Order struct {
		Orderid   string `json:"orderid"`
		Orderdate string `json:"orderdate"`
		Clientid  string `json:"clientid"`
		Goods     []struct {
			Goodid    string `json:"goodid"`
			Goodcount string `json:"goodcount"`
		} `json:"goods"`
	} `json:"order"`
}

type Client struct {
	Client struct {
		Clientid    string `json:"clientid"`
		Phonenumber string `json:"phonenumber"`
		Clientname  string `json:"clientname"`
		Firstorder  string `json:"firstorder"`
		Birthdate   string `json:"birthdate"`
	} `json:"client"`
}

type Good struct {
	Good struct {
		Goodid          string `json:"goodid"`
		Goodcategory    string `json:"goodcategory"`
		Goodname        string `json:"goodname"`
		Gooddescription string `json:"gooddescription"`
		Currentcat      string `json:"currentcat"`
		Currentcatpage  string `json:"currentcatpage"`
		Currentprice    string `json:"currentprice"`
		Goodurl			string `json:"goodurl"`
	} `json:"good"`
}

type good_mined struct {
	Data struct {
		Availability     int         `json:"Availability"`
		BadgeId          interface{} `json:"BadgeId"`
		BadgeText        string      `json:"BadgeText"`
		Brand            string      `json:"Brand"`
		CanAddToCart     bool        `json:"CanAddToCart"`
		CanAddToWishlist bool        `json:"CanAddToWishlist"`
		Categories       []struct {
			Dept *struct {
				CategoryType           int         `json:"CategoryType"`
				ChildProductCategories interface{} `json:"ChildProductCategories"`
				Description            interface{} `json:"Description"`
				Id                     int         `json:"Id"`
				ImageWidth             int         `json:"ImageWidth"`
				MenuLevel              int         `json:"MenuLevel"`
				Name                   string      `json:"Name"`
				NameMobile             interface{} `json:"NameMobile"`
				NotificationTicket     interface{} `json:"NotificationTicket"`
				ProductCount           int         `json:"ProductCount"`
				SearchWords            interface{} `json:"SearchWords"`
				ShowInMenu             bool        `json:"ShowInMenu"`
				Url                    interface{} `json:"Url"`
			} `json:"Dept"`
			Level2 struct {
				CategoryType           int         `json:"CategoryType"`
				ChildProductCategories interface{} `json:"ChildProductCategories"`
				Description            interface{} `json:"Description"`
				Id                     int         `json:"Id"`
				ImageWidth             int         `json:"ImageWidth"`
				MenuLevel              int         `json:"MenuLevel"`
				Name                   string      `json:"Name"`
				NameMobile             interface{} `json:"NameMobile"`
				NotificationTicket     interface{} `json:"NotificationTicket"`
				ProductCount           int         `json:"ProductCount"`
				SearchWords            interface{} `json:"SearchWords"`
				ShowInMenu             bool        `json:"ShowInMenu"`
				Url                    interface{} `json:"Url"`
			} `json:"Level2"`
			PDept struct {
				CategoryType           int         `json:"CategoryType"`
				ChildProductCategories interface{} `json:"ChildProductCategories"`
				Description            interface{} `json:"Description"`
				Id                     int         `json:"Id"`
				ImageWidth             int         `json:"ImageWidth"`
				MenuLevel              int         `json:"MenuLevel"`
				Name                   string      `json:"Name"`
				NameMobile             interface{} `json:"NameMobile"`
				NotificationTicket     interface{} `json:"NotificationTicket"`
				ProductCount           int         `json:"ProductCount"`
				SearchWords            interface{} `json:"SearchWords"`
				ShowInMenu             bool        `json:"ShowInMenu"`
				Url                    interface{} `json:"Url"`
			} `json:"PDept"`
		} `json:"Categories"`
		Category                string        `json:"Category"`
		ContainerInformation    string        `json:"ContainerInformation"`
		CrossSell               interface{}   `json:"CrossSell"`
		CustomerAlsoBought      []interface{} `json:"CustomerAlsoBought"`
		DeliveryType            int           `json:"DeliveryType"`
		Description             string        `json:"Description"`
		HasActiveVariant        bool          `json:"HasActiveVariant"`
		Id                      int           `json:"Id"`
		Ingredients             interface{}   `json:"Ingredients"`
		IsOnSale                bool          `json:"IsOnSale"`
		IsShadeVariant          bool          `json:"IsShadeVariant"`
		IsSizeVariant           bool          `json:"IsSizeVariant"`
		ListPrice               float64       `json:"ListPrice"`
		MarketingLabel1         string        `json:"MarketingLabel1"`
		MarketingLabel2         string        `json:"MarketingLabel2"`
		MetaDescription         string        `json:"MetaDescription"`
		MetaTitle               string        `json:"MetaTitle"`
		Name                    string        `json:"Name"`
		NewEndDate              interface{}   `json:"NewEndDate"`
		NewInclude              bool          `json:"NewInclude"`
		NewStartDate            interface{}   `json:"NewStartDate"`
		NotificationTicket      interface{}   `json:"NotificationTicket"`
		Price                   float64       `json:"Price"`
		PricePerUnitInformation string        `json:"PricePerUnitInformation"`
		ProductId               int           `json:"ProductId"`
		ProfileNumber           string        `json:"ProfileNumber"`
		Promotions              []interface{} `json:"Promotions"`
		Rating                  interface{}   `json:"Rating"`
		RatingCount             interface{}   `json:"RatingCount"`
		SaleCaption             string        `json:"SaleCaption"`
		SalePrice               float64       `json:"SalePrice"`
		SearchWords             interface{}   `json:"SearchWords"`
		ShortDescription        string        `json:"ShortDescription"`
		ShortName               string        `json:"ShortName"`
		SingleVariantSku        string        `json:"SingleVariantSku"`
		UnitPriceDetails        struct {
			MeasureUnit string  `json:"MeasureUnit"`
			UnitPrice   float64 `json:"UnitPrice"`
		} `json:"UnitPriceDetails"`
		VariantGroups []struct {
			GroupName     string      `json:"GroupName"`
			GroupPriority interface{} `json:"GroupPriority"`
			Variants      []struct {
				Availability     int         `json:"Availability"`
				CanAddToCart     bool        `json:"CanAddToCart"`
				ColorHex         string      `json:"ColorHex"`
				DeliveryType     int         `json:"DeliveryType"`
				DisplayPriority  interface{} `json:"DisplayPriority"`
				Fsc              string      `json:"Fsc"`
				IsHazmat         bool        `json:"IsHazmat"`
				LineNumber       string      `json:"LineNumber"`
				Sku              string      `json:"Sku"`
				SkuCount         int         `json:"SkuCount"`
				TryItOnMakeupSku string      `json:"TryItOnMakeupSku"`
				Type             int         `json:"Type"`
				VariantValue     string      `json:"VariantValue"`
				VolumePerUnit    float64     `json:"VolumePerUnit"`
				WeightPerUnit    float64     `json:"WeightPerUnit"`
			} `json:"Variants"`
		} `json:"VariantGroups"`
	} `json:"Data"`
	ErrorCode        int         `json:"ErrorCode"`
	ErrorId          interface{} `json:"ErrorId"`
	ErrorMessage     interface{} `json:"ErrorMessage"`
	RedirectUrl      interface{} `json:"RedirectUrl"`
	ValidationErrors interface{} `json:"ValidationErrors"`
}

func SubsToString(in []string)string{
	res:=""
	for _,s:=range in{
		res=res+s+"{br}"
	}
	return res
}

func StringToSubs(s string)[]string {
	res := strings.Split(s, "{br}")
	return res
}

type avon struct {
	XMLname xml.Name `xml:"avon"`
	Cover cover `xml:"cover"`
	Pages []page `xml:"page"`
}

type cover struct {
	Limage limage`xml:"limage"`
	Mimage mimage`xml:"mimage"`
	Simage simage`xml:"simage"`
}

type limage struct {
	Url string `xml:"url,attr"`
}

type simage struct {
	Url string `xml:"url,attr"`
}

type mimage struct {
	Url string `xml:"url,attr"`
}

type page struct {
	XMLname xml.Name `xml:"page"`
	Limage limage`xml:"limage"`
	Mimage mimage`xml:"mimage"`
	Simage simage`xml:"simage"`
	Products []product `xml:"product"`
}

type product struct {
	XMLname xml.Name `xml:"product"`
	ProductId string `xml:"id,attr"`
	ProductName string `xml:"name,attr"`
	ProductPrice string `xml:"price,attr"`
	Subproducts []subproduct `xml:"subproduct"`
}

type subproduct struct {
	XMLname xml.Name `xml:"subproduct"`
	ProductId string `xml:"id,attr"`
	ProductName string `xml:"name,attr"`
	ProductPrice string `xml:"price,attr"`
}

func TranslateCategory(in string)string{
	cats:=map[string]string{
		"ACCESSORIES":"Аксесуары",
		"BODY":"Ср-ва для тела",
		"COLOR":"Макияж",
		"FACE":"Ср-ва для лица",
		"FOOTWEAR":"Обувь",
		"FRAGRANCE":"Парфюмерия",
		"HAIR CARE":"Ср-ва для душа и волос",
		"HOME DECOR":"Домашний декор",
		"HOUSEWARES":"Для дома",
		"INNERWEAR":"Белье",
		"JEWELRY":"Ювелирные изд. и бижутерия",
		"KIDS HOME":"Для дома",
		"OUTERWEAR":"Одежда",
		"PERSONAL CARE":"Инструменты для ухода",
		"TOILETRIES":"Ср-ва для душа и волос",
		"WATCHES":"Аксесуары",
		"GIGIENIC":"Ср-ва личной гигиены",
		"UNCAT":"Без категории",
	}
	return cats[in]
}

func CategoryWithName(name string, cat string)string{
	if len(name)<5{
		return TranslateCategory("UNCAT")
	}
	res:=cat
	hair:=[]string{
		"краска для волос",
	}
	color:=[]string{
		"тени","пудра","карандаш",
	}
	gigi:=[]string{
		"интимн","мыло",
	}
	parf:=[]string{
		"туалетн","парфюм", "спрей для тела","духи",
	}
	for _,s:=range hair{
		if strings.Contains(strings.ToLower(name),s){
			res="HAIR CARE"
		}
	}
	for _,s:=range color{
		if strings.Contains(strings.ToLower(name),s){
			res="COLOR"
		}
	}
	for _,s:=range gigi{
		if strings.Contains(strings.ToLower(name),s){
			res="GIGIENIC"
		}
	}
	for _,s:=range parf{
		if strings.Contains(strings.ToLower(name),s){
			res="FRAGRANCE"
		}
	}
	return TranslateCategory(res)
}

type InstaPost struct{
	imgurl string
	imgdesc string
}

func ShowCategory(cat string){
	goodids:=DBGetCategoryGoods(cat)
	for _,g:=range goodids{
		ShowGood(g)
	}
}

func ShowGood(goodid string){
	good:=DBGetGood(goodid)
	fmt.Printf("--------------------------------\nКод товара: %v\nКатегория: %v\nНазвание товара: %v\nОписание товара: %v\nКаталог: %v\nСтраница: %v\nЦена: %v\nАдрес на сайте: %v\n", good.Good.Goodid,good.Good.Goodcategory,good.Good.Goodname,good.Good.Gooddescription,good.Good.Currentcat,good.Good.Currentcatpage,good.Good.Currentprice,good.Good.Goodurl)
}

func main(){
	//DBStart() инициализирует БД
	db=DBStart()
	cookieDir = "cookie"
	if _, err := os.Stat(cookieDir); os.IsNotExist(err) {
		err := os.Mkdir(cookieDir, os.ModePerm)
		if err != nil {
			fmt.Println("Ошибка при срздании папки с куки: ",err)
		}
	}
	cookieDir, _ = filepath.Abs(cookieDir)
	options := []chromedp.ExecAllocatorOption{
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("user-data-dir", cookieDir),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel1 := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	err:= chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	if err != nil {
		fmt.Println("При создании конекста возникла ошибка: ", err)
		return
	}
	secondCtx, cancel2 := context.WithCancel(chromeCtx)
	ctx = secondCtx
	parentCtx=chromeCtx

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		_ = <-sigChan
		fmt.Println("Закрываем конексты")
		cancel2()
		cancel1()
	}()
	defer cancel2()
	defer cancel1()

	//AVONGetGoods получает товары текущего каталога и добавляет их в бд, для использования раскомментировать
	//Важно: при первом запуске необходимо выполнить эту функцию, что бы заполнить БД товарами
	//AVONGetGoods()

	//AVONGetCatImages скачивает изображения всех страниц действующего каталога
	//AVONGetCatImages()

	//goodIDs присваем список товаров, разделенных пробелом
	//login и password - наши логин и пароль для входа в инстаграм. Можно ввести только один раз и сохранить в cookie
	//InstaStart(goddIDs string, login string, password string) запускает функцию постинга товаров в инстаграм,
	//в процессе скачивает изображение товара с сайта AVON во временный файл, после поста файл удаляется
	//для использования раскомментировать все 4 строчки
	//goddIDs:="1386829 1401001 1420694 94679 22375 1387176 1313782 1358737 1379088 49003 1457171 26640 1419764 14434 1441347 68782 20801 1409747 67604 1327226 1444822 1370547 1400305 1404669"
	//login:="instagram.login"
	//password:="instagram.password"
	//InstaStart(goddIDs, login, password)

	cancel2()
	cancel1()

	//AvonStart(login string, password string) недописанная функция для размещения заказа, пусть будет
	//avonlogin:="login"
	//avonpassword:="password"
	//AvonStart(avonlogin,avonpassword)

	//DBShowByCategoryCount() показывает количество товаров по категориям в БД, для использования раскомментировать
	//DBShowByCategoryCount()

	//ShowCategory(category string) показывает все товары введенной категории в БД, для использования раскомментировать
	//ShowCategory("Ср-ва личной гигиены")

	//ShowGood(goodid string) показывает отдельный товар по коду из БД, для использования раскомментировать
	//ShowGood("1421442")
	time.Sleep(10*time.Minute)
}