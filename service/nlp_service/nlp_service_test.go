package nlp_service

import (
	"fmt"
	"testing"

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

func (suite *NLPTestSuite) TestGetStemOfWord() {
	_ = GetStem()

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
