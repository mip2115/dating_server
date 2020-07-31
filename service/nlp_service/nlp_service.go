package nlp_service

import (
	"fmt"

	"code.mine/dating_server/mapping"
	"github.com/bbalet/stopwords"

	stemmer "github.com/agonopol/go-stem"
	"github.com/fluhus/gostuff/nlp/wordnet"
	"github.com/karan/vocabulary"
)

// "github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/comprehend"

// 	So basically consider, before sending to the model ,removing stop words and low freq words.
// Also consider diving into ngrams and getting sentiment analysis on everything.
// You can also do word freq
// You can also try to get the lemma or stem of the word using the rule
// SSES -> SS
// IES -> I
// SS -> SS
// S -> S

// RemoveStopWords from text
func RemoveStopWords(text *string) (*string, error) {
	cleanedString := stopwords.CleanString(mapping.StrToV(text), "en", false)
	return mapping.StrToPtr(cleanedString), nil
}

func GetSynset() {
	wn, _ := wordnet.Parse("../../dictionary/dict")
	catNouns := wn.Search("elephant")["n"]
	// = slice of all synsets that contain the word "cat" and are nouns.
	fmt.Println(catNouns)
}

func GetStemOfWord() string {
	word := []byte(("ponies"))
	res := string(stemmer.Stem(word))
	return res
}

func GetWordSimilarity(a string, b string) (*float64, error) {
	wn, err := wordnet.Parse("../../dictionary/dict")
	if err != nil {
		return nil, err
	}
	var score float64
	synsets1 := wn.Search(a)["n"]
	synsets2 := wn.Search(b)["n"]

	for _, v := range synsets1 {
		for _, v2 := range synsets2 {
			similarity := wn.PathSimilarity(v, v2, true)
			if similarity > score {
				score = similarity
			}
		}
	}

	//similarity := wn.PathSimilarity(cat, dog, true)
	return &score, nil
}

func GetSynonyms(word *string) ([]string, error) {
	c := &vocabulary.Config{BigHugeLabsApiKey: "", WordnikApiKey: ""}

	// Instantiate a Vocabulary object with your config
	v, err := vocabulary.New(c)
	if err != nil {
		return nil, err
	}

	wordInformation, err := v.Word(mapping.StrToV(word))
	if err != nil {
		return nil, err
	}
	return wordInformation.Synonyms, nil
}
