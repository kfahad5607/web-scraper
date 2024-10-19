package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type Link struct {
	scheme string
	domain string
	url    string
}

func main() {
	visitedLinks := map[string]bool{}
	deadLinks := []string{}

	page_url := "https://scrape-me.dreamsofcode.io/"
	visitLinks(page_url, true, &visitedLinks, &deadLinks)

	if len(deadLinks) > 0 {
		fmt.Println("Found following dead links:")
		for _, link := range deadLinks {
			fmt.Println(link)
		}
	}
}

func visitLinks(page_url string, recursive bool,  visitedLinks *map[string]bool, deadLinks *[]string){
	parsed_url, err := url.Parse(page_url)
	if err != nil{
		fmt.Println(err)
		return
	}

	fmt.Printf("visiting  %v\n", page_url)
	scheme := parsed_url.Scheme
	domain := parsed_url.Host

	resp, err := http.Get(page_url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if 400 <= resp.StatusCode && resp.StatusCode < 600 {
		*deadLinks = append(*deadLinks, page_url)
		return
	}

	if !recursive {
		return
	}
	
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	links := []Link{}
	findATag(doc, &links)

	for _, link := range links{
		if link.domain == "" {
			link.domain = domain

			if link.scheme == "" {
				link.scheme = scheme
			}

			link.url, err = url.JoinPath(link.scheme + "://", link.domain, link.url)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if _, ok := (*visitedLinks)[link.url]; ok {
			continue
		}
		(*visitedLinks)[link.url] = true

		visitLinks(link.url, link.domain == domain, visitedLinks, deadLinks)
	}
}

func findATag(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				parsed_url, err := url.Parse(attr.Val)
				if err != nil {
					fmt.Println(err)
					break
				}
				curr_url := attr.Val
				domain := parsed_url.Host
				*links = append(*links, Link{ scheme: parsed_url.Scheme, url: curr_url, domain: domain })
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findATag(c, links)
	}
}
