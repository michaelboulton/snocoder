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

		for i := 0; i < 10; i++ {
			println(generated[i])
		}

		return nil
	},
}

func gen() ([]string, error) {
	c := make(chan *sno.SequenceOverflowNotification)

	const nPartitions = 10
	const nTokens = 1000

	generators := make([]*sno.Generator, nPartitions)
	for i := 0; i < nPartitions; i++ {
		b := make([]byte, 2)
		_, err := rand.Reader.Read(b)
		if err != nil {
			return nil, err
		}

		generator, err := sno.NewGenerator(&sno.GeneratorSnapshot{
			Partition: sno.Partition{
				b[0], b[1],
			},
		}, c)
		if err != nil {
			return nil, errors.Wrap(err, "making coder")
		}

		generators = append(generators, generator)
	}


	metabytes := make([]byte, nTokens)
	times := make([]time.Time, nTokens)
	now := time.Now()
	for i := 0; i < nTokens; i++ {
		times[i] = now.Add(time.Duration(i) * time.Second)

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
		generator := generators[rand2.Int31n(int32(len(generators)))]

		newToken := generator.NewWithTime(metabytes[i], times[i])
		tokens[i] = newToken
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
