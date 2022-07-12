package src

import (
	"fmt"
	"os"
	"sf/config"
	"sf/src/boluobao"
	"strconv"
)

type AutoGenerated struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Index    int    `json:"index"`
	IsVip    bool   `json:"is_vip"`
	VolumeID string `json:"volume_id"`
	Content  string `json:"content"`
}

func GetCatalogue(BookData Books) {
	response := boluobao.Get_catalogue_detailed_by_id(BookData.NovelID)
	var orderList []string
	for _, data := range response.Data.VolumeList {
		fmt.Println("start download volume: ", data.Title)
		for _, Chapter := range data.ChapterList {
			ChapId := strconv.Itoa(Chapter.ChapID)
			if Chapter.OriginNeedFireMoney > 0 {
				orderList = append(orderList, strconv.Itoa(Chapter.ChapID))
				continue
			}
			GetContent(len(data.ChapterList), BookData.NovelName, ChapId)

		}
	}
	if len(orderList) != 0 {
		fmt.Println(len(orderList), "is no need to download")
	}

	fmt.Println("NovelName:", BookData.NovelName, "download complete!")
}

func GetContent(ChapterLength int, NovelName, cid string) {
	response := boluobao.Get_content_detailed_by_cid(cid)
	if response.Status.HTTPCode != 200 {
		if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" {
			fmt.Println(response.Status.Msg)
			os.Exit(0)
		} else {
			fmt.Println(response.Status.Msg)
		}
	} else {
		if f, err := os.OpenFile(config.Var.SaveFile+"/"+NovelName+".txt",
			os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			defer func(f *os.File) {
				err = f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(f)
			if _, ok := f.WriteString("\n\n\n" + response.Data.Title + response.Data.Expand.Content); ok != nil {
				fmt.Println(ok)
			}
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf(
		"download Volume No:%d: %s : %d/%d  %v\r",
		response.Data.Sno, response.Data.Title, response.Data.ChapOrder, ChapterLength, response.Data.Title,
	)
}
