package constant

import "time"

const UserName string = "userName"
const UserId string = "userId"
const TimeFormat string = "Jan _2 15:04:05 MST 2006"
const Following string = "Following"
const NotFollowing string = "Not Following"
const Port string = "9000"
const RelativePathForTemplate = "./web/template/"
const BackendServiceAddress = "localhost:9001"
const AuthServiceAddress = "localhost:9002"
const ContextTimeoutDuration = 5 * time.Second
const ProviderName = "memory"
const SessCookieName = "sessionId"
const MaxLifeTime = 7200
