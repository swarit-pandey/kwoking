package common

import (
	"math/rand"
	"time"
)

const (
	MaxLabelLength      = 63
	MaxLabelValueLength = 63
	MaxPrefixLength     = 253
	Prefix              = "simulated"
)

var (
	labelNameChars  = []rune("abcdefghijklmnopqrstuvwxyz0123456789-_")
	labelValueChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
)

func GenerateRandomLabels(numMaps, numLabelsPerMap int) []map[string]string {
	result := make([]map[string]string, 0, numMaps)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numMaps; i++ {
		labels := make(map[string]string, numLabelsPerMap)

		for j := 0; j < numLabelsPerMap; j++ {
			var key string
			key = Prefix + "/" + randomRunes(labelNameChars, 16, rnd)

			var value string
			value = randomRunes(labelValueChars, 16, rnd)

			if _, exists := labels[key]; exists {
				continue
			}

			labels[key] = value
		}

		result = append(result, labels)
	}

	return result
}

func randomRunes(allowed []rune, length int, rnd *rand.Rand) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = allowed[rnd.Intn(len(allowed))]
	}
	return string(b)
}
