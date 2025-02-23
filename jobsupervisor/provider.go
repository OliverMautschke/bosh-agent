//go:build !windows
// +build !windows

package jobsupervisor

import (
	"time"

	"code.cloudfoundry.org/clock"

	boshhandler "github.com/cloudfoundry/bosh-agent/handler"
	boshmonit "github.com/cloudfoundry/bosh-agent/jobsupervisor/monit"
	boshplatform "github.com/cloudfoundry/bosh-agent/platform"
	boshdir "github.com/cloudfoundry/bosh-agent/settings/directories"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const jobSupervisorListenPort = 2825

type Provider struct {
	supervisors map[string]JobSupervisor
}

func NewProvider(
	platform boshplatform.Platform,
	client boshmonit.Client,
	logger boshlog.Logger,
	dirProvider boshdir.Provider,
	handler boshhandler.Handler,
) (p Provider) {
	timeService := clock.NewClock()
	fs := platform.GetFs()
	runner := platform.GetRunner()
	monitJobSupervisor := NewMonitJobSupervisor(
		fs,
		runner,
		client,
		logger,
		dirProvider,
		jobSupervisorListenPort,
		MonitReloadOptions{
			MaxTries:               3,
			MaxCheckTries:          10,
			DelayBetweenCheckTries: 1 * time.Second,
		},
		timeService,
	)

	p.supervisors = map[string]JobSupervisor{
		"monit":      NewWrapperJobSupervisor(monitJobSupervisor, fs, dirProvider, logger),
		"dummy":      NewDummyJobSupervisor(),
		"dummy-nats": NewDummyNatsJobSupervisor(handler),
	}

	return
}

func (p Provider) Get(name string) (supervisor JobSupervisor, err error) {
	supervisor, found := p.supervisors[name]
	if !found {
		err = bosherr.Errorf("JobSupervisor %s could not be found", name)
	}
	return
}
