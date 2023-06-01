package viewcontroller

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"reflect"
	"sort"
	"spot-oceancd-cli/pkg/oceancd/builders"
	"spot-oceancd-cli/pkg/oceancd/model"
	"spot-oceancd-cli/pkg/oceancd/model/phase"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	"spot-oceancd-cli/pkg/oceancd/model/verification"
	"spot-oceancd-cli/pkg/oceancd/repositories"
	"spot-oceancd-cli/pkg/utils"
	"spot-oceancd-cli/viewcontroller/converter"
	"sync"
	"text/tabwriter"
	"time"
)

const (
	tableFormat       = "%-21s%v\n"
	columnPrefix      = "│"
	separatingRawPart = "──────────"
	rawTemplate       = "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n"
	subRowOffset      = "  "
)

// icons
const (
	iconWaiting    = "◷"
	iconOk         = "✔"
	iconFailed     = "✖"
	iconInProgress = "◌"
	iconPaused     = "॥"
	iconAborted    = "↵"
	iconCanceled   = "⊗"
	iconWarning    = "⚠"
	iconPoint      = "•"
)

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/viewcontroller/viewcontroller.go#L31
// viewController is a mini controller which allows printing of live updates to rollouts
// Allows subscribers to receive updates about
type viewController struct {
	writer io.Writer
	color  color.Color
}

func newViewController(noColor bool) *viewController {
	if noColor {
		color.NoColor = true
	}
	return &viewController{writer: os.Stdout}
}

func (c *viewController) colorizeWith(text string, col color.Attribute) string {
	return color.New(col).Sprint(text)
}

func (c *viewController) colorize(text string) string {
	switch text {
	case iconCanceled, iconPoint:
		return color.New(color.FgHiBlack).Sprint(text)
	case iconWaiting, iconInProgress:
		return color.New(color.FgCyan).Sprint(text)
	case iconOk:
		return color.New(color.FgGreen).Sprint(text)
	case iconFailed, iconAborted, iconWarning:
		return color.New(color.FgRed).Sprint(text)
	case iconPaused:
		return color.New(color.FgYellow).Sprint(text)
	default:
		return color.New(color.Reset).Sprint(text)
	}
}

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/viewcontroller/viewcontroller.go#L53
type RolloutViewController struct {
	*viewController
	rolloutId       string
	rollout         *rollout.DetailedRollout
	previousRollout *rollout.DetailedRollout
}

func NewRolloutViewController(rolloutId string, noColor bool) *RolloutViewController {
	vc := newViewController(noColor)

	return &RolloutViewController{
		viewController: vc,
		rolloutId:      rolloutId,
	}
}

func (c *RolloutViewController) GetRollout() (*rollout.DetailedRollout, error) {
	detailedRolloutBuilder := builders.NewDetailedRolloutBuilder(repositories.NewRolloutRepository())

	if c.previousRollout == nil {
		return detailedRolloutBuilder.WithStrategy().Build(c.rolloutId)
	}

	detailedRollout, err := detailedRolloutBuilder.Build(c.rolloutId)
	detailedRollout.Definition.Strategy = c.previousRollout.Definition.Strategy

	return detailedRollout, err
}

func (c *RolloutViewController) PrintRollout(detailedRollout *rollout.DetailedRollout) {
	c.rollout = detailedRollout
	fmt.Fprintf(c.writer, tableFormat, "Start Time:", c.rollout.StartTime)
	if c.rollout.EndTime != "" {
		fmt.Fprintf(c.writer, tableFormat, "End Time:", c.rollout.EndTime)
	}
	fmt.Fprintf(c.writer, tableFormat, "SpotDeploymentName:", c.rollout.SpotDeployment)
	fmt.Fprintf(c.writer, tableFormat, "Cluster ID:", c.rollout.ClusterId)
	fmt.Fprintf(c.writer, tableFormat, "Namespace:", c.rollout.Namespace)
	c.printStrategy()
	fmt.Fprintf(c.writer, tableFormat, "Status:", fmt.Sprintf("%s %s", c.statusIcon(c.rollout.Status), converter.RolloutStatus(c.rollout.Status)))
	c.printPhasesNumber()
	c.printVersions()
	if len(c.rollout.GetBackgroundVerifications()) > 0 {
		fmt.Fprintf(c.writer, "%s\n", "BackgroundVerification:")
		c.printBackgroundVerifications()
	}

	c.printPhases()
}

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/viewcontroller/viewcontroller.go#L144
func (c *RolloutViewController) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go wait.Until(func() {
		for c.processRollout() {
		}
	}, time.Second, ctx.Done())
	<-ctx.Done()
}

func (c *RolloutViewController) processRollout() bool {
	rolloutInfo, err := c.GetRollout()
	if err != nil {
		fmt.Printf("%s/n", err)
		return false
	}

	if !reflect.DeepEqual(c.previousRollout, rolloutInfo) {
		c.clear()
		c.PrintRollout(rolloutInfo)
		c.previousRollout = rolloutInfo
	}
	return true
}

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/cmd/get/get.go#L133
func (c *RolloutViewController) clear() {
	fmt.Fprint(c.writer, "\033[H\033[2J")
	fmt.Fprint(c.writer, "\033[0;0H")
}

func (c *RolloutViewController) printPhases() {
	if len(c.rollout.Phases) < 1 {
		return
	}

	writer := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', tabwriter.TabIndent)
	c.writer = writer
	c.printHeader()
	c.identifyPhaseStatuses()

	for i, rolloutPhase := range c.rollout.Phases {
		c.printPhase(rolloutPhase, columnPrefix, i)
		c.printVerifications(c.orderVerifications(rolloutPhase.Verifications), rolloutPhase)
		c.printSeparatingRaw(i == len(c.rollout.Phases)-1)
	}

	writer.Flush()
}

func (c *RolloutViewController) printHeader() {
	fmt.Fprint(c.writer, fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		c.colorize("PHASE"), c.colorize("NAME"), c.colorize("STATUS")+c.iconStub(), c.colorize("WEIGHT"),
		c.colorize("METRICS"), c.colorize("VERIFICATION"), c.colorize("VERIFICATION")+c.iconStub()))

	fmt.Fprint(c.writer, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		c.emptyCell(), c.emptyCell(), c.emptyCell()+c.iconStub(), c.emptyCell(), c.emptyCell(),
		c.colorize("PROVIDER"), c.colorize("STATUS")+c.iconStub()))
}

func (c *RolloutViewController) printPhase(phase phase.Phase, prefix string, index int) {
	raw := fmt.Sprint(prefix, " ",
		fmt.Sprintf(rawTemplate,
			c.colorizeWith(converter.PhaseIndex(index+1), color.Bold),
			c.colorizeWith(converter.PhaseName(phase), color.Bold),
			fmt.Sprintf("%s %s", c.phaseStatusIcon(phase.Status), c.colorizeWith(converter.PhaseStatus(phase), color.Bold)),
			c.colorizeWith(converter.Weight(phase), color.Bold),
			c.emptyCell(), c.emptyCell(), c.emptyCell()+c.iconStub(),
		))

	fmt.Fprint(c.writer, raw)

	if headers, ok := c.rollout.Definition.Strategy.GetHeaderRouteMatchesBySteps()[phase.Name]; ok &&
		len(phase.Verifications) == 0 &&
		len(headers) > 0 {
		headersLength := len(headers)
		headersRaw := fmt.Sprint(columnPrefix, " ",
			fmt.Sprintf(rawTemplate,
				c.colorize(fmt.Sprintf("%d %s", headersLength, utils.GetNounForm("Header", headersLength))),
				c.emptyCell(), c.emptyCell(), c.emptyCell(),
				c.emptyCell(), c.emptyCell(), c.emptyCell()+c.iconStub(),
			))

		fmt.Fprint(c.writer, headersRaw)
	}
}

func (c *RolloutViewController) printVerifications(verifications []verification.Verification, rolloutPhase phase.Phase) {
	if len(verifications) < 1 {
		return
	}

	headersCell := c.emptyCell()

	if headers, ok := c.rollout.Definition.Strategy.GetHeaderRouteMatchesBySteps()[rolloutPhase.Name]; ok && len(headers) > 0 {
		headersLength := len(headers)
		headersCell = c.colorize(fmt.Sprintf(" %d %s", headersLength, utils.GetNounForm("Header", headersLength)))
	}

	for _, verificationItem := range verifications {
		raw := fmt.Sprint(columnPrefix,
			fmt.Sprintf(rawTemplate,
				headersCell,
				c.emptyCell(), c.emptyCell()+c.iconStub(), c.emptyCell(),
				c.colorize(verificationItem.MetricName),
				c.colorize(verificationItem.Provider),
				fmt.Sprintf("%s %s", c.verificationStatusIcon(verificationItem.Status), c.colorize(converter.VerificationStatus(verificationItem))),
			))
		fmt.Fprint(c.writer, raw)
	}
}

func (c *RolloutViewController) printSeparatingRaw(isLast bool) {
	prefix := columnPrefix

	if isLast {
		prefix = "└"
	}

	fmt.Fprint(c.writer,
		fmt.Sprintf(prefix+rawTemplate,
			c.colorize(separatingRawPart), c.colorize(separatingRawPart), c.colorize(separatingRawPart)+c.iconStub(),
			c.colorize(separatingRawPart), c.colorize(separatingRawPart), c.colorize(separatingRawPart),
			c.colorize(separatingRawPart)+c.iconStub()))
}

func (c *RolloutViewController) emptyCell() string {
	return c.colorize(" ")
}

func (c *RolloutViewController) iconStub() string {
	return color.New(color.FgBlack).Sprint("")
}

func (c *RolloutViewController) identifyPhaseStatuses() {
	var fullyPromoted bool
	var canceled bool
	var aborted bool
	var phases []phase.Phase

	switch c.rollout.Status {
	case rollout.Aborted:
		aborted = true
	case rollout.Canceled:
		canceled = true
	}

	for i, rolloutPhase := range c.rollout.Phases {
		if fullyPromoted {
			rolloutPhase.Status = phase.Dropped
			if i == len(c.rollout.Phases)-1 {
				rolloutPhase.Status = phase.Finished
			}
		}

		if aborted && rolloutPhase.IsUncompleted() {
			rolloutPhase.Status = phase.Dropped
		}

		if canceled && rolloutPhase.IsUncompleted() {
			rolloutPhase.Status = phase.Canceled
		}

		switch rolloutPhase.Status {
		case phase.FullPromoted:
			fullyPromoted = true
		}

		phases = append(phases, rolloutPhase)
	}
	c.rollout.Phases = phases
}

func (c *RolloutViewController) printReplicasNumber(replicas rollout.ReplicasInfo) {
	fmt.Fprintf(c.writer, "%-21s%s %d%s %d%s %d%s %d\n", subRowOffset+"Replicas: ",
		"Desired:", replicas.Desired, " | Ready:", replicas.Ready,
		" | InProgress:", replicas.InProgress, " | Failed:", replicas.Failed,
	)
}

func (c *RolloutViewController) printVersionStatus(status rollout.VersionStatus, color color.Attribute) {
	fmt.Fprintf(c.writer, tableFormat, subRowOffset+"Version:", c.colorizeWith(status.Version, color))
	fmt.Fprintf(c.writer, tableFormat, subRowOffset+"TrafficPercentage:", status.TrafficPercentage)
	if status.K8sService != "" {
		fmt.Fprintf(c.writer, tableFormat, subRowOffset+"ServiceName:", status.K8sService)
	}
}

func (c *RolloutViewController) statusIcon(status rollout.Status) string {
	switch status {
	case rollout.Pending:
		return c.colorize(iconWaiting)
	case rollout.InProgress, rollout.Aborting, rollout.ManualPausing, rollout.Deallocating, rollout.Verifying, rollout.FailurePolicyPausing,
		rollout.BackgroundVerifying:
		return c.colorize(iconWaiting)
	case rollout.Paused, rollout.ManualPaused, rollout.FailurePolicyPaused:
		return c.colorize(iconPaused)
	case rollout.Aborted:
		return c.colorize(iconAborted)
	case rollout.Failed:
		return c.colorize(iconFailed)
	case rollout.InvalidSpec:
		return c.colorize(iconWarning)
	case rollout.Finished:
		return c.colorize(iconOk)
	case rollout.Canceled:
		return c.colorize(iconCanceled)
	default:
		return ""
	}
}

func (c *RolloutViewController) phaseStatusIcon(status phase.Status) string {
	switch status {
	case phase.Pending:
		return c.colorize(iconWaiting)
	case phase.InProgress, phase.Aborting, phase.Verifying, phase.Promoting:
		return c.colorize(iconWaiting)
	case phase.Paused:
		return c.colorize(iconPaused)
	case phase.Aborted:
		return c.colorize(iconAborted)
	case phase.Finished, phase.Promoted, phase.FullPromoted:
		return c.colorize(iconOk)
	case phase.Canceled:
		return c.colorize(iconCanceled)
	case phase.Dropped:
		return c.colorize(iconPoint)
	default:
		return ""
	}
}

func (c *RolloutViewController) verificationStatusIcon(status verification.Status) string {
	switch status {
	case verification.Successful:
		return c.colorize(iconOk)
	case verification.Failed:
		return c.colorize(iconFailed)
	case verification.Running:
		return c.colorize(iconInProgress)
	case verification.Error:
		return c.colorize(iconFailed)
	case verification.Canceled:
		return c.colorize(iconCanceled)
	default:
		return ""
	}
}

func (c *RolloutViewController) calculateActivePhase() int {
	for i, rolloutPhase := range c.rollout.Phases {
		for _, uncompletedStatus := range phase.UncompletedStatuses {
			phaseStatus := rolloutPhase.Status
			if phaseStatus == uncompletedStatus {
				return i + 1
			}
		}
	}

	return len(c.rollout.Phases)
}

func (c *RolloutViewController) orderVerifications(verifications []verification.Verification) []verification.Verification {
	sort.Slice(verifications, func(i, j int) bool {
		return verification.StatusOrder[verifications[i].Status] < verification.StatusOrder[verifications[j].Status]
	})

	return verifications
}

func (c *RolloutViewController) printPhasesNumber() {
	for _, status := range rollout.CompletedStatuses {
		if status == c.rollout.Status {
			fmt.Fprintf(c.writer, tableFormat, "Phases:", fmt.Sprintf("%d", len(c.rollout.Phases)))
			return
		}
	}
	fmt.Fprintf(c.writer, tableFormat, "Phases:", fmt.Sprintf("%d/%d", c.calculateActivePhase(), len(c.rollout.Phases)))
}

func (c *RolloutViewController) printBackgroundVerifications() {
	writer := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', tabwriter.TabIndent)
	c.writer = writer
	fmt.Fprint(c.writer, fmt.Sprintf("  %s\t%s\t%s\n", c.colorize("METRICS"), c.colorize("VERIFICATION PROVIDER"), c.colorize("VERIFICATION STATUS")+c.iconStub()))

	for _, verificationItem := range c.orderVerifications(c.rollout.GetBackgroundVerifications()) {
		raw := fmt.Sprintf("  %s\t%s\t%s\n",
			c.colorize(verificationItem.MetricName),
			c.colorize(verificationItem.Provider),
			fmt.Sprintf("%s %s", c.verificationStatusIcon(verificationItem.Status), c.colorize(converter.VerificationStatus(verificationItem))),
		)
		fmt.Fprint(c.writer, raw)
	}
}

func (c *RolloutViewController) printVersion(label string, status rollout.VersionStatus, labelColor color.Attribute) {
	fmt.Fprintf(c.writer, fmt.Sprintf("%s:\n", label))
	c.printVersionStatus(status, labelColor)
	c.printReplicasNumber(status.Replicas)
}

func (c *RolloutViewController) printVersions() {
	newVersionLabel := model.CanaryLabel
	oldVersionLabel := model.StableLabel

	if c.rollout.Strategy == model.RollingUpdateStrategyType {
		newVersionLabel = model.NewVersionLabel
		oldVersionLabel = model.OldVersionLabel
	}

	if c.rollout.Status.IsCompleted() == false {
		c.printVersion(newVersionLabel, c.rollout.NewVersionStatus, color.FgYellow)
		c.printVersion(oldVersionLabel, c.rollout.StableVersionStatus, color.FgGreen)
	} else {
		if c.rollout.Status == rollout.Finished {
			c.printVersion(newVersionLabel, c.rollout.NewVersionStatus, color.FgGreen)
		} else {
			c.printVersion(oldVersionLabel, c.rollout.StableVersionStatus, color.FgGreen)
		}
	}
}

func (c *RolloutViewController) printStrategy() {
	strategy := model.CanaryLabel

	if c.rollout.Strategy == model.RollingUpdateStrategyType {
		strategy = model.RollingUpdateStrategyTypeLabel
	}

	fmt.Fprintf(c.writer, tableFormat, "Strategy:", strategy)
}
