package mail

// MockDialer to simulate Dialer behavior in tests
// type MockDialer struct {
// 	mock.Mock
// }

// func (m *MockDialer) DialAndSend(msgs ...*gomail.Message) error {
// 	args := m.Called(msgs)
// 	return args.Error(0)
// }

// // Test for Send function
// func TestSend(t *testing.T) {
// 	mockDialer := new(MockDialer)

// 	mockDialer.On("DialAndSend", mock.Anything).Return(nil)

// 	err := Send(mockDialer, "recipient@example.com", "Test Subject", "Test Body")
// 	require.NoError(t, err)

// 	mockDialer.AssertExpectations(t)
// }
