package utils

import (
	"fmt"
	"github.com/remko/go-mkvparse"
	"os"
	"sort"
)

type mkvTagParser struct {
	mkvparse.DefaultHandler

	title *string
}

func (p *mkvTagParser) HandleString(id mkvparse.ElementID, value string, info mkvparse.ElementInfo) error {
	switch id {
	case mkvparse.TitleElement:
		p.title = &value
	}
	return nil
}

func FindAllTags(filePath string) map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
	defer file.Close()
	tagParser := mkvTagParser{}
	tagsh := mkvparse.NewTagsHandler()
	err = mkvparse.ParseSections(file, mkvparse.NewHandlerChain(&tagParser, tagsh), mkvparse.InfoElement, mkvparse.TagsElement)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}

	// Print (sorted) tags
	if tagParser.title != nil {
		fmt.Printf("- title: %q\n", *tagParser.title)
	}
	tags := tagsh.Tags()
	var tagNames []string
	for tagName := range tags {
		tagNames = append(tagNames, tagName)
	}
	sort.Strings(tagNames)
	for _, tagName := range tagNames {
		fmt.Printf("- %s: %q\n", tagName, tags[tagName])
	}
	return tags
}
