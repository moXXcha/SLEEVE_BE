#!/bin/bash

# コミットメッセージチェックスクリプト
# docs/github.md のコミットメッセージ規約に従っているかをチェック

COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# 有効なprefix一覧（docs/github.mdより）
VALID_PREFIXES=(
    "feat"
    "feature"
    "enhance"
    "fix"
    "hotfix"
    "patch"
    "delete"
    "remove"
    "cleanup"
    "refactor"
    "reorganize"
    "rename"
    "test"
    "docs"
    "readme"
    "comment"
    "config"
    "env"
    "deps"
    "build"
    "ci"
    "perf"
    "optimize"
    "security"
    "auth"
    "db"
    "migration"
    "schema"
    "style"
    "format"
    "lint"
    "chore"
)

# コミットメッセージの1行目を取得
FIRST_LINE=$(echo "$COMMIT_MSG" | head -n 1)

# マージコミットやリバートコミットはスキップ
if [[ "$FIRST_LINE" =~ ^Merge || "$FIRST_LINE" =~ ^Revert ]]; then
    exit 0
fi

# prefix:の形式かチェック
if [[ ! "$FIRST_LINE" =~ ^[a-z]+: ]]; then
    echo "❌ コミットメッセージエラー: prefixが必要です"
    echo ""
    echo "正しい形式: [prefix]: [メッセージ]"
    echo ""
    echo "有効なprefix:"
    echo "  - feat: 新機能の追加"
    echo "  - fix: バグ修正"
    echo "  - refactor: リファクタリング"
    echo "  - test: テストの追加・修正"
    echo "  - docs: ドキュメントの更新"
    echo "  - config: 設定ファイルの変更"
    echo "  - など（詳細は docs/github.md を参照）"
    echo ""
    echo "現在のメッセージ:"
    echo "$FIRST_LINE"
    exit 1
fi

# prefixを抽出
PREFIX=$(echo "$FIRST_LINE" | sed -n 's/^\([a-z]*\):.*/\1/p')

# prefixが有効かチェック
VALID=false
for valid_prefix in "${VALID_PREFIXES[@]}"; do
    if [[ "$PREFIX" == "$valid_prefix" ]]; then
        VALID=true
        break
    fi
done

if [[ "$VALID" == false ]]; then
    echo "❌ コミットメッセージエラー: 無効なprefix '$PREFIX' です"
    echo ""
    echo "有効なprefix:"
    printf "  - %s\n" "${VALID_PREFIXES[@]}"
    echo ""
    echo "詳細は docs/github.md を参照してください"
    echo ""
    echo "現在のメッセージ:"
    echo "$FIRST_LINE"
    exit 1
fi

# コロンの後にスペースとメッセージがあるかチェック
if [[ ! "$FIRST_LINE" =~ ^[a-z]+:\ .+ ]]; then
    echo "❌ コミットメッセージエラー: prefixの後に空白とメッセージが必要です"
    echo ""
    echo "正しい形式: $PREFIX: [メッセージ]"
    echo "             ↑ここに空白が必要"
    echo ""
    echo "現在のメッセージ:"
    echo "$FIRST_LINE"
    exit 1
fi

# 成功
echo "✅ コミットメッセージチェック: OK"
exit 0
