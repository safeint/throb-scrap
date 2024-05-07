package main

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const path = `C:\Users\zhuji\Videos\vr\`
const imgPath = `C:\Users\zhuji\Videos\vr\actor\`

func main() {
	var info VideoInfo
	// 创建头像文件夹
	createDirectory(imgPath)

	mp4Files := GetMp4Files()
	for key, value := range mp4Files {
		var originName string
		key = strings.TrimSuffix(key, ".mp4")
		if strings.Contains(key, "-cd") {
			re := regexp.MustCompile(`-cd\d+`)
			originName = re.ReplaceAllString(key, "")
		} else {
			originName = key
		}
		// 读取页面
		client := &http.Client{}
		url := "https://www.javbus.com/" + originName
		req, _ := http.NewRequest("GET", url, nil)
		cookieStr := "PHPSESSID=ml2lu1upcj9br86dqm04bcmk87; existmag=mag; age=verified; dv=1"
		req.Header.Set("Cookie", cookieStr)
		resp, _ := client.Do(req)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		actorList := getActor(doc)
		var names []string
		for _, item := range actorList[:min(len(actorList), 5)] {
			names = append(names, item.Name)
		}

		// 创建文件夹
		// 已存在则将nfo文件复制一份给key
		var directoryName string
		if len(names) > 0 {
			directoryName = strings.Join(names, ",")
		} else {
			directoryName = "佚名"
		}
		// 携带番号的完整路径
		fullPath := path + directoryName + string(filepath.Separator) + originName + string(filepath.Separator)
		if createDirectory(fullPath) {
			// 未存在文件夹，说明所有相关图片为下载，走生成信息
			info.Actor = actorList
			info.Art = getArt(doc, fullPath, directoryName, originName)
			info.Country = "Japan"
			info.Cover = getCover(doc, fullPath, directoryName, originName)
			info.DateAdded = getBasicInfo(doc, "發行日期:")
			info.Director = getBasicInfo(doc, "導演:")
			info.EndDate = getBasicInfo(doc, "發行日期:")
			info.FanArt = info.Art.FanArt
			info.Genre = getTag(doc)
			info.LockData = false
			info.Mpaa = "JP-18+"
			info.Maker = getBasicInfo(doc, "製作商:")
			info.OriginalTitle = getOriginalTitle(doc)
			info.Poster = info.Art.Poster
			info.Premiered = getBasicInfo(doc, "發行日期:")
			runtimeMin, _ := strconv.Atoi(strings.ReplaceAll(getBasicInfo(doc, "長度"), "分鐘", ""))
			info.Runtime = runtimeMin
			info.SortName = info.OriginalTitle
			info.Studio = getBasicInfo(doc, "發行商:")
			info.Tag = getTag(doc)
			info.Thumb = info.Cover
			info.Title = getOriginalTitle(doc)
			avYear, _ := strconv.Atoi(strings.Split(info.DateAdded, "-")[0])
			info.Year = avYear
			parseXml(fullPath, key, info)
			parseXml(fullPath, originName, info)
		} else {
			// 已存在文件夹，说明相关信息已存在，复制一个nfo文件即可
			sourceFilePath := fullPath + originName + ".nfo"
			destinationFilePath := fullPath + key + ".nfo"

			sourceFile, _ := os.Open(sourceFilePath)
			defer func(sourceFile *os.File) {
				err := sourceFile.Close()
				if err != nil {

				}
			}(sourceFile)

			destinationFile, _ := os.Create(destinationFilePath)
			defer func(destinationFile *os.File) {
				err := destinationFile.Close()
				if err != nil {

				}
			}(destinationFile)

			_, err := io.Copy(destinationFile, sourceFile)
			if err != nil {
				return
			}
		}
		MoveMp4File(value, fullPath)

	}

}

func getActor(doc *goquery.Document) []Actor {
	var actor []Actor
	doc.Find(".avatar-box").Each(func(i int, selection *goquery.Selection) {
		imgSrc, _ := selection.Find("img").First().Attr("src")
		name := selection.Find("span").First().Text()
		info := Actor{Name: name, Role: name, Type: "Actor", SortOrder: i, Thumb: fmt.Sprintf("/config/metadata/actor/%s.jpg", name)}
		if !strings.Contains(imgSrc, "http") {
			imgName := fmt.Sprintf("%s.jpg", name)
			_, err := os.Stat(imgPath + imgName)
			if err != nil {
				imgDown(imgPath, name, "https://www.javbus.com"+imgSrc)
				time.Sleep(1 * time.Second)
			}
		}
		actor = append(actor, info)
	})
	return actor
}

func getArt(doc *goquery.Document, fullPath string, directoryName string, originName string) Art {
	var fanartImgs []map[string]string
	// 花絮信息
	doc.Find("#sample-waterfall a.sample-box").Each(func(i int, selection *goquery.Selection) {
		imgSrc, _ := selection.Attr("href")
		imgTitle, _ := selection.Find("img").Attr("title")
		if !strings.Contains(imgSrc, "http") {
			imgSrc = "https://www.javbus.com" + imgSrc
		}
		fanartImgs = append(fanartImgs, map[string]string{
			"src":   imgSrc,
			"title": imgTitle,
		})
	})
	// 保存海报
	fanart := make([]string, 0)
	imgSavePath := fullPath + string(filepath.Separator)
	// 保存花絮
	for _, img := range fanartImgs {
		imgDown(imgSavePath, img["title"], img["src"])
		fanart = append(fanart, "/media/"+directoryName+"/"+originName+"/"+img["title"]+".jpg")
	}

	poster := "/media/" + directoryName + "/" + originName + "/" + originName + ".jpg"
	return Art{Poster: poster, FanArt: fanart}
}

func getBasicInfo(doc *goquery.Document, tagStr string) string {
	selector := fmt.Sprintf("span.header:contains('%s')", tagStr)
	dateStr := doc.Find(selector).Parent().Text()
	return strings.TrimSpace(strings.TrimPrefix(dateStr, tagStr))
}

func getCover(doc *goquery.Document, fullPath string, directoryName string, originalName string) string {
	coverUrl, _ := doc.Find(".bigImage").First().Attr("href")
	imgDown(fullPath, originalName, "https://www.javbus.com"+coverUrl)
	return "/media/" + directoryName + "/" + originalName + "/" + originalName + ".jpg"
}

func getTitle(doc *goquery.Document) string {
	title, _ := doc.Find(".bigImage img").Attr("title")
	return title
}

func getOriginalTitle(doc *goquery.Document) string {
	title := doc.Find("h3").First().Text()
	return title
}

func getTag(doc *goquery.Document) []string {
	var tagList []string
	doc.Find("span.genre input[name='gr_sel']").Each(func(i int, selection *goquery.Selection) {
		text := selection.Next().Text()
		tagList = append(tagList, text)
	})
	return tagList
}

func imgDown(path string, name string, url string) {
	file, _ := os.Create(path + name + ".jpg")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	response, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return
	}

	_, err := io.Copy(file, response.Body)
	if err != nil {
		return
	}
}

func parseXml(path string, name string, info VideoInfo) {
	xmlData, _ := xml.MarshalIndent(info, "", "  ")
	file, _ := os.Create(path + name + ".nfo")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	// 手动添加 XML 声明
	_, err := file.WriteString(xml.Header)
	if err != nil {
		// 处理错误
		return
	}

	_, err = file.Write(xmlData)
	if err != nil {
		return
	}
}

func createDirectory(fullPath string) bool {
	_, err := os.Stat(fullPath)
	if err != nil {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}

func GetMp4Files() map[string]string {
	dirPath := path
	// 初始化一个用于存放文件名和路径的集合
	mp4Files := make(map[string]string)
	// 递归遍历目录及其子目录下的所有文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// 检查文件是否为目标文件（以 ".mp4" 为后缀）
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp4") {
			// 将文件名和完整路径存入集合
			mp4Files[info.Name()] = path
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	// 打印结果
	fmt.Println("MP4 Files:")
	for name, path := range mp4Files {
		fmt.Printf("%s: %s\n", name, path)
	}
	return mp4Files
}

func MoveMp4File(source string, destDir string) {
	sourceFileName := filepath.Base(source)
	destPath := filepath.Join(destDir, sourceFileName)
	err := os.Rename(source, destPath)
	if err != nil {
		fmt.Println(err)
	}
}
