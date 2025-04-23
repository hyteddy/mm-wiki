package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/phachon/mm-wiki/app/fullsearch"
	"github.com/phachon/mm-wiki/app/utils"
)

var (
	DocSearchBleve = NewDocSearchBleve()
)

func NewDocSearchBleve() *docBleve {
	return &docBleve{}
}

type docBleve struct {
	index map[string]interface{}
}

func (d *docBleve) Init() {

	docRootDir := utils.Document.GetAbsRootFileByPageFile("bleve-index")
	config := fullsearch.Config{
		IndexPath: docRootDir, // 空表示内存索引
		// DictPath:  "docs/search_dict/custom.txt,docs/search_dict/dictionary.txt",
	}

	// 初始化单例
	if err := fullsearch.InitSingleton(config); err != nil {
		panic(fmt.Sprintf("初始化失败: %v", err))
	}
}

func (d *docBleve) AddIndex(id string, data interface{}) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	index := fullsearch.GetIndex()

	index.Index(id, data)
	fmt.Println("添加索引成功")

	return nil
}

// 删除索引
func (d *docBleve) DelIndex(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	index := fullsearch.GetIndex()

	index.Delete(id)

	fmt.Println("删除索引成功")

	return nil
}

// 删除所有索引
func (d *docBleve) DelAllIndex() error {
	docRootDir := utils.Document.GetAbsRootFileByPageFile("bleve-index/store")
	os.RemoveAll(docRootDir)
	fmt.Println("删除索引成功")

	return nil
}

func (d *docBleve) DocSearch(params string) *bleve.SearchResult {

	index := fullsearch.GetIndex()

	// 构建查询
	// query := bleve.NewMatchQuery(params)
	query := bleve.NewMatchPhraseQuery(params)
	// query.SetField("name")
	// query.SetField("content")
	searchRequest := bleve.NewSearchRequest(query)
	// searchRequest.Fields = []string{"name", "content"}

	// 启用高亮并配置字段
	// 添加高亮字段（默认用 <em> 标签）
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Highlight.AddField("name")    // 高亮 content 字段\
	searchRequest.Highlight.AddField("content") // 高亮 content 字段\
	searchRequest.Size = 100

	// 执行搜索
	results, err := index.Search(searchRequest)
	if err != nil {
		return nil
	}
	fmt.Printf("共找到 %d 条结果\n", results.Total)
	// for _, hit := range results.Hits {
	// 	fmt.Printf("\nID: %s (得分: %.2f)\n", hit.ID, hit.Score)

	// 	// 输出字段
	// 	if id, ok := hit.Fields["id"].(string); ok {
	// 		fmt.Printf("标题: %s\n", id)
	// 	}

	// 	// 输出高亮
	// 	if fragments, ok := hit.Fragments["content"]; ok {
	// 		fmt.Printf("内容摘要: %s\n", strings.Join(fragments, "..."))
	// 	}
	// }
	return results
}
