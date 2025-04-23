package fullsearch

import (
	"fmt"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/huichen/sego"
)

type SegoTokenizer struct {
	segmenter *sego.Segmenter
}

func NewSegoTokenizer(dictPath string) *SegoTokenizer {
	var segmenter sego.Segmenter
	segmenter.LoadDictionary(dictPath) // 加载词典文件（如未指定，使用内置词典）
	// 分词
	text := []byte("中华人民共和国中央人民政府")
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToString(segments, false))
	// fmt.Println(sego.SegmentsToString(segments, true))
	return &SegoTokenizer{segmenter: &segmenter}
}

// Tokenize 实现 bleve 的 Tokenizer 接口
func (t *SegoTokenizer) Tokenize(input []byte) analysis.TokenStream {
	text := string(input)
	segments := t.segmenter.InternalSegment([]byte(text), false)
	// segments := t.segmenter.Segment([]byte(text))
	tokens := make(analysis.TokenStream, 0, len(segments))
	for _, segment := range segments {
		token := &analysis.Token{
			Term:     []byte(segment.Token().Text()),
			Start:    segment.Start(),
			End:      segment.End(),
			Position: len(tokens) + 1,
			Type:     analysis.AlphaNumeric,
		}
		tokens = append(tokens, token)
	}
	return tokens
}

const SegoAnalyzerName = "sego_analyzer"
const SegoTokenizerName = "sego_tokenizer"

func init() {

	registry.RegisterTokenizer(SegoTokenizerName, func(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
		dictPath := "docs/search_dict/custom.txt,docs/search_dict/dictionary.txt" // Sego 词典路径
		return NewSegoTokenizer(dictPath), nil
	})

	// 注册自定义分析器
	registry.RegisterAnalyzer(SegoAnalyzerName, func(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
		tokenizer, err := cache.TokenizerNamed(SegoTokenizerName)
		if err != nil {
			return nil, err
		}
		stopFilter, err := cache.TokenFilterNamed(StopFilterName)
		if err != nil {
			return nil, err
		}
		// 创建自定义分析器
		analyzer := &analysis.DefaultAnalyzer{
			Tokenizer: tokenizer,
			TokenFilters: []analysis.TokenFilter{
				stopFilter,
			},
		}
		return analyzer, nil
	})

}
