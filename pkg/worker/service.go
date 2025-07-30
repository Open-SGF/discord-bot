package worker

import "discord-bot/pkg/worker/workerconfig"

type Service struct {
}

func NewService(config *workerconfig.Config) *Service {
	return &Service{}
}

func (s *Service) Execute() {}
