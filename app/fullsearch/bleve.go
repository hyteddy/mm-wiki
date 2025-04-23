package fullsearch

import (
	"bufio"
	"os"
	"sync"

	"github.com/blevesearch/bleve/v2" // ✅ 仅导入主包
	"github.com/blevesearch/bleve/v2/mapping"
)

var (
	instance bleve.Index
	once     sync.Once
)

// 配置结构体（使用公共接口）
type Config struct {
	IndexPath string
	DictPath  string
}

func readStopWords(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stopWords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopWords = append(stopWords, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopWords, nil
}

// 初始化单例
func InitSingleton(config Config) error {
	var err error
	once.Do(func() {
		// 1. 初始化中文分词器
		// var segmenter sego.Segmenter
		// segmenter.LoadDictionary(config.DictPath)
		// segmenter.
		// chineseStopWords, err := readStopWords("docs/search_dict/stop_tokens.txt")

		// stopFilter := tokenfilter.NewStopTokenFilter(chineseStopWords, analysis.Ideographic)
		// mapping := bleve.NewIndexMapping()
		// mapping := bleve.NewDocumentMapping()
		mapping, err := buildIndexMapping()
		if err != nil {
			panic(err)
		}

		// // 注册 Sego 分词器
		// err = mapping.AddCustomTokenMap("stop_word", map[string]interface{}{
		// 	"type":     tokenmap.Name,
		// 	"filename": "docs/search_dict/stop_tokens.txt",
		// })
		// if err != nil {
		// 	panic(err)
		// }
		// err = mapping.AddCustomTokenFilter("myStopFilter", map[string]interface{}{
		// 	"type":           stop.Name,
		// 	"stop_token_map": "stop_word",
		// })
		// err = mapping.AddCustomAnalyzer(SegoAnalyzerName, map[string]interface{}{
		// 	"type":          SegoAnalyzerName,
		// 	"token_filters": []string{"myStopFilter"},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = mapping.AddCustomTokenMap("cn_wordmap", map[string]interface{}{
		// 	"type": tokenmap.Name,
		// 	"tokens": []interface{}{
		// 		"我", "的",
		// 	},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = mapping.AddCustomTokenFilter("cn", map[string]interface{}{
		// 	"type":           "stop_tokens",
		// 	"stop_token_map": "cn_wordmap",
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err = mapping.AddCustomAnalyzer("china", map[string]interface{}{
		// 	"type":         SegoAnalyzerName,
		// 	"char_filters": []interface{}{},
		// 	"tokenizer":    SegoTokenizerName,
		// 	"token_filters": []interface{}{
		// 		`cn`,
		// 	},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// err := mapping.AddCustomTokenFilter("color_stop_filter", map[string]interface{}{
		// 	"type": stop_tokens_filter.Name,
		// 	"tokens": []interface{}{
		// 		"red",
		// 		"green",
		// 		"blue",
		// 	},
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// registry.RegisterAnalyzer("sego_word", func(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
		// 	dictPath, ok := config["dict_path"].(string)
		// 	if !ok {
		// 		return nil, fmt.Errorf("sego tokenizer requires dict_path configuration")
		// 	}
		// 	tokenizer := NewSegoTokenizer(dictPath)
		// 	// 创建自定义分析器
		// 	analyzer := &analysis.DefaultAnalyzer{
		// 		Tokenizer:    tokenizer,
		// 		TokenFilters: []analysis.TokenFilter{},
		// 	}
		// 	return analyzer, nil
		// })

		// err = mapping.AddCustomTokenFilter(StopFilterName, map[string]interface{}{"type": StopFilterName})
		// if err != nil {
		// 	panic(err)
		// }
		// err = mapping.AddCustomTokenizer(SegoTokenizerName, map[string]interface{}{"type": SegoTokenizerName})
		// if err != nil {
		// 	panic(err)
		// }

		// err = mapping.AddCustomAnalyzer(SegoAnalyzerName, map[string]interface{}{
		// 	"type":          SegoAnalyzerName,
		// 	"dict_path":     config.DictPath,
		// 	"tokenizer":     SegoTokenizerName,
		// 	"token_filters": []string{StopFilterName, "lowercase"},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// 配置文档类型的分析器
		// docMapping := bleve.NewDocumentMapping()
		// titleFieldMapping := bleve.NewTextFieldMapping()
		// titleFieldMapping.Analyzer = SegoAnalyzerName // 使用 Sego 分析器
		// docMapping.AddFieldMappingsAt("content", titleFieldMapping)
		// docMapping.AddFieldMappingsAt("name", titleFieldMapping)
		// mapping.AddDocumentMapping(mapping.DefaultType, docMapping)
		// // mapping.DefaultAnalyzer = SegoAnalyzerName

		// 配置字段使用自定义分析器

		// err := mapping.AddCustomAnalyzer("sego", map[string]interface{}{
		// 	"type": "sego",
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// tokenMap := analysis.NewTokenMap()
		// err = tokenMap.LoadFile("docs/search_dict/stop_tokens.txt")
		// if err != nil {
		// 	panic(err)
		// }
		// tokenMap.AddToken("stop_test")

		// mapping.DefaultAnalyzer = "myStopFilter"

		// cache := registry.NewCache()
		// _, err = cache.DefineTokenMap("stop_test", map[string]interface{}{
		// 	"type":   tokenmap.Name,
		// 	"tokens": []interface{}{"a", "in", "the"},
		// })

		// stopConfig := map[string]interface{}{
		// 	"type":           "stop_tokens",
		// 	"stop_token_map": "stop_test",
		// }

		// stopFilter, err := cache.DefineTokenFilter("stop_test", stopConfig)
		// if err != nil {
		// 	panic(err)
		// }

		// err := mapping.AddCustomAnalyzer("myAnalyzer", map[string]interface{}{
		// 	"type":      "sego",
		// 	"tokenizer": "sego",
		// 	"filters":   []string{"lowercase"},
		// 	// "token_filters": []string{lowercase.Name, "myStopFilter"},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// mapping.DefaultAnalyzer = "myAnalyzer"
		// 4. 创建或打开索引
		if config.IndexPath == "" {
			instance, err = bleve.NewMemOnly(mapping)
		} else {
			instance, err = bleve.Open(config.IndexPath)
			if err == bleve.ErrorIndexPathDoesNotExist {
				instance, err = bleve.New(config.IndexPath, mapping)
			}
		}
		if err != nil {
			panic(err)
		}

	})
	return err
}

func buildIndexMapping() (mapping.IndexMapping, error) {
	// a custom field definition that uses our custom analyzer
	// segoMapping := bleve.NewTextFieldMapping()
	// segoMapping.Analyzer = SegoAnalyzerName

	// breweryMapping := bleve.NewDocumentMapping()
	// breweryMapping.AddFieldMappingsAt("content", segoMapping)

	indexMapping := bleve.NewIndexMapping()
	// indexMapping.AddDocumentMapping("aa", breweryMapping)

	indexMapping.DefaultAnalyzer = "segoa"

	// var err = indexMapping.AddCustomTokenizer("segot", map[string]interface{}{
	// 	"type": SegoTokenizerName,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	var err = indexMapping.AddCustomAnalyzer("segoa",
		map[string]interface{}{
			"type":          SegoAnalyzerName,
			"tokenizer":     SegoTokenizerName,
			"token_filters": []string{},
		})
	if err != nil {
		return nil, err
	}
	// doc, err := models.DocumentModel.GetDocumentByDocumentId("4")
	// if err != nil {
	// 	panic(err)
	// }
	// content, _, err := models.DocumentModel.GetDocumentContentByDocument(doc)
	// if err != nil {
	// 	panic(err)
	// }
	// doc["content"] = content

	// tokenStream, _ := indexMapping.AnalyzeText(SegoAnalyzerName, []byte(content))
	// for _, token := range tokenStream {
	// 	fmt.Printf("[%s] ", token.Term)
	// }

	return indexMapping, nil
}

// 获取全局实例（需先调用InitSingleton）
func GetIndex() bleve.Index {
	return instance
}
