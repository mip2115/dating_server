package nlp_service

import (
	"fmt"
	"sort"
	"strings"

	"code.mine/dating_server/aws"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"github.com/bbalet/stopwords"

	stemmer "github.com/agonopol/go-stem"
	AMZN "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
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

// ProcessText –
// first always clean the text – also remove low count words
// get a score based on word count freq
// get a sentiment analysis
// get a score based on distance between distance of words (tirgams)
// gt a score based on keywords/phrases
// get a score based on entity extraction

// GetTextBreakDown –
func GetTextBreakDown(text string) (*types.TextSummary, error) {
	textSummary := &types.TextSummary{}
	text = RemoveStopWords(text)
	text = RemoveLowCountWords(text)
	// remove punctiation and stuff?
	//	text = GetStemsOfText(&text)
	// get sentiment analsys here

	// get word word rate score
	// so get the top 5 or so wrods w highest rate and see if htere's a scor for that
	topRatedWords := GetWordRateScore(text)
	textSummary.TopRatedWords = topRatedWords

	sess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}
	client := comprehend.New(sess)
	entities, err := GetEntitiesOfText(text, client)
	if err != nil {
		return nil, err
	}
	textSummary.Entities = entities

	keyphrases, err := GetKeyphrasesOfText(text, client)
	if err != nil {
		return nil, err
	}
	textSummary.Keyphrases = keyphrases
	// you have to do key phrase extraction before stemming
	// stem AFTER you have key phrase extraction
	// key phrase exctraction
	// you already cleaned the word
	// but use distance algo to detect that

	// then break into trigrams
	// do same thing here
	// ngrams := GetNGrams(text, 3)
	return textSummary, nil
}

// GetKeyphrasesOfText –
func GetKeyphrasesOfText(text string, client *comprehend.Comprehend) ([]*string, error) {
	input := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: AMZN.String("en"),
		Text:         AMZN.String(text),
	}
	result, err := client.DetectKeyPhrases(input)
	if err != nil {
		return nil, err
	}
	keyphrases := []*string{}
	for _, v := range result.KeyPhrases {
		keyphrases = append(keyphrases, mapping.StrToPtr(v.String()))
	}
	for i, v := range keyphrases {
		temp := strings.Split(*v, " ")
		for idx, val := range temp {
			temp[idx] = GetStemOfWord(val)
		}
		str := strings.Join(temp, " ")
		keyphrases[i] = &str
	}
	return keyphrases, nil
}

// GetEntitiesOfText –
func GetEntitiesOfText(text string, client *comprehend.Comprehend) ([]*string, error) {
	input := &comprehend.DetectEntitiesInput{
		LanguageCode: AMZN.String("en"),
		Text:         AMZN.String(text),
	}
	result, err := client.DetectEntities(input)
	if err != nil {
		return nil, err
	}
	entities := []*string{}
	for _, v := range result.Entities {
		entities = append(entities, mapping.StrToPtr(v.String()))
	}
	for i, v := range entities {
		entities[i] = mapping.StrToPtr(GetStemOfWord(*v))
	}
	return entities, nil
}

// GetWordRateScore –
func GetWordRateScore(text string) []*types.TopRatedWord {
	text = GetStemsOfText(text)
	frequency := map[string]float32{}
	textAsSlice := strings.Split(text, " ")
	totalWords := len(textAsSlice)

	for _, v := range textAsSlice {
		frequency[v]++
	}

	// convert frequencies to rates
	topRatedWords := []*types.TopRatedWord{}
	for k := range frequency {
		//frequency[k] = frequency[k] / float32(totalWords)
		score := frequency[k] / float32(totalWords)
		ratedWord := &types.TopRatedWord{
			Word:  &k,
			Score: &score,
		}
		topRatedWords = append(topRatedWords, ratedWord)
	}

	sort.Slice(topRatedWords, func(p, q int) bool {
		return *topRatedWords[p].Score < *topRatedWords[q].Score
	})
	return topRatedWords[:5]
}

// RemoveLowCountWords –
func RemoveLowCountWords(text string) string {
	lowCountWords := map[string]int{}

	words := strings.Split(text, " ")
	for _, word := range words {
		lowCountWords[word]++
	}

	newWords := []string{}
	for _, word := range words {
		if lowCountWords[word] > 2 {
			newWords = append(newWords, word)
		}
	}
	return strings.Join(newWords, " ")
}

// GetStemsOfText –
func GetStemsOfText(text string) string {
	s := strings.Split(text, " ")
	for i, v := range s {
		s[i] = GetStemOfWord(v)
	}
	return strings.Join(s, " ")
}

// RemoveStopWords from text
func RemoveStopWords(text string) string {
	cleanedString := stopwords.CleanString(text, "en", false)
	return cleanedString
}

func GetSynset() {
	wn, _ := wordnet.Parse("../../dictionary/dict")
	catNouns := wn.Search("elephant")["n"]
	// = slice of all synsets that contain the word "cat" and are nouns.
	fmt.Println(catNouns)
}

// GetStemOfWord –
func GetStemOfWord(word string) string {
	wordAsBytes := []byte((word))
	res := string(stemmer.Stem(wordAsBytes))
	return res
}

// GetNGrams –
func GetNGrams(input string, n int) [][]string {
	words := strings.Fields(input)
	output := [][]string{}
	endpoint := len(words) - n + 1
	for i := 0; i < endpoint; i++ {
		output = append(output, words[i:i+n])
	}
	return output
}

// GetWordSimilarity –
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
