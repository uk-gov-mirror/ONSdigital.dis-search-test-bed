// Package app runs search-relevance evaluations against a throwaway
// Elasticsearch instance: it indexes the test-set documents, evaluates every
// term with one or more search algorithms, and either logs the scores with a
// comparison of each algorithm's NDCG (Compare) or writes the ranked results
// to a CSV file (Export).
package app

import (
	"context"
	"strings"

	"github.com/ONSdigital/dis-search-test-bed/algorithm"
	"github.com/ONSdigital/dis-search-test-bed/testset/stream"
	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v4/client"
	"github.com/pkg/errors"
)

// exportAlgorithm is the algorithm whose ranking Export evaluates. The export
// CSV has no algorithm column, so the export stays single-algorithm.
const exportAlgorithm = algorithm.SearchAlgorithmBaseline

// App holds the stores the evaluation operates on.
type App struct {
	Documents  stream.Stream[stream.Item] // indexed into Elasticsearch (see storeTargets)
	Terms      stream.Stream[stream.Item] // evaluation data: query terms (not indexed)
	Judgements stream.Stream[stream.Item] // evaluation data: relevance labels (not indexed)
}

// New wires the App to the embedded fixture stores.
func New() *App {
	return &App{
		Documents:  stream.NewDocumentStore(),
		Terms:      stream.NewTermStore(),
		Judgements: stream.NewJudgementStore(),
	}
}

// Compare evaluates every test term with each of the given algorithms against
// the same index, logs their full-corpus relevance scores, then prints a table
// comparing each algorithm's NDCG. An empty algorithms slice evaluates every
// registered algorithm.
func (a *App) Compare(ctx context.Context, algorithms []algorithm.SearchAlgorithm) error {
	if len(algorithms) == 0 {
		algorithms = algorithm.AllSearchAlgorithms()
	}

	return a.withElasticsearch(ctx, func(ctx context.Context, esClient dpEsClient.Client) error {
		evaluations, err := a.evaluateTerms(ctx, esClient, algorithms)
		if err != nil {
			return err
		}
		return reportComparison(evaluations)
	})
}

// Export evaluates every test term and writes the ranked results and their
// current relevance judgements as CSV to outputPath.
func (a *App) Export(ctx context.Context, outputPath string) error {
	if strings.TrimSpace(outputPath) == "" {
		return errors.New("output path is required")
	}

	return a.withElasticsearch(ctx, func(ctx context.Context, esClient dpEsClient.Client) error {
		evaluations, err := a.evaluateTerms(ctx, esClient, []algorithm.SearchAlgorithm{exportAlgorithm})
		if err != nil {
			return err
		}
		return writeEvaluationFile(outputPath, evaluations)
	})
}
