# gitignore-merge

複数の `.gitignore` ファイルをセクション単位でインテリジェントにマージするコマンドラインツール。

## インストール

```bash
go install github.com/taku-o/gitignore-merge/cmd/gitignore-merge@latest
```

または、ソースからビルド:

```bash
git clone <repository-url>
cd gitignore-merge
go build -o gitignore-merge ./cmd/gitignore-merge/
```

## 使い方

```bash
gitignore-merge <file1> <file2> [file3...]
```

- 2つ以上の `.gitignore` ファイルパスを引数に指定する
- マージ結果は標準出力に出力される
- ファイルに書き出す場合はリダイレクトを使用する

```bash
gitignore-merge project.gitignore template.gitignore > .gitignore
```

## マージルール

### セクションマッチング

`.gitignore` ファイル内の `#` で始まるコメント行をセクションヘッダーとして認識する。セクション名はヘッダーの最初の行から `#` とスペースを除去した文字列で比較される。

以下は全て同じセクション名「Node」として扱われる:

```
# Node
## Node
#Node
```

### マージの優先順位

先頭のファイルがベースとなり、後続のファイルの内容が順番に統合される。

- **同名セクション**: 後続ファイルのパターンがベースの同名セクションに追加される
- **新規セクション**: ベースに存在しないセクションは末尾に追加される

### 重複除去

同一セクション内の完全一致するパターンは重複として除去される。

### 矛盾解決

パターン `path` と `!path` のように、完全一致する文字列の肯定・否定が矛盾する場合、ベースファイル（先頭ファイル）側のパターンが優先される。

例:
- ベースに `*.log` がある場合、後続の `!*.log` は無視される
- ベースに `!*.log` がある場合、後続の `*.log` は無視される

## 入出力例

### 入力

**file1.gitignore** (ベース):

```gitignore
# Node
node_modules/
dist/

# Logs
*.log
!important.log

# OS
.DS_Store
Thumbs.db
```

**file2.gitignore**:

```gitignore
# Node
node_modules/
build/

# Logs
*.log
!debug.log

# IDE
.vscode/
.idea/
```

### 実行

```bash
gitignore-merge file1.gitignore file2.gitignore
```

### 出力

```gitignore
# Node
node_modules/
dist/

build/
# Logs
*.log
!important.log

!debug.log
# OS
.DS_Store
Thumbs.db
# IDE
.vscode/
.idea/
```

この出力では:
- `# Node` セクション: `node_modules/` は重複除去、`build/` が追加された
- `# Logs` セクション: `*.log` は重複除去、`!debug.log` が追加された
- `# OS` セクション: ベースのみに存在するためそのまま保持
- `# IDE` セクション: file2 にのみ存在するため末尾に追加
