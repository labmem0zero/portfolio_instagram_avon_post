package main

import (
	"golang.org/x/net/html"
	"math/rand"
	"strconv"
	"time"
)

func nthChild(n *html.Node) int {
	i := 1
	for n.PrevSibling != nil {
		n = n.PrevSibling
		if n.Type == html.ElementNode {
			i++
		}
	}
	return i
}

//getSelector(n *html.Node) возвращает querySelector заданной html ноды
func getSelector(n *html.Node) string {
	if n.Parent != nil {
		s := getSelector(n.Parent)
		//nthChild(n *html.Node) возвращает n-ый номер текущей ноды среди всех сиблингов ноды
		ss := n.Data + ":nth-child(" + strconv.Itoa(nthChild(n)) + ")"
		if s != "" {
			return s + ">" + ss
		} else {
			return ss
		}
	}
	return ""
}

//Cooldown(near float64) time.Duration возвращает временной интервал с +-5% погрешностью, необходима для имитации человеческих действий
func Cooldown(near float64) time.Duration{
	zoom:=int(near/10)
	x:=rand.Intn(zoom)+int(0.95*near)
	return time.Duration(x)*time.Millisecond
}
