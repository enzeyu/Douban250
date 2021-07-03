package main

import (
	"DoubanSpirder/parse"
)

func main(){
	base_url := "https://movie.douban.com/top250"
	urls, Movies:= parse.GetPages(base_url)
	for _, j := range urls{
		if j.Url != ""{ // 去除首页爬取2次
			_, t := parse.GetPages(base_url+j.Url)
			for _,j :=range t{
				Movies = append(Movies,j)
			}
		}

	}
}
