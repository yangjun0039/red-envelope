package network


//var DefaultAdapter = NewHandlerBuilder(SwitchRouteHandler).Adapt(RequesterHandler).
//	Then(RecoveryHandler, CrossOriginHanddler, VerifySignatureHandler)

var LoginHandler = 	NewHandlerBuilder(SwitchRouteHandler).Adapt(RequesterHandler).
    Then(RecoveryHandler, CrossOriginHanddler, LoginTokenHandler)

var TestAdapter = NewHandlerBuilder(SwitchRouteHandler).Adapt(RequesterHandler).
	Then(RecoveryHandler)

var DefaultAdapter = NewHandlerBuilder(SwitchRouteHandler).Adapt(RequesterHandler).
	Then(RecoveryHandler)