package builders

import (
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	"spot-oceancd-cli/pkg/oceancd/repositories"
	"sync"
)

type DetailedRolloutBuilder struct {
	detailedRollout   *rollout.DetailedRollout
	rolloutRepository *repositories.RolloutRepository
	errors            chan error
	wg                *sync.WaitGroup
	withStrategy      bool
}

func NewDetailedRolloutBuilder(rolloutRepository *repositories.RolloutRepository) *DetailedRolloutBuilder {
	return &DetailedRolloutBuilder{
		rolloutRepository: rolloutRepository,
	}
}

func (b *DetailedRolloutBuilder) Build(rolloutId string) (*rollout.DetailedRollout, error) {
	b.detailedRollout = &rollout.DetailedRollout{}
	b.errors = make(chan error)

	b.wg = &sync.WaitGroup{}

	if b.withStrategy {
		//here b.wg begins waiting for the next goroutine: setStrategy()
		b.wg.Add(1)
		go b.setStrategy(rolloutId)
	}

	//here b.wg begins waiting for the next goroutines: setRollout(), setRolloutPhases(), setRolloutVerifications()
	b.wg.Add(3)

	go b.setRollout(rolloutId)
	go b.setRolloutPhases(rolloutId)
	go b.setRolloutVerifications(rolloutId)

	b.wg.Wait()

	select {
	case err := <-b.errors:
		return b.detailedRollout, err
	default:
		return b.detailedRollout, nil
	}
}

func (b *DetailedRolloutBuilder) setStrategy(rolloutId string) {
	defer b.wg.Done()

	strategy, err := b.rolloutRepository.GetStrategy(rolloutId)
	if err != nil {
		b.errors <- err
	}

	b.detailedRollout.Definition.Strategy = strategy
}

func (b *DetailedRolloutBuilder) setRollout(rolloutId string) {
	defer b.wg.Done()

	fetchedRollout, err := oceancd.GetRollout(rolloutId)
	if err != nil {
		b.errors <- err
	}

	b.detailedRollout.Rollout = fetchedRollout
}

func (b *DetailedRolloutBuilder) setRolloutPhases(rolloutId string) {
	defer b.wg.Done()

	phases, err := oceancd.GetRolloutPhases(rolloutId)
	if err != nil {
		b.errors <- err
	}

	b.detailedRollout.Phases = phases
}

func (b *DetailedRolloutBuilder) setRolloutVerifications(rolloutId string) {
	defer b.wg.Done()

	verifications, err := oceancd.GetRolloutVerifications(rolloutId)
	if err != nil {
		b.errors <- err
	}

	b.detailedRollout.Verifications = verifications
}

func (b *DetailedRolloutBuilder) WithStrategy() *DetailedRolloutBuilder {
	b.withStrategy = true
	return b
}
