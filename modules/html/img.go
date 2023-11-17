package _html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
	"github.com/licat233/goss/utils"
)

func (o *Html) handlerImgTag(doc *goquery.Document, htmlFilePath string) (hasModify bool) {
	// 遍历并处理img标签
	imgsE := doc.Find("img")
	length := imgsE.Length()
	// bar := progressbar.Default(int64(length)+1, utils.GetFileName(htmlFilePath))
	utils.Message("正在处理: %s的img标签", utils.GetFileName(htmlFilePath))
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
		newSrc, err := o.Bucket.UploadToOss(src)
		if err != nil {
			utils.Warning("upload image faild: %s", err)
			return
		}

		// 更新img标签的src属性
		s.SetAttr(attrName, newSrc)
		if !hasModify {
			hasModify = true
		}
	}
	imgsE.Each(func(i int, s *goquery.Selection) {
		handler(s, "src")
		handler(s, "data-src")
		bar.Increment()
	})
	bar.Increment()
	bar.Finish()
	return hasModify
}
