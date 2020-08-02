package nlp_service

import (
	"fmt"
	"io/ioutil"
	"testing"

	"code.mine/dating_server/aws"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type NLPTestSuite struct {
	suite.Suite
}

func (suite *NLPTestSuite) SetupSuite() {

	err := godotenv.Load("../../.env")
	suite.NoError(err)
	err = aws.SetAWSConnection()
	suite.NoError(err)

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

func (suite *NLPTestSuite) TestGetTextBreakDown() {
	content, err := ioutil.ReadFile("../../tests/articles/space_2.txt")
	suite.NoError(err)

	text := string(content)
	textSummary, err := GetTextBreakDown(text)
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
	suite.Run(t, new(NLPTestSuite))
}

// can you do something with levenshtein distance adn trigrams or something?
// that way you can detect phrases
