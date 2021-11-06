package main

import (
	"crypto/rand"
	rand2 "math/rand"
	"time"

	"github.com/muyo/sno"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use: "generate",
	RunE: func(cmd *cobra.Command, args []string) error {
		generated, err := gen()
		if err != nil {
			return err
		}

		for i := 0; i < len(generated); i++ {
			print(generated[i])
		}
		println()

		return nil
	},
}

func gen() ([]string, error) {
	const nPartitions = 10
	const nTokens = 1000

	var generators []*sno.Generator
	for i := 0; i < nPartitions; i++ {
		generator, err := makeGenerator()
		if err != nil {
			return nil, err
		}

		generators = append(generators, generator)
	}

	metabytes := make([]byte, nTokens)
	times := make([]time.Time, nTokens)
	now := time.Now()
	for i := 0; i < nTokens; i++ {
		times[i] = now.Add(time.Duration(rand2.Int31n(100)) * time.Second)

		read, err := rand.Reader.Read(metabytes[i : i+1])
		if err != nil {
			return nil, err
		}
		if read != 1 {
			return nil, errors.New("Not enough data read")
		}
	}

	tokens := make([]sno.ID, nTokens)
	for i := 0; i < nTokens; i++ {
		var generator *sno.Generator
		for {
			generator = generators[rand2.Int31n(int32(len(generators)))]
			newToken := generator.NewWithTime(metabytes[i], times[i])

			if rand2.Int31n(100) < 5 {
				tokens[i] = newToken
				break
			}
		}
	}

	asString := make([]string, nTokens)
	for i, t := range tokens {
		text, err := t.MarshalText()
		if err != nil {
			return nil, err
		}

		asString[i] = string(text)
	}

	return asString, nil
}

func makeGenerator() (*sno.Generator, error) {
	b := make([]byte, 2)
	_, err := rand.Reader.Read(b)
	if err != nil {
		return nil, err
	}

	c := make(chan *sno.SequenceOverflowNotification)
	generator, err := sno.NewGenerator(&sno.GeneratorSnapshot{
		Partition: sno.Partition{
			b[0], b[1],
		},
	}, c)
	if err != nil {
		return nil, errors.Wrap(err, "making coder")
	}
	return generator, nil
}
