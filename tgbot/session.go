package tgbot

type Session struct {
	ID         int64       `json:"id" comment:"用户ID"`
	LastBotId  int         `json:"last_bot_id" comment:"机器人的上一条消息ID"`
	LastUserId int         `json:"last_user_id" comment:"用户的上一条消息ID"`
	userdata   interface{} `json:"userdata" comment:"用户数据"`
	State      string      `json:"state" comment:"状态"`
	Step       int         `json:"step" comment:"步骤"`
	errorList  []int       `comment:"错误提示消息列表"`
}

func (s *Session) SaveBotID(id int) *Session {
	s.LastBotId = id
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) SaveUserID(id int) *Session {
	s.LastUserId = id
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) NewState(state string) *Session {
	s.State = state
	s.Step = 1
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) NextStep() *Session {
	s.Step++
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) ClearState() *Session {
	s.State = ""
	s.Step = 0
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) ResetStep() *Session {
	s.Step = 0
	SessionMgr.OnChange(s)
	return s
}

func (s *Session) GetUserData(val interface{}) interface{} {
	if s.userdata == nil {
		s.userdata = val
	}

	return s.userdata
}

func (s *Session) SetUserData(val interface{}) {
	s.userdata = val
	SessionMgr.OnChange(s)
}

func (s *Session) HasError() bool {
	return len(s.errorList) > 0
}

func (s *Session) AddError(err int) {
	s.errorList = append(s.errorList, err)
}

func (s *Session) GetErrors() []int {
	return s.errorList
}

func (s *Session) ClearErrors() {
	s.errorList = []int{}
}
