package nlp_service

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"os"
	"path/filepath"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/aws"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type NLPTestSuite struct {
	suite.Suite
}

func (suite *NLPTestSuite) SetupSuite() {

}
func (suite *NLPTestSuite) SetupTest() {

}

func (suite *NLPTestSuite) TearDownAllSuite() {

}

func (suite *NLPTestSuite) TearDownTest() {
}

func (suite *NLPTestSuite) TestGetSynset() {
	GetSynset()

}

// try cycling through all the words in the synset?
func (suite *NLPTestSuite) TestGetSimilarityOfSynsets() {
	num, err := GetWordSimilarity("blue", "green")
	suite.NoError(err)

	num, err = GetWordSimilarity("zebra", "horse")
	suite.NoError(err)

	num, err = GetWordSimilarity("house", "housing")
	suite.NoError(err)
	fmt.Println(num)

}

// still need better way to handle things like word not found
// and lines vs line
func (suite *NLPTestSuite) TestAddWordInformation() {

	content, err := ioutil.ReadFile("../../tests/articles/solar_1.txt")
	suite.NoError(err)

	contentAsString := RemoveStopWords(string(content))
	contentAsString = "world"
	contentAsSlice := strings.Split(contentAsString, " ")
	frequency := map[string]int{}
	for _, v := range contentAsSlice {
		index := strings.Index(v, "'")
		if index != -1 {
			v = v[:index]
		}
		frequency[v]++
	}

	for k := range frequency {

		if len(k) >= 3 {
			wordInfo, err := GetWordInformation(k)
			if err != nil {
				continue
			}
			suite.NoError(err)
			suite.NotNil(wordInfo)
		}

	}

}

// sick too https://mholt.github.io/json-to-go/
func (suite *NLPTestSuite) TestGetWordInformation() {
	wordInfo, err := GetWordInformation("cats")
	suite.NoError(err)
	suite.NotNil(wordInfo)
}

/*
func (suite *NLPTestSuite) TestGetStemOfWord() {
	_ = GetStemOfWord()

}


func (suite *NLPTestSuite) TestNGrams() {
	word := "snow is good and snow is bad"
	tokens := GetNGramOfString(word, 3)
	suite.NotNil(tokens)
}
*/

// you should consider hardcoding in the synoyms into the struct
// making it like a one time thing
// so if yo're getting a noramlized score, then you can make it count for diff things
// like different %'s.
func (suite *NLPTestSuite) TestGetSimilarityOfEntities() {
	sourceEntities := []string{"space", "satellite", "NASA", "moon"}
	candidateEntities := []string{"outer space", "satellites", "moon crater", "NASA agency"}
	scoreOne := GetSimilarityOfEntities(sourceEntities, candidateEntities)

	candidateEntities = []string{"grilled cheese", "food", "pizza", "eating"}
	scoreTwo := GetSimilarityOfEntities(sourceEntities, candidateEntities)
	suite.True(scoreOne > scoreTwo)
}

func (suite *NLPTestSuite) TestGetSimilarTexts() {
	articles := []string{}

	root := "../../tests/articles"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		content, err := ioutil.ReadFile(path)
		suite.NoError(err)

		text := string(content)
		articles = append(articles, text)
		return nil
	})

	content, err := ioutil.ReadFile("../../tests/articles/space_2.txt")
	suite.NoError(err)

	text := string(content)
	textSummary, err := GetTextBreakDown(text)
	suite.NoError(err)
	suite.NotNil(textSummary)

	suite.NoError(err)
	GetSimilarTexts(*textSummary, articles)

}

func (suite *NLPTestSuite) TestGetTextBreakDown() {
	content, err := ioutil.ReadFile("../../tests/articles/space_2.txt")
	suite.NoError(err)

	text := string(content)
	textSummary, err := GetTextBreakDown(text)
	suite.NoError(err)
	suite.NotNil(textSummary)

	content, err = ioutil.ReadFile("../../tests/articles/space_3.txt")
	suite.NoError(err)

	text = string(content)
	textSummary, err = GetTextBreakDown(text)
	suite.NoError(err)
	suite.NotNil(textSummary)
}

func (suite *NLPTestSuite) TestGetSynonyms() {
	word := "snow"
	syns, err := GetSynonyms(&word)
	suite.NoError(err)
	suite.NotNil(syns)
}

func TestNLPTestSuite(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = aws.SetAWSConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := DB.SetupDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Client.Disconnect(*db.Ctx)
	suite.Run(t, new(NLPTestSuite))
}

// can you do something with levenshtein distance adn trigrams or something?
// that way you can detect phrases
