package fullsearch

import (
	"bufio"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

// 自定义停用词过滤器
type StopTokenFilter struct {
	stopWords map[string]struct{}
}

// 从文件加载停用词
func NewStopTokenFilter(path string) (*StopTokenFilter, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	count := 0
	stopWords := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		stopWords[word] = struct{}{}
		count++
	}
	fmt.Printf("成功加载停用词数量: %d\n", count) // 检查实际加载数量
	return &StopTokenFilter{stopWords: stopWords}, nil
}

// 实现过滤器接口
func (f *StopTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	output := make(analysis.TokenStream, 0)
	for _, token := range input {
		if _, ok := f.stopWords[string(token.Term)]; !ok {
			output = append(output, token)
		}
	}
	return output
}

// 注册到 Bleve
const StopFilterName = "my_stop_filter"

func init() {
	registry.RegisterTokenFilter(StopFilterName, func(config map[string]interface{}, cache *registry.Cache) (analysis.TokenFilter, error) {
		path := "docs/search_dict/stop_tokens.txt" // 替换为实际路径
		return NewStopTokenFilter(path)
	})
}
