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

func getSelector(n *html.Node) string {
	if n.Parent != nil {
		s := getSelector(n.Parent)
		ss := n.Data + ":nth-child(" + strconv.Itoa(nthChild(n)) + ")"
		if s != "" {
			return s + ">" + ss
		} else {
			return ss
		}
	}
	return ""
}

func Cooldown(near float64) time.Duration{
	zoom:=int(near/10)
	x:=rand.Intn(zoom)+int(0.95*near)
	return time.Duration(x)*time.Millisecond
}
