
/kiro:spec-requirements "
複数の.gitignoreファイルをマージするgolangコマンドを作成する。

## 構想
- .gitignoreファイルのパスを受け取る。
- .gitignoreをパースして内容を読み込む。
- .gitignore 1、.gitignore 2、.gitignore 3と複数ファイル渡された時、先頭の.gitignoreファイルに、後続の.gitignoreファイルの内容を移植する。
- .gitignoreに#でセクションが区切られていて、同じ名前の場所が後続の.gitignoreにもあったら、同じ場所の設定を更新する。
- まったく同じパスの設定があって、互いに矛盾していたら、先頭の.gitignoreの方の設定が勝つ。

## 実装
- ライブラリは躊躇無く使って良い。
- テストも実装すること。
- A,B,Cの.gitignoreを混ぜたら、こうなります、みたいなテストと、資料を兼ねたテストも欲しい。
- ドキュメントも用意する。
"
think.

### 推奨アプローチ
- **Option B（全て新規作成）** を採用

ディレクトリ構造には、
ソースと、テストと、ドキュメント用のディレクトリを用意して。


この仕様でOK。一応、##、### と複数の場合や、#の直後に空白文字あるなしの違いがあっても、マッチして欲しい。
> 2. セクション名マッチングルール: 複数行ヘッダーの場合、最初の # 行のみで比較することを design.md で明確化する

単純ケースのみ対応
>  1. 矛盾パターンの定義範囲: path vs !path の単純ケースのみ対応するか、ワイルドカードパターン間の矛盾も扱うか →
>  現在の仕様（単純な否定パターンのみ）で進めるなら明示的にスコープを限定する

takt --task "/kiro-impl gitignore-merge 1"
/kiro-impl gitignore-merge 2

/kiro-impl gitignore-merge 3
/kiro-impl gitignore-merge 4

/kiro-impl gitignore-merge 5
/kiro-impl gitignore-merge 6

このコマンドはどこにコマンドがインストールされる？
go install gitignore-merge/cmd/gitignore-merge@latest


モジュール名を github.com/taku-o/gitignore-merge に変更。
>  リモートリポジトリで公開する場合は、go.mod のモジュール名を github.com/<user>/gitignore-merge のような形式に変更し、README
>  のインストール手順もそれに合わせる必要があります。


#で区切られている箇所があったら、マージ時に #の前の行に空白行を入れたい。
> gitignore-merge a.gitignore b.gitignore c.gitignore
> # Node
> node_modules/
> dist/
> 
> build/
> # Logs
> *.log
> !important.log
> 
> !debug.log
> # OS
> .DS_Store
> Thumbs.db
> # IDE
> .vscode/
> .idea/
> *.swp
> # Coverage
> coverage/
> *.coverprofile

Makefile


