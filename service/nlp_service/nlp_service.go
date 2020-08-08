package nlp_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/aws"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"github.com/bbalet/stopwords"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/agnivade/levenshtein"
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

// ExtractEntitiesFromArticle –
// when you process entity handling, check there if entity length is 0
// func ExtractEntitiesFromArticle(articleSummary *types.Article) ([]string, error) {
// 	entities := []string{}

// 	if articleSummary.Summary == nil {
// 		return nil, errors.New("article text summary cannot be nil")
// 	}
// 	currentEntities := articleSummary.Summary.Entities
// 	if currentEntities == nil {
// 		return nil, errors.New("entity array cannot be empty or null")
// 	}
// 	return

// }

// GetSimilarTexts –
// you are comparing texts (all the other articles) against src
func GetSimilarTexts(src types.TextSummary, texts []string) error {
	if len(texts) == 0 {
		return errors.New("can't have no texts to compare against")
	}
	if src.Entities == nil {
		return errors.New("can't have nil src entities")
	}
	if src.Keyphrases == nil {
		return errors.New("can't have nil src keyphrases")
	}
	if src.TopRatedWords == nil {
		return errors.New("can't have nil src TopRatedWords")
	}

	count := 0
	textSummariesMap := map[int]types.TextSummary{}
	for _, v := range texts {
		textSum, err := GetTextBreakDown(v)
		if err != nil {
			continue
		}
		textSummariesMap[count] = *textSum
		count++
	}
	if len(textSummariesMap) == 0 {
		return errors.New("unable to create text summaries for articles")
	}

	// map each candidate text to a score
	textToScore := map[int]float64{}
	for key, candidate := range textSummariesMap {
		// GetNormalizedEntityScore
		// GetNormalizedKeyphraseScore
		// GetNormalizedTopWordsScore

		var entityScore float64
		var topRatedWordsScore float64
		var keyPhrasesScore float64

		// process entities
		if len(src.Entities) > 0 && len(candidate.Entities) > 0 {
			maxAmtEntities := math.Max(float64(len(src.Entities)), float64(len(candidate.Entities)))

			for _, srcEntity := range src.Entities {
				for _, canEntity := range candidate.Entities {
					normalizedDistance := GetNormalizedWordSimilarityScore(canEntity, srcEntity)
					entityScore += normalizedDistance
				}
			}
			entityScore = entityScore / maxAmtEntities

			// 0.333 is temp weighted score
			textToScore[key] += entityScore * 0.33333
		}

		// process topratedword
		if len(src.TopRatedWords) > 0 && len(candidate.TopRatedWords) > 0 {
			maxAmtEntities := math.Max(float64(len(src.TopRatedWords)), float64(len(candidate.TopRatedWords)))

			for _, srcTopRatedWord := range src.TopRatedWords {
				for _, canTopRatedWord := range candidate.TopRatedWords {
					normalizedDistance := GetNormalizedWordSimilarityScore(canTopRatedWord.Word, srcTopRatedWord.Word)
					topRatedWordsScore += normalizedDistance
				}
			}
			topRatedWordsScore = topRatedWordsScore / maxAmtEntities

			// 0.333 is temp weighted score
			textToScore[key] += topRatedWordsScore * 0.33333
		}

		// process topratedword
		if len(src.Keyphrases) > 0 && len(candidate.Keyphrases) > 0 {
			maxAmtEntities := math.Max(float64(len(src.Keyphrases)), float64(len(candidate.Keyphrases)))

			for _, srcKP := range src.Keyphrases {
				for _, canKP := range candidate.Keyphrases {
					normalizedDistance := GetNormalizedWordSimilarityScore(canKP, srcKP)
					keyPhrasesScore += normalizedDistance
				}
			}
			keyPhrasesScore = keyPhrasesScore / maxAmtEntities

			// 0.333 is temp weighted score
			textToScore[key] += keyPhrasesScore * 0.33333
		}

	}
	return nil
}

func GetNormalizedWordSimilarityScore(a, b string) float64 {
	distance := float64(levenshtein.ComputeDistance(a, b))
	maxStringLength := math.Max(float64(len(a)), float64(len(b)))
	normalizedDistance := distance / maxStringLength
	return normalizedDistance
}

// GetTextBreakDown –
func GetTextBreakDown(text string) (*types.TextSummary, error) {
	textSummary := &types.TextSummary{}

	// text, err := RemoveNonAlphaNumericFromString(text)
	// if err != nil {
	// 	return nil, err
	// }
	// remove punctiation and stuff?

	sess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}
	client := comprehend.New(sess)
	topRatedWords, err := GetWordRateScore(text, client)
	if err != nil {
		return nil, err
	}
	textSummary.TopRatedWords = topRatedWords

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

	// ngrams := GetNGrams(text, 3)
	return textSummary, nil
}

// GetSimilarityOfEntities –
// https://usetrove.io/ check this for web crawling
func GetSimilarityOfEntities(sourceEntities, candidateEntities []string) float64 {
	var finalScore float64
	var totalNormalizingLength int

	// are we doing this somewhere else though when we get entities?
	for i, v := range sourceEntities {
		sourceEntities[i] = GetStemsOfText(v)
		sourceEntities[i] = strings.ToLower(sourceEntities[i])

	}
	for i, v := range candidateEntities {
		candidateEntities[i] = GetStemsOfText(v)
		candidateEntities[i] = strings.ToLower(candidateEntities[i])
	}

	for _, srcEntity := range sourceEntities {

		bestScoreCandidateEntity := math.Inf(1)
		var bestWordCandidateEntity string
		for _, candEntity := range candidateEntities {
			distance := float64(levenshtein.ComputeDistance(candEntity, srcEntity))
			if distance < bestScoreCandidateEntity {
				bestScoreCandidateEntity = distance
				bestWordCandidateEntity = candEntity
			}
		}
		finalScore += bestScoreCandidateEntity
		maxWordCount := math.Max(
			float64(len(bestWordCandidateEntity)),
			float64(len(srcEntity)))
		totalNormalizingLength += int(maxWordCount)
	}
	return 1 - float64(finalScore)/float64(totalNormalizingLength)
}

// AddSynsForText –
func AddSynsForText(text string) error {
	text = RemoveStopWords(text)
	text = strings.Replace(text, "'", "", -1)
	textAsSlice := strings.Split(text, " ")

	var wg sync.WaitGroup

	buffer := make(chan struct{}, 200)
	numErrors := 0
	for _, word := range textAsSlice {

		buffer <- struct{}{}
		go func(w string) {
			wg.Add(1)
			_, err := GetWordInformation(w)
			if err != nil {
				numErrors++
			}
			<-buffer
			wg.Done()
		}(word)
	}
	if numErrors == len(textAsSlice) {
		return errors.New("no words could be captured by thesaurus api")
	}
	wg.Wait()
	return nil
}

// work on this tomorrow
// GetWordInformation –
func GetWordInformation(word string) (*types.WordInformation, error) {
	wordInfoFromDB, err := GetWordFromCollection(word)
	if err != nil {
		return nil, err
	}
	if wordInfoFromDB != nil {
		return wordInfoFromDB, nil
	}

	// erase dogs
	// wordInfo := types.WordInformation{}
	wordInfo := types.WordInformation{}

	baseURL := `https://words.bighugelabs.com/api/1`
	apiKey := os.Getenv("BIG_HUGE_LABS_KEY")
	url := fmt.Sprintf(`%s/%s/%s/json`, baseURL, apiKey, word)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	values := []string{}
	err = json.NewDecoder(resp.Body).Decode(&values)
	if err != nil {
		return nil, err
	}

	wordInfo.Word = word
	wordInfo.WordList = values
	wordInfo.WordStem = GetStemOfWord(word)

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	wordInfo.CreatedAt = now
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	wordInfo.UUID = newUUID.String()

	err = AddWordToCollection(&wordInfo)
	if err != nil {
		return nil, err
	}

	return &wordInfo, nil
}

// GetWordFromCollection –
// return the list of syns if it exists
// TODO – logical OR filter for stem or word
func GetWordFromCollection(word string) (*types.WordInformation, error) {
	c, err := DB.GetCollection("word_synonyms")
	if err != nil {
		return nil, err
	}

	wordStem := GetStemOfWord(word)
	res := c.FindOne(context.Background(), bson.M{"wordStem": wordStem})
	err = res.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	wordInfo := &types.WordInformation{}
	wordInfo.WordStem = wordStem
	err = res.Decode(wordInfo)
	if err != nil {
		return nil, err
	}
	return wordInfo, nil
}

// AddWordToCollection –
func AddWordToCollection(wordInfo *types.WordInformation) error {
	c, err := DB.GetCollection("word_synonyms")
	if err != nil {
		return err
	}
	_, err = c.InsertOne(context.Background(), wordInfo)
	if err != nil {
		return err
	}
	return nil
}

// GetKeyphrasesOfText –
func GetKeyphrasesOfText(text string, client *comprehend.Comprehend) ([]string, error) {
	input := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: AMZN.String("en"),
		Text:         AMZN.String(text),
	}
	result, err := client.DetectKeyPhrases(input)
	if err != nil {
		return nil, err
	}
	sort.Slice(result.KeyPhrases, func(p, q int) bool {
		return *result.KeyPhrases[p].Score < *result.KeyPhrases[q].Score
	})

	keyphrases := []string{}
	for _, v := range result.KeyPhrases {
		temp := strings.Split(*v.Text, " ")
		for idx, val := range temp {
			temp[idx] = GetStemOfWord(val)
		}
		str := strings.Join(temp, " ")
		keyphrases = append(keyphrases, str)
	}
	length := len(keyphrases)
	if length > 5 {
		length = 5
	}
	return keyphrases[:length], nil
}

// GetEntitiesOfText –
func GetEntitiesOfText(text string, client *comprehend.Comprehend) ([]string, error) {
	input := &comprehend.DetectEntitiesInput{
		LanguageCode: AMZN.String("en"),
		Text:         AMZN.String(text),
	}
	result, err := client.DetectEntities(input)
	if err != nil {
		return nil, err
	}
	sort.Slice(result.Entities, func(p, q int) bool {
		return *result.Entities[p].Score < *result.Entities[q].Score
	})

	entities := []string{}
	set := map[string]bool{}
	for _, v := range result.Entities {
		word := *v.Text
		if !set[word] {
			entities = append(entities, *v.Text)
			set[word] = true
		}
	}

	length := len(entities)
	if length > 5 {
		length = 5
	}
	return entities[:length], nil
}

// GetWordRateScore –
func GetWordRateScore(text string, client *comprehend.Comprehend) ([]types.TopRatedWord, error) {
	text = RemoveStopWords(text)
	text = RemoveLowCountWords(text)
	text = strings.Replace(text, "'", "", -1)

	syntaxInput := &comprehend.DetectSyntaxInput{
		LanguageCode: AMZN.String("en"),
		Text:         AMZN.String(text),
	}
	output, err := client.DetectSyntax(syntaxInput)
	if err != nil {
		return nil, err
	}

	textAsSlice := []string{}
	for _, tkn := range output.SyntaxTokens {
		if *tkn.PartOfSpeech.Tag == "PROPN" ||
			*tkn.PartOfSpeech.Tag == "PRON" ||
			*tkn.PartOfSpeech.Tag == "NOUN" ||
			*tkn.PartOfSpeech.Tag == "SYM" {
			textAsSlice = append(textAsSlice, *tkn.Text)
		}
	}
	// text = strings.Join(words, " ")
	for i := range textAsSlice {
		textAsSlice[i] = GetStemOfWord(textAsSlice[i])
	}
	// text = GetStemsOfText(text)
	frequency := map[string]float32{}
	//textAsSlice := strings.Split(text, " ")
	totalWords := len(textAsSlice)

	for _, v := range textAsSlice {
		frequency[v]++
	}

	// convert frequencies to rates
	topRatedWords := []types.TopRatedWord{}
	for k := range frequency {
		//frequency[k] = frequency[k] / float32(totalWords)
		score := frequency[k] / float32(totalWords)
		ratedWord := types.TopRatedWord{}
		ratedWord.Score = score
		ratedWord.Word = k
		topRatedWords = append(topRatedWords, ratedWord)
	}

	sort.Slice(topRatedWords, func(p, q int) bool {
		return topRatedWords[p].Score < topRatedWords[q].Score
	})
	length := len(topRatedWords)
	if length >= 5 {
		length = 5
	}
	return topRatedWords[:length], nil
}

// RemoveNonAlphaNumericFromString –
func RemoveNonAlphaNumericFromString(text string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	processedString := reg.ReplaceAllString(text, "")
	return processedString, nil
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
