# test

## 実行

- `go test`でテスト
  - テストファイルはテスト対象と同じパッケージに置く
- 構造体`testing.T`の関数
  - `Error` = (`Log`エラーログを記録と`Fail`失敗を記録するが実行継続を許す)
  - `Fatal` = (`Log`と`FailNow`失敗を記録し実行を停止)


```
❯ ls
README.md  go.mod  main.go  main_test.go  post.json

❯ go test
PASS
ok      forest.work/test        0.003
```

- 詳細なあああ出力は，`-v`verboseフラグを付ける
- ソースコードのどの部分のテストが実行されたか出力するには，`-cover`カバレッジフラグ(cover)を付ける

```
❯ go test -v -cover
=== RUN   TestDecode
--- PASS: TestDecode (0.00s)
=== RUN   TestEncode
    main_test.go:21: Skipping encoding for now
--- SKIP: TestEncode (0.00s)
PASS
coverage: 50.0% of statements
ok      forest.work/test        0.002s
```

### テストをスキップする場合

- `-short` オプションを付けることで，テスト関数のはじめに`if testing.Short() { t.Skip("skip message") }`を書けばスキップできる
- 
```
❯ go test -v -cover -short
=== RUN   TestDecode
--- PASS: TestDecode (0.00s)
=== RUN   TestEncode
    main_test.go:22: Skipping encoding for now
--- SKIP: TestEncode (0.00s)
=== RUN   TestLongRunningTest
    main_test.go:27: Skipping long-running test in short mode
--- SKIP: TestLongRunningTest (0.00s)
PASS
coverage: 50.0% of statements
ok      forest.work/test        0.002s
```