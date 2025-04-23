package controllers

import "github.com/phachon/mm-wiki/app/services"

type SearchController struct {
	BaseController
}

func (this *SearchController) UpdateAllIndex() {

	// spaceId := strings.TrimSpace(this.GetString("space_id", ""))
	services.DocIndexService.UpdateAllDocIndex(30)
	this.jsonSuccess("添更新索引成功！", nil)

}

func (this *SearchController) RebuildIndex() {
	services.DocIndexService.DelAllDocIndex()
	services.DocIndexService.UpdateAllDocIndex(30)
	this.jsonSuccess("重建索引成功！", nil)

}
