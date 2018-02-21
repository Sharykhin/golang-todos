package handler

//func TestIndex(t *testing.T) {
//	savedTODOIndex := todoIndex
//	defer func() {
//		todoIndex = savedTODOIndex
//	}()
//	tt := []struct {
//		name         string
//		limitParam   string
//		offsetParam  string
//		todoFunc     func(ctx context.Context, limit, offset int) ([]entity.Todo, int, error)
//		expectStatus int
//	}{
//		{
//			name:         "bad limit",
//			limitParam:   "ds",
//			expectStatus: http.StatusBadRequest,
//		},
//		{
//			name:         "bad offset",
//			limitParam:   "1",
//			offsetParam:  "dfg",
//			expectStatus: http.StatusBadRequest,
//		},
//		{
//			name:        "in func",
//			limitParam:  "1",
//			offsetParam: "2",
//			todoFunc: func(ctx context.Context, limit, offset int) ([]entity.Todo, int, error) {
//				return nil, 0, fmt.Errorf("test case")
//			},
//			expectStatus: http.StatusInternalServerError,
//		},
//		{
//			name:        "all ok",
//			limitParam:  "1",
//			offsetParam: "2",
//			todoFunc: func(ctx context.Context, limit, offset int) ([]entity.Todo, int, error) {
//				return []entity.Todo{}, 12, nil
//			},
//			expectStatus: http.StatusOK,
//		},
//	}
//	var wg sync.WaitGroup
//
//	for _, tc := range tt {
//		wg.Add(1)
//		t.Run(tc.name, func(t *testing.T) {
//			defer wg.Done()
//			todoIndex = tc.todoFunc
//			e := httpexpect.WithConfig(httpexpect.Config{
//				Client: &http.Client{
//					Transport: httpexpect.NewBinder(Handler()),
//					Jar:       httpexpect.NewJar(),
//				},
//				Reporter: httpexpect.NewAssertReporter(t),
//			})
//			e.Request(http.MethodGet, "/").
//				WithQuery("limit", tc.limitParam).
//				WithQuery("offset", tc.offsetParam).
//				Expect().Status(tc.expectStatus)
//		})
//	}
//	wg.Wait()
//}
