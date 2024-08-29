package tgbot

type Session struct {
	LastBotId  int         `json:"last_bot_id" comment:"机器人的上一条消息ID"`
	LastUserId int         `json:"last_user_id" comment:"用户的上一条消息ID"`
	UserData   interface{} `json:"userdata" comment:"用户数据"`
	State      string      `json:"state" comment:"状态"`
	Step       int         `json:"step" comment:"步骤"`
}

func (s *Session) NewState(state string) *Session {
	s.State = state
	s.Step = 1
	return s
}

func (s *Session) NextStep() *Session {
	s.Step++
	return s
}

func (s *Session) ClearState() *Session {
	s.State = ""
	s.Step = 0
	return s
}

func (s *Session) ResetStep() *Session {
	s.Step = 0
	return s
}

func (s *Session) GetUserData(val interface{}) interface{} {
	if s.UserData == nil {
		s.UserData = val
	}

	return s.UserData
}
