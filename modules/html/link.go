package _html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
	"github.com/licat233/goss/utils"
)

func (o *Html) handlerLinkTag(doc *goquery.Document, htmlFilePath string) (hasModify bool) {
	// 遍历并处理link标签
	linksE := doc.Find("link")
	length := linksE.Length()
	// bar := progressbar.Default(int64(length)+1, utils.GetFileName(htmlFilePath))
	utils.Message("正在处理: %s的link标签", utils.GetFileName(htmlFilePath))
	bar := pb.StartNew(length + 1)
	var handler = func(s *goquery.Selection, attrName string) {
		src, exists := s.Attr(attrName)
		if !exists {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}
		newSrc := src
		newSrc, err := o.uploadToOss(src)
		if err != nil {
			utils.Warning("upload image faild: %s", err)
			return
		}

		// 更新link标签的href属性
		s.SetAttr(attrName, newSrc)
		if !hasModify {
			hasModify = true
		}
	}
	linksE.Each(func(i int, s *goquery.Selection) {
		handler(s, "href")
		bar.Increment()
	})
	bar.Increment()
	bar.Finish()
	return hasModify
}
