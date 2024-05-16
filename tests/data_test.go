package tests

//func (s *APITestSuite) TestSetData() {
//	r := s.Require()
//
//	login, password := "testingSetData", "testingSetData"
//
//	pass := s.serv.GeneratePasswordHash(password)
//
//	user := entity.User{Login: login, DBPassword: pass}
//
//	id, err := s.serv.CreateUser(context.Background(), user)
//	if err != nil {
//		s.Fail("Error checking data: " + err.Error())
//	}
//
//	jwt, err := s.serv.GenerateJWTToken(id)
//	if err != nil {
//		s.Fail("Error generating JWT: " + err.Error())
//	}
//
//	payload, testEvent := "test data", "testEvent"
//	inputBody := fmt.Sprintf(`{"payload":"%s"}`, payload)
//
//	req := httptest.NewRequest("POST", "api/set/{event}", bytes.NewReader([]byte(inputBody)))
//	req.Header.Set("Authorization", "Bearer "+jwt)
//
//	resp := httptest.NewRecorder()
//
//	rctx := chi.NewRouteContext()
//	rctx.URLParams.Add("event", testEvent)
//
//	req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//
//	route := chi.NewRouter()
//	s.handler.Register(route)
//
//	route.ServeHTTP(resp, req)
//
//	r.Equal(http.StatusOK, resp.Result().StatusCode)
//}
